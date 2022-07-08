package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/interfaces"
)

type AccrualClient struct {
	client  *http.Client
	db      interfaces.DB
	address string
	qu      chan string
}

func NewAccrualClient(address string, db interfaces.DB, qu chan string) *AccrualClient {
	return &AccrualClient{
		client:  &http.Client{},
		address: address,
		db:      db,
		qu:      qu,
	}
}

func (c *AccrualClient) SentOrder(order string) (int, error) {
	url := fmt.Sprint(c.address, "/api/orders/", order)

	//resp, err := c.client.R().SetContext(ctx).SetPathParams(map[string]string{"orderNumber": task.NumOrder}).Get(c.accrualAddress + "/api/orders/{orderNumber}")
	resp, err := c.client.Get(url)
	if err != nil {
		c.qu <- order
		return 0, err
	}
	defer resp.Body.Close()
	var accrual dto.AccrualResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.qu <- order
		return 0, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &accrual)
		if err != nil {
			c.qu <- order
			return 0, err
		}
		if accrual.OrderStatus == dto.StatusProcessed || accrual.OrderStatus == dto.StatusInvalid {
			if err = c.db.UpdateAccrualOrder(accrual); err != nil {
				c.qu <- order
				return 0, err
			}
		}

		if accrual.OrderStatus == dto.StatusProcessing || accrual.OrderStatus == dto.StatusRegistered {
			if accrual.OrderStatus == dto.StatusProcessing {
				if err = c.db.UpdateAccrualOrder(accrual); err != nil {
					c.qu <- order
					return 0, err
				}
			}
			c.qu <- order
		}

	case http.StatusTooManyRequests:
		wait, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return wait, nil
	default:
		c.qu <- order

	}
	return 0, nil
}
