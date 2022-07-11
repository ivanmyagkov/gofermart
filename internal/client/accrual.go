package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"ivanmyagkov/gofermart/internal/dto"
)

type AccrualClient struct {
	client  *http.Client
	address string
	qu      chan dto.Order
}

func NewAccrualClient(address string, qu chan dto.Order) *AccrualClient {
	return &AccrualClient{
		client:  &http.Client{},
		address: address,
		qu:      qu,
	}
}

func (c *AccrualClient) SentOrder(order dto.Order) (dto.Order, int, error) {

	url := fmt.Sprint(c.address, "/api/orders/", order.Number)

	resp, err := c.client.Get(url)
	if err != nil {
		return order, 0, err
	}
	defer resp.Body.Close()
	var accrual dto.Order
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return order, 0, err
	}
	err = json.Unmarshal(body, &accrual)
	if err != nil {
		return order, 0, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return accrual, 0, nil
		//if accrual.Status == dto.StatusProcessed || accrual.Status == dto.StatusInvalid {
		//	if err = c.db.UpdateAccrualOrder(accrual); err != nil {
		//		c.qu <- accrual
		//		return 0, err
		//	}
		//}
		//
		//if accrual.Status == dto.StatusProcessing || accrual.Status == dto.StatusRegistered {
		//	if accrual.Status == dto.StatusProcessing {
		//		if err = c.db.UpdateAccrualOrder(accrual); err != nil {
		//			c.qu <- accrual
		//			return 0, err
		//		}
		//	}
		//	c.qu <- accrual
		//}

	case http.StatusTooManyRequests:
		wait, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return order, wait, nil
	}
	return order, 0, nil
}
