package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/interfaces"
)

const StatusRegistered string = "REGISTERED"
const StatusInvalid string = "INVALID"
const StatusProcessing string = "PROCESSING"
const StatusProcessed string = "PROCESSED"

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
	log.Println("hiii")
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
	switch resp.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &accrual)
		if err != nil {
			c.qu <- order
			return 0, err
		}
		if accrual.OrderStatus == StatusProcessed || accrual.OrderStatus == StatusInvalid {
			if err = c.db.UpdateAccrualOrder(accrual); err != nil {
				c.qu <- order
				return 0, err
			}
		}

		if accrual.OrderStatus == StatusProcessing {
			accrual.OrderStatus = dto.StatusNew
			if err = c.db.UpdateAccrualOrder(accrual); err != nil {
				c.qu <- order
				return 0, err
			}
		}
		if err = c.db.UpdateAccrualOrder(accrual); err != nil {
			c.qu <- order
			return 0, err
		}

	case http.StatusTooManyRequests:
		wait, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return wait, nil
	default:
		c.qu <- order

	}
	return 0, nil
}
