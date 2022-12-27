package app

type Message struct {
	ClientName string `json:"client_name"`
	ChanelName string `json:"chanel_name"`
	Text       string `json:"text"`
}
