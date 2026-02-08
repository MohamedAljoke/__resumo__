# Content Ideas – Blog & Short Videos

This file is a living backlog of ideas for:

- Blog posts
- Short videos (YouTube / Shorts / Reels)
- Code showcases
- Architecture breakdowns

Status legend:

- 💡 Idea
- 🧪 Experiment
- ✍️ Writing
- 🎥 Recording
- ✅ Published

---

## 🔧 Libraries & Open Source

### GoQueue – DDD-friendly Job Queue in Go

Status: 💡

**Angles:**

- Why I built a queue instead of using X
- Applying DDD to infrastructure code
- Concurrency without shared state
- Workers, backpressure, retries
- Tradeoffs vs channels vs message brokers

**Formats:**

- Blog: Architecture + decisions
- Video: “Let’s build a queue in Go in 10 minutes”
- Video: Concurrency mistakes this library avoids

---

## 🎮 Game Development & Clean Code

### Social Deduction Game (Backend-focused)

Status: 💡

**Angles:**

- Modeling game rules with DDD
- Game loops as state machines
- Event-driven game logic
- Separating rules from I/O (clean architecture)
- Testing game logic without a UI

**Formats:**

- Blog: Clean Architecture applied to a game
- Video: “Your backend is just a game engine”
- Video: Turning game rules into pure functions

---

## 🔁 Concurrency & Async Thinking

### JS Generators ↔ Go Routines: Same Idea, Different Clothes

Status: 💡

**Core analogy:**

- `yield` vs `channel <-`
- `for await` vs `range channel`
- Backpressure in both worlds

**Formats:**

- Blog: Mental model comparison
- Video: “If you understand generators, you already understand Go”
- Code: Side-by-side examples (Node vs Go)

---

## 💰 Financial Systems & Transactions

### Idempotency in Financial APIs

Status: 💡

**Angles:**

- Why retries are dangerous in finance
- Idempotency keys explained
- Designing safe APIs for payments
- Common mistakes I’ve seen in production

**Formats:**

- Blog: Deep dive
- Video: “This bug costs millions”

---

### Eventual Consistency in Money Systems

Status: 💡

**Angles:**

- Why strong consistency is expensive
- Ledgers vs balances
- Compensations instead of rollbacks
- When _not_ to use transactions

**Formats:**

- Blog: Real-world examples
- Video: “Why your balance is lying to you”

---

### Ledger-first Architecture

Status: 💡

**Angles:**

- Append-only logs
- Rebuilding state from events
- Auditing and traceability
- Why finance loves immutability

**Formats:**

- Blog + diagrams
- Video: “Your database should be boring”

---

## ☁️ AWS & Distributed Systems

### When NOT to Use SQS

Status: 💡

**Angles:**

- FIFO vs Standard gotchas
- Visibility timeout traps
- Poison messages
- Why queues don’t fix bad architecture

---

### Lambda Is Not a Microservice

Status: 💡

**Angles:**

- Function ≠ service
- Hidden coupling
- Cold starts as design pressure
- When Lambdas shine (and when they don’t)

---

### Designing Retry-Safe Systems on AWS

Status: 💡

**Angles:**

- At-least-once delivery
- Idempotent handlers
- DLQs as observability tools
- Retrying humans vs machines

---

## 🧠 DDD & Architecture (Practical, Not Academic)

### Aggregates Are Not Database Tables

Status: 💡

**Angles:**

- Why most DDD implementations fail
- Invariants > schemas
- Transaction boundaries explained simply

---

### Domain Events vs Integration Events

Status: 💡

**Angles:**

- Internal vs external meaning
- Naming things properly
- Avoiding event spaghetti

---

### “Your Service Is Too Smart”

Status: 💡

**Angles:**

- Anemic domains vs god services
- Where logic actually belongs
- Refactoring stories

---

## ⚙️ Go & Node in the Real World

### Why Go Feels Boring (and That’s a Feature)

Status: 💡

**Angles:**

- Explicitness
- Predictability
- Why finance loves Go

---

### Node.js Is Great — Until It Isn’t

Status: 💡

**Angles:**

- CPU-bound workloads
- Event loop myths
- When to move to Go (and when not to)

---

### Writing Code for Humans, Not Frameworks

Status: 💡

**Angles:**

- Minimal abstractions
- Removing cleverness
- Code review lessons

---

## 🎥 Short Video Ideas (60–90s)

- “This retry bug will bankrupt you”
- “Generators are just Go routines in disguise”
- “Why balances are derived data”
- “Clean Architecture explained with a game”
- “Most people misuse DDD like this”
- “AWS won’t save bad design”
- “Your queue is lying to you”

---

## 📚 Book Notes → Content

Books / Topics to Extract Ideas From:

- Distributed systems
- Financial ledgers
- Event-driven architecture
- Reliability engineering

Rule:

> Every chapter should generate at least one blog or video idea.

---

## 🗑️ Parking Lot (Messy Thoughts)

- Weird bugs I’ve seen in prod
- Things I disagree with but need to articulate
- Analogies that helped juniors understand systems
