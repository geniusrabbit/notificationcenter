# Notificationcenter pub/sub library

[![Build Status](https://github.com/geniusrabbit/notificationcenter/workflows/run%20tests/badge.svg)](https://github.com/geniusrabbit/notificationcenter/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/notificationcenter)](https://goreportcard.com/report/github.com/geniusrabbit/notificationcenter)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/notificationcenter?status.svg)](https://godoc.org/github.com/geniusrabbit/notificationcenter)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/notificationcenter/badge.svg)](https://coveralls.io/github/geniusrabbit/notificationcenter)

> License Apache 2.0

Publish/subscribe messaging, or pub/sub messaging, is a form of asynchronous
service-to-service communication used in serverless and microservices architectures.
In a pub/sub model, any message published to a topic is immediately received by all
of the subscribers to the topic. Pub/sub messaging can be used to enable event-driven
architectures, or to decouple applications in order to increase performance,
reliability and scalability.

Library provides basic primitives to use different queue implementations behind,
simplify writing pub/sub-services.

- [Using examples](#Using-examples)
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
  nc "github.com/geniusrabbit/notificationcenter"
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
  nc "github.com/geniusrabbit/notificationcenter"
  "github.com/geniusrabbit/notificationcenter/nats"
)

func main() {
  ctx := context.Background()
  events := nats.MustNewSubscriber(nats.WithTopics("events"),
    nats.WithNatsURL("nats://connection"), nats.WithGroupName(`group`))
  nc.Register("events", events)
  nc.Register("refresh", interval.NewSubscriber(time.Minute * 5))

  // Add new receiver to process the stream "events"
  nc.Subscribe("events", nc.FuncReceiver(ctx, func(msg nc.Message) error {
    fmt.Printf("%v\n", msg.Data())
    return nil
  }))

  // Add new time interval receiver to refresh the data every 5 minutes
  nc.Subscribe("refresh", nc.FuncReceiver(ctx, func(msg nc.Message) error {
    return db.Reload()
  }))

  // Run subscriber listeners
  nc.Listen(ctx)
}
```

## TODO

* [ ] Add support Amazon SQS queue
* [X] Add support Redis queue
* [ ] Add support RabbitMQ queue
* [ ] Add support MySQL notifications queue
* [X] Add support PostgreSQL notifications queue
* [X] Remove metrics from the queue (DEPRECATED)
* [X] Add support NATS & NATS stream
* [X] Add support kafka queue
* [X] Add support native GO chanels
* [X] Add support native GO time interval
