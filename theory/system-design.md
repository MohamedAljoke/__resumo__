# System Design

- [Synchronous vs Asynchronous](#synchronous-vs-asynchronous)
- [Orchestrated vs Coordination](#orchestrated-vs-coordination)
- [Elasticity vs Escalability](#elasticity-vs-escalability)
- [Static coupling vs Dynamic coupling](#static-coupling-vs-dynamic-coupling)
- [CAP Theorem](#cap-theorem-consistency-availability-partition-tolerance)
- [ACID vs BASE](#acid-vs-base-cap-theorem-related)
- [Blue/Green Deployment vs Canary Deployment](#bluegreen-deployment-vs-canary-deployment)
- [Monolithic Event-Driven and Microservices Event-Driven](#monolithic-event-driven-and-microservices-event-driven)
- [Reliability, Fault Tolerance, and Redundancy](#reliability-fault-tolerance-and-redundancy)
- [Throughput and Latency](#throughput-and-latency)
- [The Open Graph Protocol](#the-open-graph-protocol-ogp)
- [Compiled vs interpreted languages](#compiled-vs-interpreted-languages)
- [TCP vs UDP](#tcp-vs-udp)
- [Optimistic Locking vs Pessimistic Locking](#optimistic-locking-vs-pessimistic-locking)
- [Deadlocks](#what-are-deadlocks)
- [HTTP vs WebSocket](#http-vs-websocket)
- [Polling vs Long Polling vs Server sent events](#polling-vs-long-polling-vs-server-sent-events)
- [Server Sent events VS WebSockets](#server-sent-events-vs-websockets)
- [REST VS GRAPHQL VS GRPC](#rest-vs-graphql-vs-grpc)


System Design is all about planning how different parts of a system work together to handle scalability, performance, reliability, and maintainability.

## Synchronous vs Asynchronous

#### 5 years old explanation

Synchronous is like waiting in line for ice cream. You have to wait for the person in front of you to finish before you can get yours.

Asynchronous is like ordering a pizza. You call and place your order, but while waiting for the pizza to arrive, you can play, watch TV, or do something else. When it's ready, they bring it to you.

### Synchronous:

Synchronous operations happen one after another, in a blocking manner. Each task must complete before the next one begins. This is common in sequential programming where each step depends on the previous step finishing.

### Asynchronous:

Asynchronous operations allow tasks to happen independently and not block the execution of the program. Instead of waiting, the program can continue running other tasks while waiting for the operation to complete.

### Trade-offs Between Synchronous and Asynchronous

| **Factor**          | **Synchronous**                                               | **Asynchronous**                                                                         |
| ------------------- | ------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| **Execution Order** | Tasks run sequentially, one after another.                    | Tasks can run independently without waiting.                                             |
| **Performance**     | Slower for tasks that involve waiting (e.g., I/O operations). | More efficient for handling multiple tasks concurrently.                                 |
| **Complexity**      | Easier to understand and debug.                               | More complex due to callbacks, promises, and concurrency issues.                         |
| **Blocking**        | Blocks execution until a task is completed.                   | Does not block; other tasks can continue running.                                        |
| **Error Handling**  | Errors are handled sequentially, making debugging simpler.    | Requires handling multiple possible execution paths (e.g., promise chains, async/await). |
| **Use Case**        | Best for tasks that must happen in order, like calculations.  | Best for tasks like API calls, file I/O, or user interactions.                           |
| **Resource Usage**  | Can lead to inefficiency if waiting on slow operations.       | Uses resources efficiently by not blocking the main thread.                              |
| **Scalability**     | Harder to scale due to blocking nature.                       | More scalable, especially for high-load applications.                                    |

## Orchestrated vs Coordination

#### 5 years old explanation

Orchestrated is like a conductor leading an orchestra. The conductor tells each musician when to play so that everything happens in a structured way.

Coordination is like a group of friends playing music together without a conductor. They listen to each other and adjust as they go, but no one is in complete control.

### Orchestrated:

In an orchestrated system, there is a central authority that controls and directs the flow of tasks. Each component follows predefined instructions, much like musicians following a conductor. This approach ensures a structured, predictable sequence of actions.

- example: aws step function

### Coordination:

In a coordinated system, components work together without a central authority dictating their actions. Each part is aware of others and adjusts dynamically based on communication and shared state.

- Orchestration a central service controls the work flow, it easy to track workflows but adds a single point of control (and potential failure).
- Choreography Services react to events and determine their next step independently (like event-driven architectures) scales better and reduces dependencies but can become chaotic without clear event contracts.

## Elasticity vs Escalability

- Escalability how much the service can handle simultaneos users.
- Elasticity how much it can handle an explosion of users
  ![alt text](images/scalability_elasticity.png)

## Static coupling vs Dynamic coupling

static coupling describes how services are wired together _compilation phase_, whereas dynamic coupling
describes how services call one another at runtime.
it is important to know that a database an SQS queue are static coupling components, when we call the SQS that is a dynamic coupling

## CAP Theorem (Consistency, Availability, Partition Tolerance)

- **Consistency** means all nodes in a distributed system see the same data at the same time.
- **Availability** means every request gets a response, even if some of the data might not be up-to-date.
- **Partition Tolerance** ensures the system continues to function even if network partitions occur between nodes.

## ACID vs BASE (CAP Theorem related)

- **ACID** refers to the properties (Atomicity, Consistency, Isolation, Durability) of traditional databases ensuring reliable transactions.
- **BASE** (Basically Available, Soft state, Eventually consistent) is a model used in NoSQL systems where availability is prioritized over consistency.

## Blue/Green Deployment vs Canary Deployment

### **Blue/Green Deployment**

involves running two environments (blue and green) and switching traffic between them to minimize downtime during updates.

- When minimal downtime is critical.
- Major application updates requiring a quick rollback.
- Ensuring a stable environment before switching traffic.

### **Canary Deployment**

gradually rolls out changes to a small subset of users before deploying to everyone, allowing for early detection of issues.

- When testing updates with real users but in a controlled manner.
- When monitoring performance and errors in a production-like environment.
- When reducing risk of major failures with incremental rollout.

## Monolithic Event-Driven and Microservices Event-Driven:

- Monolithic Event-Driven: Inside a monolith, you can implement event-driven behavior to reduce tight coupling, making your app more scalable, though it's still running as a single process.
- Microservices Event-Driven: In microservices, each service is independent and communicates through events, enabling decoupling and better scalability across machines or containers.

## Reliability, Fault Tolerance, and Redundancy:

- Reliability: the system's ability to perform its intended functions without failure or errors over a specified period of time.
- Fault Tolerance: refers to how well the system can detect and heal itself from a problem, i.e. disable a function, revert to a different mode, switch to a different server.
- Redundancy: This redundancy is provided by our backup server which essentially "shadows" the contents of a server. We don't need this server, but it only comes into play if our primary server fails. Only by having this redundancy, are we able to have fault tolerance. In this case, the second server is simply a backup. But, what if we had two servers that were both active? This would be called active-active redundancy.

## Throughput and Latency:

- Throughput: refers to the amount of data or operations we can handle over some period of time. The throughput of a client making requests to a server would be measured through the number of requests per second, Database(queries/second),Network(bytes/second).
- Latency: refers to the delay between the client making the request and the server responding to that request.

## The Open Graph Protocol (OGP):

By embedding specific meta tags in the <head> of an HTML document, you can control how your content appears when shared on platforms like Facebook, LinkedIn, Twitter (with some differences), and others.

#### Example:

```html
<meta property="og:title" content="My Awesome Website" />
<meta
  property="og:description"
  content="This is a description of my awesome website."
/>
<meta property="og:image" content="https://example.com/image.jpg" />
<meta property="og:url" content="https://example.com" />
<meta property="og:type" content="website" />
<meta property="og:site_name" content="MySite" />
```

## Compiled vs interpreted languages:

#### 5 years old explanation

Imagine you're reading a storybook. In one case, the story is written in a special way that the book is turned into a picture book, and you can read it right away—like magic! This is similar to a compiled language. Before you can play with it, someone (a compiler) reads the whole book and turns it into something you can use directly.

In another case, you have a person reading the story to you one page at a time. Each time you want to hear a part of the story, the reader (the interpreter) reads it aloud to you, page by page, without turning the whole book into something you can read on your own. This is like an interpreted language. It reads and understands the story one part at a time while you play with it!

## TCP vs UDP

#### 5 years old explanation

Imagine you are sending letters to your friend:

- TCP (Like sending a letter with tracking)
  When you send a letter, you make sure your friend gets it. If the letter gets lost, you send it again. Your friend also tells you, "Hey, I got your letter!" so you know it arrived safely. This is like TCP—it makes sure all messages are delivered in the right order and without mistakes.
- UDP (Like throwing paper airplanes)
  Instead of mailing a letter, you throw paper airplanes with messages. Some might land perfectly, but others might get lost or crumpled. You don't wait for your friend to say, "I got it!" You just keep throwing. This is like UDP—it's super fast, but messages might get lost.

### TCP

- It sends the data _IN ORDER_
- Ensures delivery, if a packet fails it resends it
- slower
- HTTP, WS, SMTP

### UDP

- Not Reliabile
- Not in order
- ALOT faster
- streaming data, gaming

## Optimistic Locking vs Pessimistic Locking

### Pessimistic Locking:

- locks the data (like database row) whe you read it, no one else can read or update that row
- you can imagin that this makes cost, because everyone else will need to wait till you finish
- Safer when high concurrency is expected
- Avoids race conditions.
- Risk of deadlocks. Can reduce performance due to waiting and blocking.

### Optimistic Locking:

- when you read a row you will have like a version, when you are going to save the update it checks if the version is the same,
  if it is the same you can, if another query raced and changed the data before you the version changed and your query needs to redo.
- Reads data without locking.
- On save, checks if data has changed.
- If changed, throws error or retries.
- Requires conflict handling logic.

## what are deadlocks:

- A deadlock is a situation where two or more operations are stuck waiting for each other, and none can proceed. It's like a circular wait where everyone is holding something the other one needs.
- like if you have a database and you have an operation that does something, but is waiting for the other one doing the same operation to finish, but this other is waiting for your operation to finish.
- it can happen in integration tests, if you are not running in band, and tests are using same db and ur tests are waiting for the other one to finish they might end up waiting for each other

## HTTP vs WebSocket:

### HTTP:
HTTP (Like knocking on a door every time you want to talk)
Every time you want to say something, you knock on your friend's door, wait for them to open it, say your thing, and then leave. Then if you want to say something else, you do the whole thing again. That's HTTP — it opens a connection every time it sends something, then closes it.
- One-time connection per request
- Client talks first, server replies (no surprises!)
- Slower for real-time stuff
- Great for: loading web pages, APIs, forms

### WebSocket (Like a walkie-talkie):
You and your friend both have walkie-talkies. You connect once and can keep talking back and forth as much as you want without knocking again. That's WebSocket — it's like having a constant open line to chat freely!
- Always-on connection after the first handshake
- Both sides can talk any time
- Super fast for live updates
- Great for: chat apps, games, live data dashboards

## Polling vs Long Polling vs Server sent events:

### Polling:
Polling (Like asking "Did you write me a letter?" every minute)
You keep asking the mailman every minute: "Is there a letter for me?" Even if there's nothing new, you still ask over and over. That's polling — your app keeps checking the server even if there's nothing to say.
- Client asks the server repeatedly
- Easy to implement
- Can waste resources
- Good for: simple checks, legacy systems
### Long Polling :
Long Polling (Like waiting at the mailbox until a letter comes)
You ask the mailman, "Is there a letter?" and if there isn't one, he says, "Wait here, I'll tell you when it comes." When it finally arrives, you go home and then come back to wait again. That's long polling — slower, but less asking.

- Client waits until server has something to reply
- Less traffic than regular polling
- Still needs to reconnect every time
- Good for: chat apps before WebSocket

### Server sent events:
Server-Sent Events (Like a radio where only the server talks)
You have a radio, and your friend is a radio DJ. They can talk to you any time, but you can't talk back. That's Server-Sent Events — one-way messages from the server to the client.
- Server pushes updates to client
- One-way (server → client)
- Only works over HTTP
- Good for: notifications, live feeds


## Server Sent events VS WebSockets:
SSE uses plain HTTP, so it's easier to set up and works naturally with things like:
Proxies,Firewalls,Load balancers (when you use aws api gateway for example for websocket you need to use another especific for it) 
- SSE has built-in auto-reconnect:
- With WebSocket, you need to handle reconnection and state manually.
- SSE is perfect when only the server needs to push data, like: (lighter)


## REST VS GRAPHQL VS GRPC

### REST
- Uses HTTP and typically JSON format.
- Each resource has its own endpoint (e.g., /users, /posts).
- Can lead to over-fetching or under-fetching data.
- Easy to use and widely supported.
- Great for simple CRUD APIs.

### GraphQL
- Uses a single HTTP endpoint.
- Allows clients to request exactly the data they need.
- Uses a strongly typed schema.
- Reduces over-fetching and under-fetching.
- Excellent for frontend-heavy applications.

### gRPC
- Uses HTTP/2 and Protocol Buffers (Protobuf) for data exchange.
- Defines services and methods using .proto files.
- High performance and low latency.
- Strong typing and efficient for internal microservice communication.
- More setup required compared to REST or GraphQL.


