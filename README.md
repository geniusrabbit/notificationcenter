# Notificationcenter

[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/notificationcenter)](https://goreportcard.com/report/github.com/geniusrabbit/notificationcenter)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/notificationcenter?status.svg)](https://godoc.org/github.com/geniusrabbit/notificationcenter)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/notificationcenter/badge.svg)](https://coveralls.io/github/geniusrabbit/notificationcenter)

> License Apache 2.0

The union eventstream interface.

### Create new stream processor

```go
// Create new stream processor
eventStream, err = nats.NewLogger([]string{"event"}, "nats://hostname:4222")
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
events := notificationcenter.LoggerByName("events")
events.Send(message{title: "event 2"})
```
