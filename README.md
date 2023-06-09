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

### Phase 4

- Copy `.env.example` to `.env` and fill in the values. To set up multiple relays, set the value of `RELAY_URLS` to a comma-separated list of relay URLs.
- Start the RabbitMQ and cockroach DB by: `docker compose up -d`
- Start the aggregator by: `make aggregate`
- To see the events stored in DB: `make view`

### Phase 5

- Copy `.env.example` to `.env` and fill in the values. To set up multiple relays, set the value of `RELAY_URLS` to a comma-separated list of relay URLs.
- Start the relay and aggregator by: `docker compose up -d`

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

### Phase 4

Please provide a short writeup of why you chose a particular Queue or Event Stream system for Phase 4, answering the following questions:

1. Why did you choose this solution?

➡️ I choose RabbitMQ rather than other kinds of Queue, like Kafka, because I think it is the best fit for this assignment. Using Kafka might be a little bit overkill. Also, storing the events in a queue is like a buffer, the data don't need to be stored after being processed, but because of the retention setting of Kafka, the data will be stored for a period of time. I think RabbitMQ is a better choice for this assignment.

2. If the number of events to be stored will be huge, what would you do to scale your chosen solution?

➡️ I will consider multiple steps to scale. First, I will increase the memory size of the RabbitMQ instance and run more consumers. Second, using a cluster of RabbitMQ instances. Third, using a cluster of RabbitMQ instances and sharding the events by the event type. Finally, if the whole event throughput is higher than the database capability, I will consider scaling up the database discussed earlier.

Self-reflection:

I forgot to check the test cases listed on the Notion page, so the aggregator only meets the requirements of the assignment. I will try to fix this in the following days.

### Phase 5

Please provide a short writeup of why you chose these metrics, answering the following questions:

1. Why did you choose these 5 metrics?

➡️ I host my relay on an Ubuntu server on GCP, and I choose Datadog trail for collecting the metrics. Setting up the datadog agent is quite easy on the Ubuntu server, but installing the integration like Nginx and RabbitMQ took me a while. The datadog agent and the integrations already collect lots of metrics, including: 1. the usage of CPU, memory, disk, and network of the host server, 2. the connection status of the nginx, 3. memory usage, queue size, and connection status of the RabbitMQ. I think these metrics are enough for me to monitor the status of the relay and aggregator. These metrics can help me understanding the load of the server and the status of the services running on the server. For example, if the CPU usage is high, I might need to scale up the server or the RabbitMQ instance, or if connections of the Nginx is high, I might need to scale up the relay.

2. What kind of errors or issues do you expect them to help you monitor?

➡️ I expect them to monitor the server and the service health. For example, if the storage of the server is full, it might mean the Cockroach DB is storing to much data of the events, and I might need to scale up the storage of my VM to let the aggregator keeps storing new events.

3. If you had more time, what other tools would you use or metrics would you instrument and why?

➡️ Besides getting metrics from Nginx and RabbitMQ, I also tried to use OpenTelemetry to collect traces of my relay server. Sadly, because of limited time I did not pass those data to the datadog successfully. As a result, I would spend some time digging the usage of the OpenTelemetry. Also, I might consider using Prometheus and Grafana to monitor the metrics. They are also popular & free tools for monitoring and I want to try them out too.

## Datadog Dashboards

- [System - Metrics](https://p.ap1.datadoghq.com/sb/020745ac-0353-11ee-a94b-da7ad0900009-640899023508a05cdb5d8ffcdd515956)
- [NGINX - Metrics](https://p.ap1.datadoghq.com/sb/020745ac-0353-11ee-a94b-da7ad0900009-565bf790ee2917c03d02f2a57f8b63cf)
- [RabbitMQ Overview (OpenMetrics Version)](https://p.ap1.datadoghq.com/sb/020745ac-0353-11ee-a94b-da7ad0900009-90cc081ce68838d2d0664e2a4ef6c170)

## Deployment

Relay and aggregator are deployed on GCP.
Relay URL: `wss://andyhsu10.sefo.fi`

## References

- [nbd-wtf/go-nostr](https://github.com/nbd-wtf/go-nostr)
- [NIP-01](https://github.com/nostr-protocol/nips/blob/master/01.md)
- [Build a Simple Websocket Chat Server in Go with Gin](https://lwebapp.com/en/post/go-websocket-chat-server)
- [Multi-room Chat Application With WebSockets In Go And Vue.js (Part 2)](https://dev.to/jeroendk/multi-room-chat-application-with-websockets-in-go-and-vue-js-part-2-3la8)
- [Cockroach Labs - Start a single-node cluster](https://www.cockroachlabs.com/docs/stable/start-a-local-cluster-in-docker-linux.html#start-a-single-node-cluster)
- [Build a Go App with CockroachDB and GORM](https://www.cockroachlabs.com/docs/v22.2/build-a-go-app-with-cockroachdb-gorm)
- [RabbitMQ tutorial - "Hello World!"](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
- [OpenTelemetry - Getting Started](https://opentelemetry.io/docs/instrumentation/go/getting-started/)
- [OpenTelemetry Tracing API for Go](https://uptrace.dev/opentelemetry/go-tracing.html)
