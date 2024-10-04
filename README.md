# Notificationcenter Pub/Sub Library

[![Build Status](https://github.com/geniusrabbit/notificationcenter/workflows/Tests/badge.svg)](https://github.com/geniusrabbit/notificationcenter/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/notificationcenter)](https://goreportcard.com/report/github.com/geniusrabbit/notificationcenter)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/notificationcenter?status.svg)](https://godoc.org/github.com/geniusrabbit/notificationcenter)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/notificationcenter/badge.svg)](https://coveralls.io/github/geniusrabbit/notificationcenter)

> License: Apache 2.0

The **NotificationCenter** library provides a unified interface for publish/subscribe (pub/sub) messaging in Go applications. It simplifies asynchronous communication between services in serverless and microservices architectures by abstracting the complexities of various message brokers.

With NotificationCenter, you can seamlessly integrate different pub/sub backends like Kafka, NATS, Redis, PostgreSQL, and more without altering your application logic. This promotes decoupled architectures, enhancing performance, reliability, and scalability.

## Table of Contents

- [Features](#features)
- [Supported Modules](#supported-modules)
- [Installation](#installation)
- [Usage Examples](#usage-examples)
  - [Import the Package](#import-the-package)
  - [Create a Publisher](#create-a-publisher)
  - [Publish Messages](#publish-messages)
  - [Subscribe to Messages](#subscribe-to-messages)
- [TODO](#todo)
- [License](#license)

## Features

- **Unified Interface**: Interact with multiple pub/sub backends using a consistent API.
- **Easy Integration**: Quickly set up publishers and subscribers with minimal configuration.
- **Backend Flexibility**: Swap out message brokers without changing your application code.
- **Event-Driven Architecture**: Facilitate loosely coupled communication between services.
- **Scalability**: Improve performance and reliability by decoupling application components.

## Supported Modules

- [Kafka](kafka)
- [NATS](nats)
- [NATS Streaming](natsstream)
- [PostgreSQL](pg)
- [Redis](redis)
- [Go Channels](gochan)
- [Time Interval Executor](interval)

## Installation

Install the library using `go get`:

```bash
go get github.com/geniusrabbit/notificationcenter/v2
```

## Usage Examples

Below are basic examples demonstrating how to use NotificationCenter in your Go application.

### Import the Package

```go
import (
  nc "github.com/geniusrabbit/notificationcenter/v2"
)
```

### Create a Publisher

Create a new publisher using one of the supported backends. For example, using **NATS**:

```go
import (
  "github.com/geniusrabbit/notificationcenter/v2/nats"
  "log"
)

// Create a new NATS publisher
eventStream, err := nats.NewPublisher(
  nats.WithNatsURL("nats://hostname:4222"),
)
if err != nil {
  log.Fatal(err)
}

// Register the publisher with NotificationCenter
err = nc.Register("events", eventStream)
if err != nil {
  log.Fatal(err)
}
```

### Publish Messages

You can publish messages using global functions or by obtaining a publisher interface.

**Using Global Functions:**

```go
import (
  "context"
)

// Define your message structure
type Message struct {
  Title string
}

// Publish a message globally
nc.Publish(context.Background(), "events", Message{Title: "Event 1"})
```

**Using Publisher Interface:**

```go
// Get the publisher interface
eventsPublisher := nc.Publisher("events")

// Publish a message
eventsPublisher.Publish(context.Background(), Message{Title: "Event 2"})
```

### Subscribe to Messages

Create a subscriber and register it with NotificationCenter.

```go
import (
  "context"
  "fmt"
  "github.com/geniusrabbit/notificationcenter/v2"
  "github.com/geniusrabbit/notificationcenter/v2/nats"
  "github.com/geniusrabbit/notificationcenter/v2/interval"
  "time"
)

func main() {
  ctx := context.Background()

  // Create a NATS subscriber
  eventsSubscriber := nats.MustNewSubscriber(
    nats.WithTopics("events"),
    nats.WithNatsURL("nats://hostname:4222"),
    nats.WithGroupName("group"),
  )
  nc.Register("events", eventsSubscriber)

  // Optional: Create a time interval subscriber (e.g., for periodic tasks)
  refreshSubscriber := interval.NewSubscriber(5 * time.Minute)
  nc.Register("refresh", refreshSubscriber)

  // Subscribe to the "events" stream
  nc.Subscribe("events", func(msg nc.Message) error {
    // Process the received message
    fmt.Printf("Received message: %v\n", msg.Data())

    // Acknowledge the message if necessary
    return msg.Ack()
  })

  // Subscribe to the "refresh" stream for periodic tasks
  nc.Subscribe("refresh", func(msg nc.Message) error {
    // Perform your periodic task here
    fmt.Println("Performing periodic refresh")
    return msg.Ack()
  })

  // Start listening for messages
  nc.Listen(ctx)
}
```

## TODO

- [ ] Add support for **Amazon SQS**
- [x] Add support for **Redis** queue
- [ ] Add support for **RabbitMQ**
- [ ] Add support for **MySQL notifications**
- [x] Add support for **PostgreSQL notifications**
- [x] ~~Remove deprecated metrics from the queue~~
- [x] Add support for **NATS & NATS Streaming**
- [x] Add support for **Kafka** queue
- [x] Add support for native **Go channels**
- [x] Add support for **Time Interval Execution**

## License

NotificationCenter is licensed under the [Apache 2.0 License](LICENSE).

---

By using NotificationCenter, you can focus on building the core functionality of your application without worrying about the intricacies of different messaging infrastructures. Feel free to contribute to the project or report any issues you encounter.
