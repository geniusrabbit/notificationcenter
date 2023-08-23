# Notificationcenter pub/sub library

[![Build Status](https://github.com/geniusrabbit/notificationcenter/workflows/Tests/badge.svg)](https://github.com/geniusrabbit/notificationcenter/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/notificationcenter)](https://goreportcard.com/report/github.com/geniusrabbit/notificationcenter)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/notificationcenter?status.svg)](https://godoc.org/github.com/geniusrabbit/notificationcenter)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/notificationcenter/badge.svg)](https://coveralls.io/github/geniusrabbit/notificationcenter)

> License Apache 2.0

Publish/subscribe messaging, often referred to as pub/sub messaging, serves as a pivotal form of asynchronous communication between services within serverless and microservices architectures. Operating on a pub/sub model, this approach entails the instantaneous transmission of any published message to all subscribers associated with the corresponding topic. The utility of pub/sub messaging extends to enabling event-driven architectures and the seamless decoupling of applications, ultimately yielding improvements in performance, reliability, and scalability.

At its core, this mechanism involves the interaction between publishers, who disseminate messages, and subscribers, who receive and act upon these messages. By employing this model, systems can leverage the power of loosely coupled communication, enhancing the adaptability of individual components within the broader infrastructure.

To streamline the implementation of this messaging paradigm, libraries provide essential foundational elements that facilitate the utilization of various queue implementations. These libraries abstract the complexities of interacting with diverse queuing systems, thereby simplifying the development of pub/sub services. This not only promotes efficient communication between services but also empowers developers to concentrate on the business logic and functionality of their applications without becoming entangled in the intricacies of messaging infrastructures.

- [Using examples](#using-examples)
  - [Create new publisher processor](#create-new-publisher-processor)
  - [Send event by the notification publisher](#send-event-by-the-notification-publisher)
  - [Subscribe by the specific notification publisher](#subscribe-by-the-specific-notification-publisher)
- Modules
  - [Kafka](kafka)
  - [NATS](nats)
  - [NATS Stream](natstream)
  - [PostgreSQL](pg)
  - [Redis](redis)
  - [Golang Chanels implementation](gochan)
  - [Golang time interval executor](interval)
- [TODO](#todo)

## Using examples

Basic examples of usage.

```go
import(
  nc "github.com/geniusrabbit/notificationcenter/v2"
)
```

### Create new publisher processor

```go
// Create new publisher processor
eventStream, err = nats.NewPublisher(nats.WithNatsURL("nats://hostname:4222/group?topics=event"))
if err != nil {
  log.Fatal(err)
}

// Register stream processor
err = nc.Register("events", eventStream)
if err != nil {
  log.Fatal(err)
}
```

### Send event by the notification publisher

```go
// Send by global functions
nc.Publish(context.Background(), "events", message{title: "event 1"})

// Send by logger interface
events := nc.Publisher("events")
events.Publish(context.Background(), message{title: "event 2"})
```

### Subscribe by the specific notification publisher

```go
import (
  nc "github.com/geniusrabbit/notificationcenter/v2"
  "github.com/geniusrabbit/notificationcenter/v2/nats"
)

func main() {
  ctx := context.Background()
  events := nats.MustNewSubscriber(nats.WithTopics("events"),
    nats.WithNatsURL("nats://connection"), nats.WithGroupName(`group`))
  nc.Register("events", events)
  nc.Register("refresh", interval.NewSubscriber(time.Minute * 5))

  // Add new receiver to process the stream "events"
  nc.Subscribe("events", func(msg nc.Message) error {
    fmt.Printf("%v\n", msg.Data())
    return msg.Ack()
  })

  // Add new time interval receiver to refresh the data every 5 minutes
  nc.Subscribe("refresh", func(msg nc.Message) error {
    return db.Reload()
  })

  // Run subscriber listeners
  nc.Listen(ctx)
}
```

## TODO

- [ ] Add support Amazon SQS queue
- [X] Add support Redis queue
- [ ] Add support RabbitMQ queue
- [ ] Add support MySQL notifications queue
- [X] Add support PostgreSQL notifications queue
- [X] ~~Remove metrics from the queue (DEPRECATED)~~
- [X] Add support NATS & NATS stream
- [X] Add support kafka queue
- [X] Add support native GO chanels
- [X] Add support native GO time interval
