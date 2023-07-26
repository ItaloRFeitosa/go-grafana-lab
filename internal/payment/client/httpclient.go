package client

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/italorfeitosa/go-grafana-lab/internal/env"
	"github.com/italorfeitosa/go-grafana-lab/internal/payment/model"
	"github.com/italorfeitosa/go-grafana-lab/pkg/httpclient"
)

type Client struct {
	resty *resty.Client
}

func New() *Client {
	return &Client{
		resty: httpclient.NewResty(env.PaymentURL),
	}
}

func (c *Client) CreatePayment(ctx context.Context, paym model.Payment) error {
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
