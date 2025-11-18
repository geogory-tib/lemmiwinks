package inbox

// plan on adding a message cache to save conversations
type Message_t struct {
	Content string //`json:"content"`
	Time    string //`json:"time"`
	From    string //`json:"from"`
}
type Inbox_T struct {
	Inboxes  map[string][]Message_t
	Username string
}

func Init_Inbox(username string) (ret Inbox_T) {
	ret.Username = username
	ret.Inboxes = make(map[string][]Message_t)
	return
}
