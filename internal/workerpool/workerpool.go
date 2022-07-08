package workerpool

import (
	"context"
	"time"

	"ivanmyagkov/gofermart/internal/client"
)

type OutputWorker struct {
	ch     chan string
	client *client.AccrualClient
	ctx    context.Context
}

func NewWorker(ch chan string, client *client.AccrualClient, ctx context.Context) *OutputWorker {
	return &OutputWorker{
		ch:     ch,
		ctx:    ctx,
		client: client,
	}
}

func (w *OutputWorker) Do() error {
	for {
		select {
		case <-w.ctx.Done():
			return nil
		case order := <-w.ch:
			wait, err := w.client.SentOrder(order)
			if err != nil {
				return err
			}
			if wait != 0 {
				time.Sleep(time.Duration(wait) * time.Second)
			}
		}
	}
}
