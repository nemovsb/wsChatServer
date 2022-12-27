package app

type Client struct {
	Nickname string
}

func NewClient(nick string) *Client {
	return &Client{
		Nickname: nick,
	}
}
