package jsontypes

type Message_i interface {
	Get_Content() string
	Get_Time() string
}

type User_J struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

type Message_J struct {
	To       string `json:"to"`
	Contents string `json:"content"`
	From     string `json:"from"`
	Time     string `json:"time"`
}

func (msg Message_J) Get_Content() string {
	return msg.Contents
}
func (msg Message_J) Get_Time() string {
	return msg.Time
}

type Server_message_J struct {
	Code     int    `json:"code"`
	Contents string `json:"contents"`
	Time     string `json:"time"`
}
