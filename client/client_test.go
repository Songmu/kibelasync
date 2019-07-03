package client

func Test(cli Doer) *Client {
	return &Client{cli: cli}
}
