# ReactJS

### Q1: What’s the difference between `useMemo` and `useCallback`?
- `useMemo`: memoizes the **result of a function**.
- `useCallback`: memoizes the **function itself**.
- Use `useMemo` to avoid recalculating expensive values, and `useCallback` to avoid re-creating functions unnecessarily.

---

### Q2: Why do we use `key` in React lists?
- Keys help React identify which items changed, were added, or removed.
- Without keys, React re-renders unnecessarily and may produce bugs (like losing input focus).
- Keys must be unique and stable (not array index if list changes).

---

### Q3: What’s the difference between controlled and uncontrolled components?
- **Controlled**: React controls the form input via state (`value` + `onChange`).
- **Uncontrolled**: DOM maintains the value (`defaultValue`, use `ref`).
- Controlled = predictable, better for validation. Uncontrolled = simpler, less overhead.

---

### Q4: How do you optimize performance in React?
- `React.memo` to prevent unnecessary re-renders.
- `useMemo` / `useCallback` for expensive calculations or stable references.
- Code splitting (dynamic import).
- Virtualization for large lists.
- Avoid inline functions/objects where possible.

---

### Q5: What’s reconciliation in React?
- Reconciliation = React’s process of diffing Virtual DOM trees to update the real DOM efficiently.
- Uses heuristics (keys, element types) to decide if it reuses or re-creates elements.

---

### Q6: Explain how the Virtual DOM works in React.js.
The Virtual DOM is an in-memory representation of the actual DOM. When a state change occurs in a React component, React first updates the Virtual DOM and then compares it with the previous Virtual DOM (using a process called “diffing”). React only applies the differences to the real DOM, optimizing updates.

---

### Q7: React.lazy E SUSPENSE?
In summary, React Suspense is a mechanism for handling asynchronous operations and providing loading states or fallback UIs, while lazy loading is a technique for deferring the loading of resources until they are needed to improve performance