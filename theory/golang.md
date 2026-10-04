# Go

Interview review: each question has the answer to say out loud. **Bold = the keywords to hit.**
Based on *100 Go Mistakes and How to Avoid Them* (Teiva Harsanyi). Outputs verified on Go 1.25.

## 🔗 Syllabus

- [Core Concept: Everything Is a Copy](#core-concept-everything-is-a-copy)
- [Slices and Arrays](#slices-and-arrays)
- [Maps](#maps)
- [range Loops](#range-loops)
- [Memory Leaks](#memory-leaks)
- [Strings and Runes](#strings-and-runes)
- [Methods and Receivers](#methods-and-receivers)
- [Interfaces and nil](#interfaces-and-nil)
- [Error Handling](#error-handling)
- 🚧 Concurrency (coming next)

---

## Core Concept: Everything Is a Copy

**Q: Does Go pass by reference?**
No. **Everything is passed by value**: every assignment, argument, `range` value, and `m[k]` is a copy. But some values *contain* an address, so the copy **shares the target**:

| Type | A copy is… | After copying |
|---|---|---|
| `[N]T` array | **all the elements** | independent |
| `[]T` slice | **header `{ptr, len, cap}`** | **shares** the backing array |
| `map`, `chan`, `*T` | an address | **shares** the target |
| `struct` | all the fields | independent (pointer fields still share) |

**Q: If two variables share the same target, why doesn't reassigning one affect the other?**
Each holds **its own copy of the address**. They don't point to each other.
- **Write through** the address (`s[i] = x`, `*p = x`) → visible to both.
- **Reassign** the variable (`s = append(…)`, `p = &y`) → only that variable changes.
```go
q := p
p = &y      // *q still reads x
```

**Q: What do `&` and `*` mean?**
`&x` = **the address of** x. `*T` in a type = **pointer to T**. `*p` on a value = **go to that address**.
`*p = 7` changes the target. `p = &y` only changes what `p` holds.

**Q: Which types can be nil?**
Only **address types**: pointer, slice, map, chan, func, interface. A struct **can't be nil**: `var p MyStruct = nil` doesn't compile, and its empty state is the zero value `MyStruct{}`.

---

## Slices and Arrays

**Q: Array vs slice?**
An array `[N]T` **is** its elements: fixed size, and the size is part of the type. A slice `[]T` is a **header `{ptr, len, cap}`** that **views** a backing array. Assigning an array **copies every element**. Assigning a slice **copies only the header**, so both slices **share** the array.

**Q: What are `len` and `cap`? What are they after `s[lo:hi]`?**
`len` = elements you can index. `cap` = room until the **end of the backing array**. For `s[lo:hi]`: **len = hi−lo, cap = cap(s)−lo**. **Slicing never copies.**

**Q: Walk me through `append`.**
It checks for **room** (len < cap):
- **Room** → writes **in place** into the existing array, possibly **overwriting data another slice sees**.
- **No room** → allocates a **new, bigger array**, copies over, and returns a header to it. From then on, **no more sharing**.

Always write `s = append(s, x)`.

**Q: What does this print?**
```go
s := []int{1, 2, 3}
b := s[:2]
b = append(b, 99)
fmt.Println(s)
```
**`[1 2 99]`.** `b` has len 2 but **cap 3**, so `append` writes in place into `s[2]`.
*Follow-up: how do you prevent it?* → `s[:2:2]` (full slice expression, cap = 2) or `slices.Clip` forces a new array. `slices.Clone` gives full independence.
*Why is it dangerous?* → It **depends on spare capacity**, which nothing in the code shows. The same code can **pass tests and fail in production**.

**Q: What does this print?**
```go
func grow(s []int) {
	s = append(s, 4)
	s[0] = 99
}

a := []int{1, 2, 3}
grow(a)
fmt.Println(a)
```
**`[1 2 3]`.** `a` has cap 3, so `append` **allocated a new array**, and `s[0] = 99` wrote into **that** one. A function can change the caller's elements only while it still shares the array, and it can **never** change the caller's `len`. To grow a caller's slice, **return it**.

**Q: What does this print?**
```go
s := make([]int, 3)
s = append(s, 1)
fmt.Println(s)
```
**`[0 0 0 1]`.** `make([]int, 3)` creates **len 3** (three zeros), and `append` adds **after** them.
Use `make([]int, 0, 3)` (len 0, cap 3) when you plan to `append`, or `make([]int, 3)` + `s[i] = …` when you'll index.

**Q: How many elements does `copy(dst, src)` copy?**
**`min(len(dst), len(src))`.** It uses **len, not cap**:
```go
dst := make([]int, 0)       // len 0
copy(dst, []int{1, 2, 3})   // copies 0 → dst is []
```
Give `dst` a length first: `make([]int, len(src))`, or just use `slices.Clone(src)`.

**Q: nil slice vs empty slice?**
nil: `var s []T`, with **no backing array**. Empty: `[]T{}` or `make([]T, 0)`. Both have **len 0** and work the same with `len`, `range`, and `append`.
They differ in: **JSON** (nil → `null`, empty → `[]`) and **`reflect.DeepEqual`** / testify (not equal).
*Best practice:* check **`len(s) == 0`**, never `s == nil`. Return `nil` by default.

**Q: Can you compare slices with `==`?**
**No**, only against `nil`. Slices, maps, and funcs aren't comparable. Use **`slices.Equal`** / **`maps.Equal`** (Go 1.21). **Arrays are comparable**: `[3]int{1, 2} == [3]int{1, 2, 0}` is `true`. Avoid `reflect.DeepEqual` in hot code, because it's slow.

---

## Maps

**Q: Why doesn't `m["bob"].balance += 10` compile?**
`m[k]` returns a **copy**, and **map values aren't addressable**: the map can **move entries** when it grows.
Fix: `v := m[k]; v.balance += 10; m[k] = v`, or use `map[string]*Account`.
*Follow-up: why does `s[i].balance += 10` work for a slice?* → A slice element has a **fixed address** in its array.

**Q: What happens when you read from and write to a nil map?**
**Reading is fine** (zero value, `len` 0). **Writing panics**: `assignment to entry in nil map`. Always `make` before writing.

**Q: How do you tell "missing key" from "zero value"?**
The **comma-ok** idiom: `v, ok := m[k]`. `ok` is false if the key is missing.

**Q: Is map iteration ordered?**
**No. The order is random, on purpose**, and changes between runs. A key **added during iteration may or may not appear**. For a stable order, sort the keys: `slices.Sorted(maps.Keys(m))` (Go 1.23).

**Q: What can be a map key?**
Any **comparable** type: numbers, strings, pointers, channels, interfaces, **arrays**, structs made of comparable fields. **Not** slices, maps, or funcs.

**Q: Do maps shrink when you delete keys?**
**No, never.** `delete`, `clear(m)`, and even `maps.Clone` keep the full size. A map that spiked during peak traffic keeps that memory forever.
Fix: **rebuild** with `make(map[K]V, len(m))` + copy and drop the old map. Or store **pointers** for big values so the values themselves can be freed.

---

## range Loops

**Q: What does this print?**
```go
accounts := []account{{100}, {200}}
for _, a := range accounts {
	a.balance += 1000
}
fmt.Println(accounts)
```
**`[{100} {200}]`**, unchanged. **The range value is a copy.** Fix: **write through the index**: `accounts[i].balance += 1000`.
*And with `[]*account`?* → It **does** change, because `a` is a copy of the **pointer**.
*Did Go 1.22 change this?* → No. 1.22 gives a **fresh variable per iteration**, but it's **still a copy**.

**Q: How is the range expression evaluated?**
**Once, before the loop**, into a **hidden copy**. The loop iterates over that copy.

| `range` over | Hidden copy | Element writes during the loop | Reassign / `append` |
|---|---|---|---|
| array `a` | all the elements | ❌ invisible | ❌ |
| `&a` | the address | ✅ visible | ❌ |
| slice `s` | the header | ✅ visible | ❌ old len, old array |
| chan `ch` | the address | – | ❌ `ch = ch2` ignored |

**Q: Does this loop end?**
```go
s := []int{0, 1, 2}
for range s {
	s = append(s, 10)
}
```
**Yes, after 3 iterations** (the hidden header has len 3). The classic `for i := 0; i < len(s); i++` version **never ends**, because `len(s)` is re-evaluated each time.

**Q: What does this print?**
```go
a := [3]int{0, 1, 2}
for i, v := range a {
	a[2] = 10
	if i == 2 { fmt.Println(v) }
}
```
**`2`.** `range` copied the **whole array** first. With a **slice** → `10` (the header shares the array). With **`range &a`** → `10`.

**Q: What does this print?**
```go
s := []int{0, 1, 2}
for i, v := range s {
	if i == 0 {
		s = append(s, 99) // cap full → new array
		s[2] = 10
	}
	if i == 2 { fmt.Println(v) }
}
```
**`2`.** `s` moved to a **new array**, while the hidden header still reads the **old** one. Without the `append` → `10`.

---

## Memory Leaks

**Q: How can a slice cause a memory leak?**
**The GC frees whole arrays, never part of one.** Keeping `small := big[:5]` of a 1 MB buffer keeps **the whole 1 MB alive**.
Fix: **copy what you keep**: `bytes.Clone`, `strings.Clone`, `slices.Clone`.
*Does `big[:5:5]` fix it?* → **No.** It limits `append`, not the GC.

**Q: And with a slice of pointers?**
`s = s[:2]` hides the other slots, but **the GC still scans the whole array**, so what they point to stays alive. Fix: **set removed slots to nil** before shrinking, or clone. (`slices.Delete` / `slices.Compact` zero them since Go 1.22.)

**Q: Name real code where this happens.**
An ID kept from a request body, a substring of a big file, map keys from `strings.Split`, `regexp` submatches, popping from a queue of pointers (`q[0] = nil` before `q = q[1:]`).
**Reused buffers** are the opposite problem: with `bufio.Scanner.Bytes()`, `sql.RawBytes`, or `sync.Pool`, your piece **silently changes** later. → **Clone what you keep.**
*What already copies?* → `string(b)`, `[]byte(s)`, `scanner.Text()`.

---

## Strings and Runes

**Q: What is a string in Go?**
A **read-only sequence of bytes**, not characters: a header **`{ptr, len}`** (a slice without `cap`). It's **immutable**.

**Q: What does `len` return for a string?**
**Bytes, not characters.** The bug hides until the text has `é`, `汉`, or an emoji.
```go
len("hello")  // 5
len("hêllo")  // 6   ê = 2 bytes
len("汉")     // 3   1 character, 3 bytes
len("汉字!")  // 7
```
To count characters: **`utf8.RuneCountInString(s)`** → `"汉字!"` = 3.

**Q: What is a rune?**
**One character**: a Unicode **code point**, `type rune = int32`. UTF-8 stores it in **1–4 bytes**:
```
h  [68]           1 byte
ê  [C3 AA]        2 bytes
汉 [E6 B1 89]     3 bytes
😀 [F0 9F 98 80]  4 bytes
```
**Charset** (Unicode) = characters with numbers (`汉` = `U+6C49`). **Encoding** (UTF-8) = how a number becomes bytes (`E6 B1 89`).

**Q: What does this print?**
```go
s := "hêllo"
for i, r := range s {
	fmt.Print(i, ":", string(r), " ")
}
fmt.Printf("%c", s[1])
```
**`0:h 1:ê 3:l 4:l 5:o`**, then **`Ã`**.
- `range` over a string yields **runes**, and `i` is each rune's **starting byte index**, so it **jumps 1 → 3** (ê is 2 bytes).
- `s[i]` is a **byte**, not a rune. `s[1]` is only the first byte of ê, printed as `Ã`.
- To index by character: `[]rune(s)[1]` (this allocates).

**Q: Are Go strings always UTF-8?**
**No.** **Literals** are UTF-8 (source files are). Strings from files, the network, or `string(bytes)` can hold **any bytes**, and nothing checks them.

**Q: `strings.TrimRight` vs `strings.TrimSuffix`?**
```go
strings.TrimRight("123oxo", "xo")   // "123"   removes any of the chars {x, o} from the right, repeatedly
strings.TrimSuffix("123oxo", "xo")  // "123o"  removes the exact suffix "xo" once
```
**`TrimRight`/`TrimLeft` take a set of characters. `TrimSuffix`/`TrimPrefix` take an exact string.**

---

## Methods and Receivers

**Q: Value receiver vs pointer receiver?**
**Calling a method copies the receiver.** **Value** → the method gets a **copy of the struct**, and changes are thrown away. **Pointer** → it gets the **address**, so changes reach the caller.
*Does the call look different?* → **No.** `c.add(50)` works for both (Go passes `&c` for you). **Only the signature** shows it.

**Q: What does this print?**
```go
type customer struct{ balance float64 }
func (c customer) add(v float64) { c.balance += v }

c := customer{balance: 100}
c.add(50)
fmt.Println(c.balance)
```
**`100`.** Fix: `func (c *customer) add`.

**Q: A value receiver, but the struct has a pointer field. Can the method change the caller's data?**
```go
type customer struct{ data *data }
func (c customer) add(v float64) {
	c.data.balance += v           // ✅ writes THROUGH the pointer → caller sees it
	c.data = &data{balance: 999}  // ❌ reassigns the copy's own field → caller doesn't
}
```
The copy has **its own `data` field** (created **at the call**) holding **the same address**.

**Q: When must a receiver be a pointer or a value?**
- **Must be pointer:** the method **mutates** · the struct holds a **`sync.Mutex`** (copying breaks it)
- **Should be pointer:** **large** struct
- **Must be value:** **map / func / chan** receiver types · enforcing **immutability**
- **Should be value:** `int`, `string`, **`time.Time`**, small value-like structs

**Default: value. If you're unsure, use a pointer.** Don't mix the two on one type.

**Q: Can you call a method on a nil pointer?**
**Yes.** A method is a function whose first parameter is the receiver. `var f *Foo; f.Bar()` works if `Bar` doesn't dereference `f`.

---

## Interfaces and nil

**Q: What is an interface value under the hood?**
**A box with two slots: `[ type | value ]`.** It's **nil only when *both* slots are empty.**

**Q: What does this print?**
```go
type MyErr struct{}
func (e *MyErr) Error() string { return "boom" }

func validate() error {
	var e *MyErr   // nil pointer
	return e
}

fmt.Println(validate() == nil)
```
**`false`.** Returning `e` as an `error` **fills the type slot**:
```
return nil  →  [ nil    | nil ]  → err == nil: true
return e    →  [ *MyErr | nil ]  → err == nil: false
```
Fix: **never return a typed nil pointer as an interface. Return a literal `nil`.** This applies to **any interface**, not only `error`.

**Q: If `p` is a nil pointer, where does the type in `var err error = p` come from?**
**From the compiler.** A pointer has **no type slot** (8 bytes, just the address). An interface stores **type + value (16 bytes)**. At `err = p`, the compiler writes **`p`'s declared type** (`*MyErr`, the **concrete** type) into the type slot, even though the address is nil.

**Q: Can a `MyErr` *value* be used as an `error` here?**
**No.** `Error()` has a **pointer receiver**, so only `*MyErr` has the method: `MyErr does not implement error (method Error has pointer receiver)`.

---

## Error Handling

**Q: `%w` vs `%v` in `fmt.Errorf`?**
**`%w` wraps**: the source error stays **inside**, so `errors.Is/As` can find it. **`%v` transforms**: the source becomes **text only**. **Both print the same message.**
*Why not always `%w`?* → It makes the inner error **part of your API**. If callers check `errors.Is(err, sql.ErrNoRows)` and you switch databases, their check **silently breaks**.
*Best pattern:* **translate** implementation errors into your own:
```go
if errors.Is(err, sql.ErrNoRows) {
	return fmt.Errorf("get user %d: %w", id, ErrNotFound)  // wrap YOUR sentinel
}
return fmt.Errorf("get user %d: %v", id, err)              // hide the driver's error
```

**Q: What does this print?**
```go
err := fmt.Errorf("ctx: %w", transientError{…})
switch err.(type) {
case transientError: fmt.Println("503")
default:             fmt.Println("400")
}
```
**`400`.** A **type switch only sees the outer box** (`*fmt.wrapError`). The same goes for `err == ErrX`: with `ErrX` wrapped, **`==` is false**.

**Q: `errors.Is` vs `errors.As`?**
| | Asks | Use for |
|---|---|---|
| **`errors.Is(err, ErrNotFound)`** | is this **exact value** anywhere in the chain? | **sentinels**: `var ErrX = errors.New(…)`, `os.ErrNotExist`, `io.EOF` |
| **`errors.As(err, &target)`** | is there an error **of this type**? If so, fill `target` | **custom types** with fields: `*os.PathError` → `.Path` |
```go
var pe *os.PathError
if errors.As(err, &pe) {        // the variable's type = the type in the chain
	fmt.Println(pe.Path)
}
```
*Gotcha:* the target **must be a pointer** (`&pe`). `errors.As(err, pe)` compiles but **panics** (`go vet` catches it).
🆕 Go 1.26: `pe, ok := errors.AsType[*os.PathError](err)`.

**Q: Sentinel error vs error type: when to use which?**
**Sentinel** (`var ErrNotFound = errors.New(…)`) for **expected** errors that callers branch on, checked with `errors.Is`. **Custom type** when the error needs to **carry data** (status code, path, retry-after), checked with `errors.As`.

**Q: How do you combine several errors?**
**`errors.Join(e1, e2)`** (Go 1.20). The message is each one on its own line, and **`errors.Is` / `errors.As` match any of them**.

**Q: Should you log an error and also return it?**
**No. Handle an error once:** either **log it** (and stop), **or return it** (with context). Doing both prints the same failure several times up the call stack and makes logs confusing.

**Q: Can you ignore an error?**
Only **explicitly**: `_ = notify()` with a comment saying why. A silently dropped error looks like a bug to the next reader.
Watch `defer f.Close()` on **writes**: `Close` can fail (unflushed data), so handle its error, for example by assigning it to a named result in the deferred func.

**Q: When should you `panic` instead of returning an error?**
Almost never. Only for **programmer bugs** (impossible state, a nil dependency at startup) or failing **mandatory setup** (`regexp.MustCompile`, config at init). Everything that can fail at runtime (I/O, input, network) returns an **`error`**.

---

## 30-second summary
> **Go copies on every assignment, and what's copied depends on the type.** Arrays copy their elements, while slices, maps, channels, and pointers copy an address that **shares the target**. That explains `append` overwriting other slices, `range` values being copies, and small slices keeping big buffers alive. **Reassigning a variable never changes anyone else's copy.** A string is **bytes**, so `len` counts bytes and `range` gives runes. An **interface is `[type | value]`**, nil only when both are empty, so never return a typed nil pointer. Wrap errors with **`%w`** only when callers should see the source, and check them with **`errors.Is`** (values) and **`errors.As`** (types), never with `==` or a type switch.
