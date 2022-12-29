package app

type Message struct {
	FromConnId int
	ClientName string `json:"client_name"`
	ChanelName string `json:"chanel_name"`
	Text       string `json:"text"`
}

type SendTask struct {
	M  Message
	To Connect
}
