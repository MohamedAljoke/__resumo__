# JavaScript / TypeScript

## 🔗 Syllabus

- [What is sourceMap in typescript options](#what-is-sourcemap-in-typescript-options)
- [Ecmascript vs javascript](#ecmascript-vs-javascript)
- [What does typescript add that is not in javascript](#what-does-typescript-add-that-is-not-in-javascript)
- [Arrow function vs normal function](#arrow-function-vs-normal-function)
- [Closures](#closures)
- [Generators](#generators)
- [Abstract Class vs Interface](#abstract-class-vs-interface)

## External links
- [Closures - Lydia Hallie](https://youtu.be/6Ixyltr8_R0?si=YyHWy_TFFGWXqkpn)
- [Event loop - Boot dev](https://youtu.be/WNrHrwm1wkU?si=ozwVbhNY3UXRrWLo)

## What is sourceMap in typescript options:

- The sourceMap option in the tsconfig.json file enables the generation of .map files, which allow debuggers (like those in VS Code or Chrome DevTools) to map the transpiled JavaScript code back to the original TypeScript source. in this way it will show the typescript file not the compiled one that gave the error

```json
{
  "sourceMap": true /* Generates source map files (.map) that allow debugging of TypeScript code in browsers by mapping compiled JavaScript back to original TypeScript source. This helps with debugging by showing original TypeScript lines in dev tools */
}
```
## Ecmascript vs javascript

- ECMAScript (ES) is the standardized specification for scripting languages, maintained by ECMA International (TC39 committee).
  like Logical assignment operators (&&=, ||=, ??= they are in ecmascript 2021.
  </br>
  If you use like classes with old ecma it will compile to functions, but with new the compiled code can be a class
- JavaScript (JS) is an implementation of ECMAScript with additional features: like WEB API DOM, FETCH, NODE

## What does typescript add that is not in javascript

### Interfaces & Enums:

in javascript you would need to use objects, interface is a typescript implementation COMPILED then to objcts in javascript

```typescript
interface User {
  name: string;
  age: number;
  isAdmin?: boolean;
}
```

```typescript
enum Status {
  Success = 200,
  NotFound = 404,
}
```

### Access Modifiers (public, private, protected):

JavaScript does not have built-in private properties (before ES2020 #privateField).

```typescript
class Person {
  private ssn: string;
  public name: string;

  constructor(name: string, ssn: string) {
    this.name = name;
    this.ssn = ssn; // ❌ Not accessible outside
  }
}

const user = new Person("Alice", "123-45-6789");
```

## Arrow function vs normal function

Arrow functions do not have their own this. They inherit this from the surrounding scope.

### example:

- greetNormal is a regular function attached to user, Since it's a normal function, when called as user.greetNormal(), this correctly refers to the user object.
- fn1 is a normal function inside greetNormal (another function not an object), Since it's a standalone function call (fn1();), this defaults to the global object (window in browsers or undefined in strict mode)
- fn2 is an arrow function inside greetNormal (another function not an object),Arrow functions do not have their own this, Instead, they inherit this from the enclosing scope, which in this case is greetNormal. Since greetNormal's this correctly refers to user, fn2 also uses user's this, so it prints "Alice".
- greetArrow is an arrow function in the root object. Since greetArrow is an arrow function, it inherits this from the outer lexical scope However, greetArrow is defined at the object level, not inside another function.
  Its outer scope is the global scope (window in browsers, undefined in strict mode), so this.name is undefined.

- ps: Objects in JavaScript do not create their own scope.

```typescript
const user = {
  name: "Alice",
  greetNormal: function () {
    console.log(`Hello root function, ${this.name}`);
    const fn1 = function () {
      console.log(`Hello function inside function: `, this.name);
    };
    fn1();
    const fn2 = () => {
      console.log(`Hello arrow function inside function: `, this.name);
    };
    fn2();
  },
  greetArrow: () => {
    console.log(`Hello from root arrow function, ${this.name}`);
  },
};

user.greetNormal();
user.greetArrow();
```

## Closures

- A closure is a function that remembers the variables from its outer scope, even after the outer function has finished executing.
- All functions in javascript are closures, they have a static scope chain that is defined on the moment of its creation
- [Watch the video](https://youtu.be/6Ixyltr8_R0?si=YyHWy_TFFGWXqkpn)

### example:

```typescript
function outerFunction(outerVariable) {
  const insideVar = "insideVar";
  return function innerFunction(innerVariable) {
    console.log(
      `Outer: ${outerVariable}, Inner: ${innerVariable}, insideVar: ${insideVar}`
    );
  };
}

const newFunction = outerFunction("Hello");
const insideVar = "changed"; //even tho we "change the insideVar" it does not change in the function cuz it uses the var on the instance it was created
newFunction("World"); // Output: Outer: Hello, Inner: World, insideVar: insideVar
```
## generators

JavaScript generators are special functions that allow you to pause and resume execution. They are defined using the function* syntax and use the yield keyword to return values incrementally.

### examples:

```typescript
function* fetchData() {
  console.log("Fetching...");
  yield new Promise((resolve) =>
    setTimeout(() => resolve("Data Loaded"), 2000)
  );
}

const dataGen = fetchData();
dataGen.next().value.then(console.log); // After 2 sec -> "Data Loaded"
```

```typescript
const express = require("express");
const app = express();

function* dataGenerator() {
  for (let i = 1; i <= 10; i++) {
    yield `Data chunk ${i}\n`;
  }
}

app.get("/stream", async (req, res) => {
  res.setHeader("Content-Type", "text/plain");
  res.setHeader("Transfer-Encoding", "chunked");

  const generator = dataGenerator();
  for (const chunk of generator) {
    res.write(chunk);
    await new Promise((resolve) => setTimeout(resolve, 1000)); // Simulate delay
  }

  res.end();
});

app.listen(3000, () => console.log("Server running on port 3000"));
```

## Abstract Class vs Interface

- **Abstract Classes**:

  - Can have both implemented and abstract (unimplemented) methods.
  - Can have properties with logic.
  - Supports access modifiers (`public`, `protected`, `private`).
  - Used when multiple related classes share common behavior.

- **Interfaces**:
  - Only define method signatures and properties, no implementations.
  - Cannot have access modifiers; all properties and methods are implicitly public.
  - Supports multiple inheritance (a class can implement multiple interfaces).
  - Used to enforce a structure without defining behavior.

### Example

```typescript
abstract class Animal {
  constructor(public name: string) {}

  abstract makeSound(): void; // Must be implemented by subclasses

  move(): void {
    console.log(`${this.name} is moving`);
  }
}

class Dog extends Animal {
  makeSound(): void {
    console.log("Bark");
  }
}

interface Flyable {
  fly(): void;
}

class Bird extends Animal implements Flyable {
  makeSound(): void {
    console.log("Chirp");
  }

  fly(): void {
    console.log(`${this.name} is flying`);
  }
}
```


