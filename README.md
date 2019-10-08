# Notificationcenter stream library

[![Build Status](https://travis-ci.org/geniusrabbit/notificationcenter.svg?branch=master)](https://travis-ci.org/geniusrabbit/notificationcenter)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/notificationcenter)](https://goreportcard.com/report/github.com/geniusrabbit/notificationcenter)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/notificationcenter?status.svg)](https://godoc.org/github.com/geniusrabbit/notificationcenter)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/notificationcenter/badge.svg)](https://coveralls.io/github/geniusrabbit/notificationcenter)

> License Apache 2.0

The union eventstream wrapper over nifferent stream implementations.

## Using examples

### Create new stream processor

```go
// Create new stream processor
eventStream, err = nats.NewStream([]string{"event"}, "nats://hostname:4222")
if err != nil {
  log.Fatal(err)
}

// Register stream processor
err = notificationcenter.Register("events", eventStream)
if err != nil {
  log.Fatal(err)
}
```

### Send event to the notification stream

```go
// Send by global functions
notificationcenter.Send("events", message{title: "event 1"})

// Send by logger interface
events := notificationcenter.StreamByName("events")
events.Send(message{title: "event 2"})
```

### Subscribe for the specific notification stream

```go
import (
  nc "github.com/geniusrabbit/notificationcenter"
  "github.com/geniusrabbit/notificationcenter/nats"
)

func main() {
  events := nats.MustNewSubscriber("nats://connection", "group", []string{"events"})
  nc.Register("events", events)

  // Add new handler to process the stream "events"
  nc.Subscribe("events", nc.FuncHandler(func(msg nc.Message) error {
    fmt.Printf("%v\n", msg.Data())
    return nil
  }))

  // Run seubscribers listeners
  nc.Listen()
}
```

## TODO

* [ ] Add support Amazon SQS queue
* [ ] Add support Redis queue
* [ ] Add support RabbitMQ queue
* [ ] Add support MySQL notifications queue
* [X] Add support PostgreSQL notifications queue
* [X] Remove metrics from the queue (DEPRECATED)
* [X] Add support NATS & NATS stream
* [X] Add support kafka queue
