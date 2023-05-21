# 2023 RISE - Distributed Systems Assignment

## Prerequisites

- Go: v1.20.4

## Getting Started

### Phase 1

- Example public and private keys
  - Public key: `04d7f44979ac5f7738bd54161f7a2ccbce353c1450c2507d8ddcc0a2f5bbcb8e`
  - Private key: `78fe12de0e7a943e0c2df19bf3878bf137bf11db62e981396e983e506caa300c`
- Copy `.env.example` to `.env` and fill in the values
- Start the example client by: `go run example_client.go`

### Phase 2

- Copy `.env.example` to `.env` and fill in the values
- Start the cockroach DB by: `docker compose up -d`
- Start the relay by: `make run`
- The relay will be running on "ws://localhost:8080"
- To view the built-in CockroachDB admin UI, visit "http://localhost:8081"

### Phase 3

- Copy `.env.example` to `.env` and fill in the values
- Start the RabbitMQ and cockroach DB by: `docker compose up -d`
- Start the aggregator by: `make aggregate`
- To see the events stored in DB: `make view`

## Questions

### Phase 1

1. What are some of the challenges you faced while working on Phase 1?

➡️ The biggest challenge to me would be getting familiar with golang. I have only some experiences on developing in golang. Though I use the package, go-nostr, founded on GitHub to solve the problem, I spent some time on understanding their example code and adjusting the code to meet the requirements of the assignment.

2. What kind of failures do you expect to a project such as DISTRISE to encounter?

➡️ The failures might come from the number of WebSocket concurrent connections. Compared to traditional HTTP requests, WebSocket connections consume more resources, which can lead to performance degradation if the number of concurrent connections is too high. In such cases, traditional server architectures may not be able to handle the load and could crash, requiring specialized architectures to mitigate the issue.

### Phase 2

1. Why did you choose this database?
   ➡️ I choose Cockroach DB simply because it is a distributed SQL database and I want to try it out. 🙂

2. If the number of events to be stored will be huge, what would you do to scale the database?
   ➡️ From the document of Cockroach Labs, I might consider using load balancing. As the number of events increases, a multi-node CockroachDB cluster can distribute the traffic from clients. Ref: [Deploy CockroachDB On-Premises](https://www.cockroachlabs.com/docs/stable/deploy-cockroachdb-on-premises.html#step-6-set-up-load-balancing)

Questions:

1. Is there any free Cockroach DB GUI client for Mac or Windows? Seems like the Cockroach DB admin UI cannot query and show the data. I have also tried using PgAdmin to connect to the Cockroach DB but failed too.

Self-reflection:

Thanks to lots of online golang & websocket tutorials and sharings from peers save me, golang newbie 😅, a lot of time on this assignment. I underestimated the difficulty of the tech stack I selected to finish this assignment. I need to spend more time understanding the basics of golang, like goroutine, channel, etc. Also, spending time using docker compose to build a Cockroach DB cluster would be a good idea too. It might be a fundamental requirement for future assignments.

### Phase 3

1. Why did you choose this database? Is it the same or different database as the one you used in Phase 2? Why is it the same or a different one?
   ➡️ I'm still using Cockroach DB, same as Phase 2. The reason for the same one simply because I haven't been familiar with it, so I want to see how it works more and try to use it in the following assignments.

2. If the number of events to be stored will be huge, what would you do to scale the database?
   ➡️ I might consider using sharding to scale the database. Shard the database by the event type, and each shard can be stored in a different node.

## References

- [nbd-wtf/go-nostr](https://github.com/nbd-wtf/go-nostr)
- [NIP-01](https://github.com/nostr-protocol/nips/blob/master/01.md)
- [Build a Simple Websocket Chat Server in Go with Gin](https://lwebapp.com/en/post/go-websocket-chat-server)
- [Multi-room Chat Application With WebSockets In Go And Vue.js (Part 2)](https://dev.to/jeroendk/multi-room-chat-application-with-websockets-in-go-and-vue-js-part-2-3la8)
- [Cockroach Labs - Start a single-node cluster](https://www.cockroachlabs.com/docs/stable/start-a-local-cluster-in-docker-linux.html#start-a-single-node-cluster)
- [Build a Go App with CockroachDB and GORM](https://www.cockroachlabs.com/docs/v22.2/build-a-go-app-with-cockroachdb-gorm)
- [RabbitMQ tutorial - "Hello World!"](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
