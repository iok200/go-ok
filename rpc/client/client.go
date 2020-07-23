package client

import "google.golang.org/grpc"

type Client struct {
	addr string
}

func New() *Client {
	return &Client{addr: addr}
}

func (this *Client) Conn() error {
	conn, err := grpc.Dial(this.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	return nil
}
