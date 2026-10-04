# OOP, SOLID, DDD, Clean Code

## 🔗 Syllabus

- [Dependency inversion vs dependency injection](#dependency-inversion-vs-dependency-injection)
- [Polymorphism vs Inheritance](#polymorphism-vs-inheritance)
- [Value object vs Entites](#value-object-vs-entites)
- [Subdomains](#subdomains)

## Dependency inversion vs dependency injection

#### 5 years old explanation

You want to build a LEGO car, but you don't have wheels. Instead of making the wheels yourself, your mom gives you the wheels, and you just attach them. Now your car can move!

Now, what if you have different types of wheels? Big wheels, small wheels, even futuristic hover wheels!

Instead of your LEGO car needing a specific type of wheels, you say:
"I don't care what kind of wheels they are, as long as they can spin!"

### Dependency Injection:

A class receives its dependencies from the outside (they are instantiated elsewhere and injected).
Promotes flexibility and testability.

#### Example:

```typescript
// Dependency injection (Logger)
class Logger {
  log(message: string): void {
    console.log("Log:", message);
  }
}
// Class that depends on Logger
class UserService {
  private logger: Logger;

  constructor(logger: Logger) {
    this.logger = logger;
  }

  createUser(name: string): void {
    this.logger.log(`User ${name} created.`);
  }
}
// Dependency is instantiated outside and injected
const logger = new Logger();
const userService = new UserService(logger);
userService.createUser("John Doe");
```

### Dependency Inversion:

High-level modules depend on abstractions (interfaces), not concrete implementations.
Decouples classes by relying on interfaces, not direct instances.
Dependency Inversion adds more boiler code (more code to maintain), and the code becomes harder to navigate

#### Example:

```typescript
// Dependency inversion (Logger)
interface ILogger {
  log(message: string): void;
}

// Low-level module (implementation of ILogger)
class ConsoleLogger implements ILogger {
  log(message: string): void {
    console.log("Console:", message);
  }
}

// High-level module depends on the interface, not a concrete class
class UserService {
  private logger: ILogger;

  constructor(logger: ILogger) {
    this.logger = logger;
  }

  createUser(name: string): void {
    this.logger.log(`User ${name} created.`);
  }
}

// Inject any implementation of ILogger
const logger = new ConsoleLogger();
const userService = new UserService(logger);
userService.createUser("Jane Doe");
```

## Polymorphism vs Inheritance:

### Inheritance
One class (child) gets behavior and properties from another class (parent). Useful for code reuse and shared behavior.
```typescript
// Parent class
class Animal {
  name: string;

  constructor(name: string) {
    this.name = name;
  }

  makeSound(): void {
    console.log("Some generic sound");
  }
}

class Dog extends Animal {
  makeSound(): void {
    console.log("Woof!");
  }
}

const dog = new Dog("Buddy");
dog.makeSound();
```

### Polymorphism
Different classes can share the same method name but behave differently. You can treat them as the same "type" (like Animal), but they act based on their own implementation.
```typescript
class Cat extends Animal {
  makeSound(): void {
    console.log("Meow!");
  }
}

class Duck extends Animal {
  makeSound(): void {
    console.log("Quack!");
  }
}

const animals: Animal[] = [
  new Dog("Buddy"),
  new Cat("Whiskers"),
  new Duck("Daffy"),
];

animals.forEach(animal => animal.makeSound());
/*
  Woof!
  Meow!
  Quack!
*/

```
## Value object vs Entites

#### 5 years old explanation

An entity is like your favorite teddy bear. Even if there are other teddy bears that look exactly the same, you can IDENTIFY yours.
A value object is like a balloon. If you have a red balloon and I give you another red balloon, you don't care if it's the exact same balloon. They are the same thing because their color and size are what matter, not their identity.

#### Value Object:

A Value Object is defined by its attributes rather than a unique identity. They are immutable and used for representing concepts like Money, Address, or Date Range.

#### Example:

```typescript
class CPF {
  private constructor(private readonly value: string) {
    if (!CPF.isValid(value)) {
      throw new Error("Invalid CPF");
    }
  }

  static create(value: string): CPF {
    return new CPF(CPF.clean(value));
  }

  static isValid(value: string): boolean {
    return /^\d{11}$/.test(CPF.clean(value));
  }

  static clean(value: string): string {
    return value.replace(/\D/g, "");
  }

  static format(value: string): string {
    return value.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, "$1.$2.$3-$4");
  }

  equals(other: CPF): boolean {
    return this.value === other.value;
  }
}
```

#### Entities:

An Entity is an object that is defined by a unique identity rather than just its attributes. Entities typically represent domain concepts like Users, Orders, or Products.

#### Example:

```typescript
class User {
  constructor(
    public readonly id: string,
    public name: string,
    public cpf: CPF,
    private isActive: boolean = true
  ) {}

  deactivate(): void {
    this.isActive = false;
  }

  activate(): void {
    this.isActive = true;
  }

  status(): string {
    return this.isActive ? "Active" : "Inactive";
  }

  equals(other: User): boolean {
    return this.id === other.id;
  }
}

// Example usage:
const cpf = CPF.create("123.456.789-09");
const user = new User("1", "John Doe", cpf);

console.log(user.status()); // "Active"
user.deactivate();
console.log(user.status()); // "Inactive"
```
## Subdomains:

#### Core Subdomain (Subdomínio Principal):

- This is the most critical part of the system that provides a unique competitive advantage.

#### example:

in a fraud detection sistem that the evaluation is made manualy buy employees, the work done by the employees is the principal, the software is a suporte subdomain

#### Generic Subdomain:

- A common reusable domain that does not provide competitive advantage
- Can be replaced with a third-party services

#### example:

- Authentication & Authorization in most applications is generic because it's not unique to the business.
- Payment processing (if using a third-party like Stripe) can be considered generic.

#### Supporting Subdomain:

- Important for the business but not the main competitive advantage.
- It supports the core subdomain but does not need to be as optimized.

#### example:

In an e-commerce system, the customer support ticket system is important but does not give a unique market advantage.

