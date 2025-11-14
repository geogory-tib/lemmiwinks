package main

import (
	_ "embed"

	"lemiwinks/server"
)

//go:embed defaultconf.json
var defaultconf []byte

func main() {
	srv := server.New_Server_Instance(defaultconf)
	srv.Server_Main()
}
