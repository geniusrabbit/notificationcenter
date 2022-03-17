package gochan

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	nc "github.com/geniusrabbit/notificationcenter/v2"
)

func jsonEncoder(msg any, wr io.Writer) error {
	enc := json.NewEncoder(wr)
	enc.SetEscapeHTML(false)
	return enc.Encode(msg)
}

type encoder func(msg any, wr io.Writer) error

// Publisher writer for GO proxy server
type Publisher struct {
	proxy *Proxy
}

// Publish one or more messages to the pub-service
func (p Publisher) Publish(ctx context.Context, messages ...any) error {
	for _, msg := range messages {
		if err := p.proxy.write(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// Proxy object with native implementation of message pools in GO
type Proxy struct {
	nc.ModelSubscriber
	encoder encoder
	pool    chan message
}

// New proxy object
func New(poolSize int, encoder ...encoder) *Proxy {
	proxy := &Proxy{pool: make(chan message, poolSize)}
	if len(encoder) > 0 && encoder[0] != nil {
		proxy.encoder = encoder[0]
	} else {
		proxy.encoder = jsonEncoder
	}
	return proxy
}

// Publisher writer accessor
func (p *Proxy) Publisher() Publisher {
	return Publisher{proxy: p}
}

func (p *Proxy) write(ctx context.Context, msg any) error {
	buff := &bytes.Buffer{}
	if err := p.encoder(msg, buff); err != nil {
		return err
	}
	p.pool <- message{ctx: ctx, data: buff.Bytes()}
	return nil
}

// Listen starts processing queue
func (p *Proxy) Listen(ctx context.Context) error {
	for msg := range p.pool {
		if err := p.ProcessMessage(msg); err != nil {
			return err
		}
	}
	return nil
}

// Close current proxy pool
func (p *Proxy) Close() error {
	close(p.pool)
	return nil
}

var _ nc.Subscriber = (*Proxy)(nil)
