package app

type Connect struct {
	Id int
	*Client
}

func NewConnect(id int, client *Client) *Connect {
	return &Connect{
		Id:     id,
		Client: client,
	}
}
