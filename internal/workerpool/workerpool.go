package workerpool

import (
	"context"
	"time"

	"ivanmyagkov/gofermart/internal/client"
	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/interfaces"
)

type OutputWorker struct {
	ch     chan dto.Order
	db     interfaces.DB
	client *client.AccrualClient
	ctx    context.Context
}

func NewWorker(ch chan dto.Order, client *client.AccrualClient, ctx context.Context, db interfaces.DB) *OutputWorker {
	return &OutputWorker{
		ch:     ch,
		ctx:    ctx,
		client: client,
		db:     db,
	}
}

func (w *OutputWorker) Do() error {
	for {
		select {
		case <-w.ctx.Done():
			return nil
		case order := <-w.ch:
			newOrder, wait, err := w.client.SentOrder(order)
			if err != nil {
				w.ch <- newOrder
				return err
			}
			if wait != 0 {
				w.ch <- newOrder
				time.Sleep(time.Duration(wait) * time.Second)
				return nil
			}

			if newOrder.Status == dto.StatusProcessed || newOrder.Status == dto.StatusInvalid {
				if err = w.db.UpdateAccrualOrder(newOrder); err != nil {
					w.ch <- newOrder
					return err
				}
			}

			if newOrder.Status == dto.StatusProcessing || newOrder.Status == dto.StatusRegistered {
				if newOrder.Status == dto.StatusProcessing {
					if order.Status != newOrder.Status {
						if err = w.db.UpdateAccrualOrder(newOrder); err != nil {
							w.ch <- newOrder
							return err
						}
					}
				}
				w.ch <- newOrder
			}

		}
	}
}
