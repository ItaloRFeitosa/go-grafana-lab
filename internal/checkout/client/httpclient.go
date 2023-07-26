package client

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/italorfeitosa/go-grafana-lab/internal/checkout/model"
	"github.com/italorfeitosa/go-grafana-lab/internal/env"
	"github.com/italorfeitosa/go-grafana-lab/pkg/httpclient"
)

type Client struct {
	resty *resty.Client
}

func New() *Client {
	return &Client{
		resty: httpclient.NewResty(env.CheckoutURL),
	}
}
func (c *Client) StartCheckout(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/checkouts/%s", id)

	req := c.resty.R().SetContext(ctx)

	res, err := req.Put(endpoint)

	if err != nil {
		return fmt.Errorf("error on start checkout: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on start checkout, status code: %d", res.StatusCode())
	}

	return nil
}

func (c *Client) FinishCheckout(ctx context.Context, payment model.Payment) error {
	endpoint := fmt.Sprintf("/checkouts/%s/finish", payment.CorrelationID)

	req := c.resty.R().SetContext(ctx).SetBody(payment)

	res, err := req.Patch(endpoint)

	if err != nil {
		return fmt.Errorf("error on finish checkout: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on finish checkout, status code: %d", res.StatusCode())
	}

	return nil
}

func (c *Client) StartCheckoutPayment(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/checkouts/%s/payments", id)

	req := c.resty.R().SetContext(ctx)

	res, err := req.Post(endpoint)

	if err != nil {
		return fmt.Errorf("error on start checkout payment: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on start checkout payment, status code: %d", res.StatusCode())
	}

	return nil
}
