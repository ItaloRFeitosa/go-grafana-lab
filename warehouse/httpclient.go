package warehouse

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
		resty: httpclient.NewResty(env.WarehouseURL),
	}
}

func (c *Client) OrderDispatch(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/warehouse/orders/%s/dispatch", id)
	req := c.resty.R().SetContext(ctx)
	res, err := req.Patch(endpoint)
	if err != nil {
		return fmt.Errorf("error on dispatch: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on dispatch, status code: %d", res.StatusCode())
	}

	return nil
}

func (c *Client) OrderPrepare(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/warehouse/orders/%s/prepare", id)
	req := c.resty.R().SetContext(ctx)
	res, err := req.Patch(endpoint)
	if err != nil {
		return fmt.Errorf("error on prepare order: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error on prepare order, status code: %d", res.StatusCode())
	}

	return nil
}
