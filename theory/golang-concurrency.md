# Go: Concurrency and Channels

Interview review: each question has the answer to say out loud. **Bold = the keywords to hit.**
Based on *100 Go Mistakes and How to Avoid Them* (Teiva Harsanyi), chapters 8–9. Outputs verified on Go 1.25.
Part of the [Go review](./golang.md).

## 🔗 Syllabus

- [Concurrency vs Parallelism](#concurrency-vs-parallelism)
- [Goroutines and the Scheduler](#goroutines-and-the-scheduler)
- [Data Races and Race Conditions](#data-races-and-race-conditions)
- [The Memory Model](#the-memory-model)
- [Mutex or Channel?](#mutex-or-channel)
- [Channels](#channels)
- [select](#select)
- [Channel Size](#channel-size)
- [Context](#context)
- [Goroutine Leaks](#goroutine-leaks)
- [The sync Package](#the-sync-package)
- [Shared Slices and Maps](#shared-slices-and-maps)
- [Patterns](#patterns)

---

## Concurrency vs Parallelism

**Q: What's the difference between concurrency and parallelism?**
**Concurrency is about structure**: splitting a problem into independent steps that coordinate. **Parallelism is about execution**: running the same step many times at once. *"Dealing with lots of things at once vs doing lots of things at once"* (Rob Pike).
Coffee shop: separate accept order → grind → brew with queues between the steps = **concurrency**. Adding 3 grinders = **parallelism** in one step. **Concurrency enables parallelism.**

**Q: Is a concurrent version always faster?**
**No.** The book's parallel merge sort, with a goroutine for every split, was **~8× slower** than the sequential one (17.5 s vs 2.3 s). Each goroutine had **too little work** for what it costs to create and schedule. Adding a **threshold** (sequential below 2048 elements) made it **~40% faster**. Start sequential, then **benchmark and profile**.

---

## Goroutines and the Scheduler

**Q: Goroutine vs OS thread?**
A goroutine is scheduled by the **Go runtime** onto OS threads, and an OS thread is scheduled by the OS onto a core. A goroutine starts with a **~2 KB stack** that grows as needed; an OS thread usually has ~2 MB. Switching goroutines is much cheaper, but **not free**.

**Q: Explain G / M / P.**
**G** = goroutine, **M** = OS thread (machine), **P** = logical processor. **`GOMAXPROCS`** = the number of Ps, which caps how many Ms run Go code at the same time (more Ms can exist when threads are blocked in syscalls).
Each P has a **local run queue**, and there is one **global queue**. An idle P **steals work** from other Ps. The scheduler has been **preemptive since Go 1.14**: a goroutine that runs for ~10 ms can be switched out.

**Q: How many goroutines for a worker pool?**
- **CPU-bound** → about **`runtime.GOMAXPROCS(0)`** (not `NumCPU`, since GOMAXPROCS can be lower). More only adds context switches.
- **I/O-bound** → decided by **what the external system can handle** (DB pool size, API rate limit).

🆕 Go 1.25: the default GOMAXPROCS is **container-aware** (it follows the cgroup CPU limit).

---

## Data Races and Race Conditions

**Q: Data race vs race condition?**
- **Data race**: **≥2 goroutines access the same memory at the same time, and at least one writes.** The result is undefined. `i++` is read → increment → write, so two goroutines can both write `1`.
- **Race condition**: the **behaviour depends on timing or order** you don't control.

Code can be **free of data races and still have a race condition**:
```go
go func() { mu.Lock(); i = 1; mu.Unlock() }()
go func() { mu.Lock(); i = 2; mu.Unlock() }()
// no data race (it's locked), but i ends as 1 OR 2 → race condition
```
A lock fixes a data race. A race condition needs **coordination** (ordering).

**Q: Three ways to fix `i++` from two goroutines?**
1. **`sync/atomic`**: `var i atomic.Int64; i.Add(1)` (only for numbers and pointers)
2. **Mutex** around the critical section
3. **Channel**: the goroutines send `1`, and **one goroutine owns `i`** and adds them up

**Q: How do you find data races?**
`go test -race` / `go run -race`. It only reports races that **actually happen** during that run, so run it on tests that exercise concurrency.

---

## The Memory Model

**Q: What does "happens-before" guarantee across goroutines?**
- **`go f()` happens before f starts.** A goroutine **exiting** guarantees nothing, so `go func(){ i++ }(); fmt.Println(i)` is a race.
- **A send happens before the matching receive completes.**
- **`close(ch)` happens before a receive that sees the close.**
- **Unbuffered: the receive happens before the send completes.**

**Q: Is this a data race?**
```go
i := 0
ch := make(chan struct{}, 1)   // buffered
go func() { i = 1; <-ch }()
ch <- struct{}{}
fmt.Println(i)
```
**Yes.** A **buffered** send doesn't wait for the receive, so `i = 1` and `Println(i)` can run at the same time (`-race` reports it). **Make it unbuffered and the race is gone**: the send can't finish until the receive happens, and the receive comes after `i = 1`.

---

## Mutex or Channel?

**Q: When do you use a mutex, and when a channel?**
- **Parallel goroutines** (the same step sharing state, like HTTP handlers updating a cache) need **synchronization** → **mutex**.
- **Concurrent goroutines** (different steps) need **coordination** or **transfer of ownership** (passing a result on to the next stage) → **channel**.

Don't force channels everywhere because of "share memory by communicating". They **complement** each other.

---

## Channels

**Q: Unbuffered vs buffered?**
- **Unbuffered** (`make(chan T)`): the sender **blocks until a receiver takes it**. You get **synchronization**: both goroutines are at a known point.
- **Buffered** (`make(chan T, n)`): the sender only blocks when the **buffer is full**. You get **communication only, no synchronization**.

**Q: What does this print?**
```go
ch := make(chan int)
close(ch)
fmt.Println(<-ch, <-ch)
v, ok := <-ch
fmt.Println(v, ok)
```
**`0 0`** then **`0 false`**. A **receive from a closed channel never blocks** and returns the **zero value**. Use **`v, ok`** to tell a real `0` from a close. (`range ch` stops on close by itself.)

**Q: What happens on a nil channel? A closed one?**
| Operation | nil channel | closed channel |
|---|---|---|
| send | **blocks forever** | **panic** |
| receive | **blocks forever** | zero value, `ok = false` |
| close | **panic** | **panic** |

**Q: Why would you ever *want* a nil channel?**
To **disable a case in a `select`**. Merging two channels: when one closes, **set it to `nil`** so the select stops picking it. Without that you'd get endless zero values or a **busy loop**.
```go
for ch1 != nil || ch2 != nil {
	select {
	case v, ok := <-ch1:
		if !ok { ch1 = nil; break }
		out <- v
	case v, ok := <-ch2:
		if !ok { ch2 = nil; break }
		out <- v
	}
}
close(out)
```

**Q: How do you send a signal with no data?**
**`chan struct{}`**. `struct{}` is **0 bytes**, and the type says "no meaning in the value". `chan bool` makes readers wonder what `false` means. The same idea gives a set: `map[K]struct{}`.

**Q: Who should close a channel?**
**The sender** (the owner), and **only once**. The receiver never closes. With several senders, coordinate (WaitGroup, then close).

---

## select

**Q: If several cases are ready, which one runs?**
**One at random** (uniform), **not the first in source order**. This prevents starvation.
*Trap:* a producer sends 10 messages to a buffered `messageCh` and then sends on `disconnectCh`. The receiver may handle 0, 5 or 10 messages before it sees the disconnect.
*Fixes:*
- **Single producer:** make `messageCh` unbuffered, or use **one channel** for both kinds of event (a channel is FIFO).
- **Multiple producers:** on disconnect, **drain** the rest:
```go
case <-disconnectCh:
	for {
		select {
		case v := <-messageCh: handle(v)
		default: return            // runs only when nothing else is ready
		}
	}
```

---

## Channel Size

**Q: How do you choose a buffer size?**
**Unbuffered by default** (it synchronizes, and deadlocks show up right away). If you need a buffer, **start at 1**. Use a different size only with a reason:
- a **worker pool** → the number of workers
- **rate limiting** → the limit
- anything else → **benchmark it and comment why** (`make(chan int, 40)`: why 40?)

Queues are almost always **nearly full or nearly empty**. A big buffer usually **hides back-pressure** rather than solving it.

---

## Context

**Q: What does a `context.Context` carry?**
A **deadline**, a **cancellation signal**, and **request-scoped values**, across API boundaries. A function that callers wait on should take `ctx` as its **first parameter**.

**Q: Why `defer cancel()`?**
`WithTimeout`, `WithDeadline` and `WithCancel` hold resources (a **timer** and a link to the parent) until the context is cancelled. `defer cancel()` releases them **as soon as you return**. `go vet` (lostcancel) warns when it's missing.

**Q: How does a goroutine notice cancellation?**
`ctx.Done()` returns a channel that gets **closed**, because **only a close reaches every receiver**. Then `ctx.Err()` is `context.Canceled` or `context.DeadlineExceeded`. Never block on a bare channel operation inside a context-aware function:
```go
select {
case <-ctx.Done():
	return ctx.Err()
case ch <- v:
}
```

**Q: How should you define a context value key?**
With an **unexported custom type**: `type key string; const userKey key = "user"`. A plain string key can **collide** with another package's key. Use values only for **request-scoped data** (trace ID, auth info), never for ordinary function parameters.

**Q: What's wrong with this handler?**
```go
go publish(r.Context(), resp)   // async, to Kafka
writeResponse(w, resp)
```
**The request context is cancelled once the response is written** (also when the client disconnects). The publish races the response and often gets cancelled. Use **`context.WithoutCancel(r.Context())`** (Go 1.21): it **keeps the values** (trace ID) and **drops cancellation and the deadline**. `context.Background()` would lose the values.

**Q: Which context to use when you don't have one yet?**
**`context.TODO()`**, which marks it as "not decided yet". `context.Background()` is for the top level (main, tests, the start of a request).

---

## Goroutine Leaks

**Q: What does a goroutine leak cost?**
Its stack (2 KB and growing), **every heap object it references**, and the **resources it holds** (connections, files, sockets). Leaks build up until you run out of memory.

**Q: How do you avoid leaks?**
**Every `go` statement needs a stop plan**: *what makes it return, and who waits for it?*
- `for v := range ch` exits only when **someone closes `ch`**.
- A goroutine blocked on a send or receive that will never happen **leaks forever**.
- Signalling with a context **isn't enough** if the goroutine owns resources, because cancel doesn't **wait** for cleanup. Give the type a `Close()` that signals **and blocks until the goroutine is done**, then `defer w.Close()` in main.

Testing tools: `go.uber.org/goleak`, and 🆕 **`testing/synctest`** (GA in Go 1.25) for deterministic concurrent tests.

**Q: What did this print before Go 1.22?**
```go
for _, i := range []int{1, 2, 3} {
	go func() { fmt.Print(i) }()
}
```
Often **`333`** or `233`, because every closure shared **one `i`**. 🆕 **Since Go 1.22, each iteration has its own variable**, so it prints 1, 2 and 3 in some order. Old fix: `i := i`, or pass it as an argument.

---

## The sync Package

**Q: What's wrong here, and what does it print?**
```go
var wg sync.WaitGroup
var v atomic.Int64
for range 3 {
	go func() {
		wg.Add(1)
		v.Add(1)
		wg.Done()
	}()
}
wg.Wait()
fmt.Println(v.Load())
```
**Anything from 0 to 3.** On Go 1.25 it printed **`0` in 20 of 20 runs**. `wg.Add` runs **inside** the goroutine, so `Wait` sees a counter of 0 and returns before the goroutines have even started. **Call `Add` in the parent, before `go`.**
🆕 Go 1.25: **`wg.Go(func(){...})`** does Add, go and Done for you, and `go vet` flags *"WaitGroup.Add called from inside new goroutine"*.

**Q: Why does this deadlock?**
```go
func (c *Customer) UpdateAge(age int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if age < 0 {
		return fmt.Errorf("bad age for %v", c)   // %v calls c.String()
	}
	...
}
func (c *Customer) String() string {
	c.mu.RLock()                                 // already locked → waits forever
	...
}
```
**`%v` / `%s` call `String()`**, and **Go mutexes aren't reentrant**. Result: `fatal error: all goroutines are asleep - deadlock!` *Fix:* validate **before** locking, or format the fields directly (`c.id`).

**Q: What's wrong with this?**
```go
type Counter struct {
	mu sync.Mutex
	n  map[string]int
}
func (c Counter) Inc(k string) { c.mu.Lock(); defer c.mu.Unlock(); c.n[k]++ }
```
The **value receiver copies the mutex**, so each call locks **its own copy** and nothing is protected → data race. **Never copy** `Mutex`, `RWMutex`, `WaitGroup`, `Cond`, `Once`, `Map` or `Pool`. Use a **pointer receiver**. `go vet` says: *"Inc passes lock by value: Counter contains sync.Mutex"*.

**Q: Mutex vs RWMutex?**
`RWMutex` allows **many readers at once or one writer**. It helps when reads far outnumber writes; otherwise a plain `Mutex` is simpler and often faster.

**Q: When do you need `sync.Cond`?**
To **broadcast repeatedly to many waiting goroutines**. A channel message reaches **only one** receiver, and a `close` reaches everyone but **only once**.
```go
c.L.Lock()
for !condition() {   // ALWAYS a loop: re-check after every wake-up
	c.Wait()         // unlock → sleep → re-lock
}
c.L.Unlock()
// updater: lock, change state, unlock, c.Broadcast()
```
Downside: if nobody is waiting, the **Broadcast is lost**.

---

## Shared Slices and Maps

**Q: Is this a data race?**
```go
s := make([]int, 0, 1)
go func() { _ = append(s, 1) }()
go func() { _ = append(s, 2) }()
```
**Yes.** There's **spare capacity**, so both appends write **index 0 of the same backing array**. With `make([]int, 1)` (full) each append gets a **new array** → no race. Don't rely on that difference: **never `append` to a shared slice from several goroutines.** Copy it first.

**Q: Does this protect the map?**
```go
c.mu.RLock()
balances := c.balances   // copy under the lock
c.mu.RUnlock()
for _, b := range balances { sum += b }
```
**No.** Assigning a map (or slice) copies only the **header**, so both variables share the same data, and reading it outside the lock races with writers. *Fix:* hold the lock for the **whole loop**, or make a **real copy** under the lock (`maps.Clone`) when the work after it is slow.

**Q: Which concurrent slice and map accesses are races?**
- Slice: **same index** with a writer → race. **Different indexes** → fine.
- Map: **any write** while another goroutine reads or writes → race, **even with different keys**.

---

## Patterns

**Q: Write a worker pool.**
```go
jobs := make(chan Job)            // or buffered with size = workers
var wg sync.WaitGroup
for range runtime.GOMAXPROCS(0) { // CPU-bound
	wg.Go(func() {
		for j := range jobs {     // exits when jobs is closed
			process(j)
		}
	})
}
for _, j := range all { jobs <- j }
close(jobs)                       // the sender closes
wg.Wait()
```

**Q: Fan-out calls where any error should cancel the rest?**
**`errgroup`** (`golang.org/x/sync/errgroup`):
```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(10)                         // at most 10 at once
for i, c := range circles {
	g.Go(func() error {
		r, err := call(ctx, c)
		if err != nil { return err }   // cancels ctx for the others
		results[i] = r                 // a different index per goroutine → no race
		return nil
	})
}
if err := g.Wait(); err != nil { return nil, err }   // the first error
```
Cancellation only helps if `call` **respects `ctx`**.

**Q: Fan-in (merge channels)?**
One goroutine `select`s over the inputs, **sets each closed input to `nil`**, and **closes the output** when every input is nil (see [Channels](#channels)).

---

## 30-second summary
> **Concurrency is structure, parallelism is execution**, and goroutines are cheap but **not free**, so benchmark before going parallel. Size pools to **GOMAXPROCS** for CPU work and to the **external limit** for I/O. A **data race** is unsynchronized access with a writer; a **race condition** is timing-dependent behaviour, and you can have one without the other. Use a **mutex for shared state** and a **channel for coordination**. **Unbuffered channels synchronize, buffered ones don't**: start unbuffered, or with a buffer of 1. A closed channel returns **zero values** (check `ok`); a **nil channel blocks**, which **disables a select case**; `select` picks **randomly**. Pass **`ctx`**, always **`defer cancel()`**, and use **`WithoutCancel`** for work that must outlive the request. **Every goroutine needs a stop plan.** Call `wg.Add` **before** `go` (or use `wg.Go`), **never copy sync types**, remember that **`%v` calls `String()`**, and never `append` to or header-copy a shared slice or map without a lock.
