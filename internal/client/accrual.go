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
	qu      chan dto.AccrualResponse
}

func NewAccrualClient(address string, qu chan dto.AccrualResponse) *AccrualClient {
	return &AccrualClient{
		client:  &http.Client{},
		address: address,
		qu:      qu,
	}
}

func (c *AccrualClient) SentOrder(order dto.AccrualResponse) (dto.AccrualResponse, int, error) {

	url := fmt.Sprint(c.address, "/api/orders/", order.NumOrder)

	resp, err := c.client.Get(url)
	if err != nil {
		return order, 0, err
	}
	defer resp.Body.Close()
	var accrual dto.AccrualResponse
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
	case http.StatusTooManyRequests:
		wait, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return order, wait, nil
	}
	return order, 0, nil
}
