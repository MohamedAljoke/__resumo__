# Interview Review

My notes for reviewing before an interview: questions with the answer to say out loud, with **bold keywords** to hit. Code outputs are checked by running them.

## Structure

`languages/` has one folder per programming language. `general/` holds topics that don't depend on a language.

- **[languages/](./languages/README.md)**
  - **[go/](./languages/go/README.md)**
    - [Fundamentals](./languages/go/fundamentals.md): copies, slices, maps, range, strings, receivers, interfaces & nil, errors
    - [Concurrency & Channels](./languages/go/concurrency.md): scheduler, races, memory model, channels, select, context, sync, patterns
  - **[javascript/](./languages/javascript/README.md)**
    - [JavaScript / TypeScript](./languages/javascript/javascript-typescript.md)
    - [Node.js](./languages/javascript/nodejs.md)
    - [React](./languages/javascript/react.md)
- **[general/](./general/README.md)**
  - [System Design](./general/system-design.md)
  - [OOP / SOLID / DDD / Clean Code](./general/oop-solid-ddd-clean-code.md)
  - [SQL](./general/sql.md)
- [Work Experiences](./memory/work-experiences.md): real problems and incidents from work, to use as behavioural answers
- `images/`: diagrams used in the notes

## Adding notes

- A new language gets a folder under `languages/`, with a `README.md` index. Topics that don't depend on a language go in `general/`.
- Each file: syllabus at the top, `**Q: ...**` followed by a short answer, and a 30-second summary at the end.

## License

Personal learning repository. Feel free to use it as a reference, but it may be incomplete or contain mistakes.
