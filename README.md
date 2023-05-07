# 2023 RISE - Distributed Systems Assignment

## Prerequisites

* Go: v1.20

## Getting Started

* Example public and private keys
  - Public key: `04d7f44979ac5f7738bd54161f7a2ccbce353c1450c2507d8ddcc0a2f5bbcb8e`
  - Private key: `78fe12de0e7a943e0c2df19bf3878bf137bf11db62e981396e983e506caa300c`
* Copy `.env.example` to `.env` and fill in the values
* Start the example client by: `go run example_client.go`

## Questions

### Phase 1

1. What are some of the challenges you faced while working on Phase 1?

➡️ The biggest challenge to me would be getting familiar with golang. I have only some experiences on developing in golang. Though I use the package, go-nostr, founded on GitHub to solve the problem, I spent some time on understanding their example code and adjusting the code to meet the requirements of the assignment.

2. What kind of failures do you expect to a project such as DISTRISE to encounter?

➡️ The failures might come from the number of WebSocket concurrent connections. Compared to traditional HTTP requests, WebSocket connections consume more resources, which can lead to performance degradation if the number of concurrent connections is too high. In such cases, traditional server architectures may not be able to handle the load and could crash, requiring specialized architectures to mitigate the issue.

## References

* [nbd-wtf/go-nostr](https://github.com/nbd-wtf/go-nostr)
* [NIP-01](https://github.com/nostr-protocol/nips/blob/master/01.md)
