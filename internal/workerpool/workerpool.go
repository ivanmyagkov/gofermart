package workerpool

import (
	"context"
	"log"
	"time"

	"ivanmyagkov/gofermart/internal/client"
)

type OutputWorker struct {
	ch     chan string
	done   chan struct{}
	client *client.AccrualClient
	ctx    context.Context
	ticker *time.Ticker
}

func NewWorker(ch chan string, client *client.AccrualClient, ctx context.Context) *OutputWorker {
	ticker := time.NewTicker(10 * time.Second)
	return &OutputWorker{
		ch:     ch,
		ctx:    ctx,
		client: client,
		ticker: ticker,
	}
}

func (w *OutputWorker) Do() error {
	for {
		select {
		case <-w.ctx.Done():
			w.ticker.Stop()
			return nil
		case order := <-w.ch:
			wait, err := w.client.SentOrder(order)
			if err != nil {
				return err
			}
			log.Println(wait)
			if wait != 0 {
				time.Sleep(time.Duration(wait) * time.Second)
			}
		case <-w.ticker.C:
			if len(w.ch) == 0 {
				break
			}
			wait, err := w.client.SentOrder(<-w.ch)
			if err != nil {
				return err
			}
			if wait != 0 {
				time.Sleep(time.Duration(wait) * time.Second)
			}

		}
	}
}
