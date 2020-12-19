package exactonline_bq

import (
	eo "github.com/leapforce-libraries/go_exactonline_new"
	ec "github.com/leapforce-libraries/go_exactonline_new/crm"
)

type Client struct {
	clientID    string
	exactOnline *eo.ExactOnline
}

func NewClient(clientID string, exactOnline *eo.ExactOnline) *Client {
	return &Client{clientID, exactOnline}
}

func (c *Client) CRMClient() *ec.Client {
	return c.exactOnline.CRMClient
}
func (c *Client) ClientID() string {
	return c.clientID
}
