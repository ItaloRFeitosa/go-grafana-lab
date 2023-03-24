package payment

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/italorfeitosa/go-grafana-lab/env"
	"github.com/italorfeitosa/go-grafana-lab/httpclient"
)

type Client struct {
	resty *resty.Client
}

func NewClient() *Client {
	return &Client{
		resty: httpclient.NewResty(env.PaymentURL),
	}
}

func (c *Client) CreatePayment(ctx context.Context, paym Payment) error {
	req := c.resty.R().SetContext(ctx).SetBody(paym)

	res, err := req.Post("/payments")

	if err != nil {
		return fmt.Errorf("error on create payment: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on create payment, status code: %d", res.StatusCode())
	}

	return nil
}
