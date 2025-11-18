package listener

import (
	"gerbil/inbox"
	"gerbil/jsontypes"
	"gerbil/util"
	"net"
)

// listens to incoming messages
type Listener_T struct {
	Server_Conn net.Conn
	Msg_Chan    chan jsontypes.Message_i
}

func Login_and_Init() (listener Listener_T, inbox inbox.Inbox_T) {
	util.Todo()
	return
}
