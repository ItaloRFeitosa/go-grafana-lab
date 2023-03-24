package webhook

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/italorfeitosa/go-grafana-lab/env"
	"github.com/italorfeitosa/go-grafana-lab/httpclient"
)

type CheckoutClient struct {
	resty *resty.Client
}

func NewCheckoutClient() *CheckoutClient {
	return &CheckoutClient{
		resty: httpclient.NewResty(env.CheckoutURL),
	}
}

func (c *CheckoutClient) FinishCheckout(ctx context.Context, payment Payment) error {
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
