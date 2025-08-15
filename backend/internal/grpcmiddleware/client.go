package grpcmiddleware

import (
	"context"
	"github.com/asim/go-micro/v3/client"
)

type Client struct {
	grpcClient client.Client
}

func NewClient() *Client {
	return &Client{
		grpcClient: client.NewClient(),
	}
}

func (c *Client) NewRequest(service, method string, req interface{}) client.Request {
	return c.grpcClient.NewRequest(service, method, req)
}

func (c *Client) Call(ctx context.Context, req client.Request, rsp interface{}) error {
	return c.grpcClient.Call(ctx, req, rsp)
}
