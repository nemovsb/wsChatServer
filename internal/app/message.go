package app

type Message struct {
	//Id         int
	ClientName string `json:"client_name"`
	ChanelName string `json:"chanel_name"`
	Text       string `json:"text"`
}
