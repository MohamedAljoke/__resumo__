# NodeJS

### Q1: Explain the event loop.
- Event loop handles async operations in phases (timers, I/O callbacks, idle, poll, check, close).
- Node.js is single-threaded but uses libuv + event loop to handle concurrency.
- That’s why it’s efficient for I/O-heavy tasks but not CPU-heavy ones.

---

### Q2: Difference between `process.nextTick()`, `setImmediate()`, and `setTimeout()`?
- `nextTick`: runs after current operation, before next event loop phase.
- `setImmediate`: runs at the check phase, after I/O events.
- `setTimeout(fn, 0)`: runs after timers phase, not immediate.

---

### Q3: How do you handle errors in async/await code?
- Use `try/catch` for sync + async code.
- Centralized error middleware in Express.
- Always return meaningful error messages, not stack traces.

---

### Q4: How do you secure an API in Node.js?
- JWT/OAuth2 authentication.
- Input validation/sanitization.
- Rate limiting, CORS.
- HTTPS only.
- Use parameterized queries (avoid SQL injection).

---

### Q5: Why is Node.js great for I/O but not CPU-heavy tasks?
- Because it’s single-threaded — CPU-bound tasks block the event loop.
- Use worker threads or offload to background services for CPU-heavy work.