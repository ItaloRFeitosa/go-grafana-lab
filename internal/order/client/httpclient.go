package client

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/italorfeitosa/go-grafana-lab/internal/env"
	"github.com/italorfeitosa/go-grafana-lab/internal/order/model"
	"github.com/italorfeitosa/go-grafana-lab/pkg/httpclient"
)

type Client struct {
	resty *resty.Client
}

func New() *Client {
	return &Client{
		resty: httpclient.NewResty(env.OrderURL),
	}
}

func (c *Client) GetOrder(ctx context.Context, id string) (model.Order, error) {
	var order model.Order

	endpoint := fmt.Sprintf("/orders/%s", id)
	req := c.resty.R().SetContext(ctx).SetResult(model.Order{})
	res, err := req.Get(endpoint)
	if err != nil {
		return order, fmt.Errorf("error on get order: %w", err)
	}

	if res.IsError() {
		return order, fmt.Errorf("error on get order, status code: %d", res.StatusCode())
	}

	order = *res.Result().(*model.Order)

	return order, nil
}

func (c *Client) CreateOrder(ctx context.Context, id string) error {
	req := c.resty.R().SetContext(ctx).SetBody(model.Order{ID: id})

	res, err := req.Post("/orders")

	if err != nil {
		return fmt.Errorf("error on create order: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on create order, status code: %d", res.StatusCode())
	}

	return nil
}

func (c *Client) ApproveOrder(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/orders/%s/approve", id)

	req := c.resty.R().SetContext(ctx)

	res, err := req.Patch(endpoint)

	if err != nil {
		return fmt.Errorf("error on approve order: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on approve order, status code: %d", res.StatusCode())
	}

	return nil
}

func (c *Client) FailOrder(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/orders/%s/fail", id)

	req := c.resty.R().SetContext(ctx)

	res, err := req.Patch(endpoint)

	if err != nil {
		return fmt.Errorf("error on fail order: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on fail order, status code: %d", res.StatusCode())
	}

	return nil
}
