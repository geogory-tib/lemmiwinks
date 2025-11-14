package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"lemiwinks/mperror"
	"lemiwinks/server_cmds"
	"lemiwinks/util"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

type Message struct {
	To       string `json:"to"`
	Contents string `json:"content"`
	From     string `json:"from"`
	Time     string `json:"time"`
}

type Server_message struct {
	Code     int    `json:"code"`
	Contents string `json:"contents"`
	Time     string `json:"time"`
}
type active_manager struct {
	active map[string]net.Conn
	mtx    sync.Mutex
}
type user_manager struct {
	users map[string][]byte
	mtx   sync.Mutex
}

type Server_t struct {
	active   active_manager
	user     user_manager
	listener net.Listener
	log_file *os.File
}

func New_Server_Instance(default_conf []byte) (ret Server_t) {
	var config_json []byte
	file, err := os.Open("config.json")
	if err != nil {
		log.Println("Error Loading Config File using embedded defaults")
		file = nil
		config_json = default_conf
	} else {
		defer file.Close()
		config_data, err := io.ReadAll(file)
		if err != nil {
			log.Panic(err)
		}
		config_json = config_data
	}
	config_struct := struct {
		Ip       string `json:"ip"`
		Port     string `json:"port"`
		Log_path string `json:"log_file"`
	}{}
	json.Unmarshal(config_json, &config_struct)
	log_file, err := os.Create(config_struct.Log_path)
	if err != nil {
		log.Panic(log_file)
	}
	ret.log_file = log_file
	log.SetOutput(log_file)
	port_and_ip := fmt.Sprintf("%s:%s", config_struct.Ip, config_struct.Port)
	ret.listener, err = net.Listen("tcp", port_and_ip)
	if err != nil {
		log.Panic(err)
		log_file.Close()
	}
	ret.active.active = make(map[string]net.Conn)
	ret.user.users = make(map[string][]byte)
	return
}

func (srv *Server_t) Server_Main() {
	defer srv.log_file.Close()
	defer srv.listener.Close()
	for {
		conn, err := srv.listener.Accept()
		if err != nil {
			log.Println(err)
		} else {
			go srv.handle_conn(conn)
		}
	}
}

func (srv *Server_t) handle_conn(conn net.Conn) {
	buffer := make([]byte, 6144)
	user_name := ""
	conn_reader := bufio.NewReader(conn)
	for err := srv.handle_login(conn, conn_reader, buffer, &user_name); err != nil; {
	}
	mesage_buffer := Message{}
	for {
		n, err := util.ReadFullJson(conn_reader, buffer)
		if err != nil {
			log.Print(err)
		}
		json.Unmarshal(buffer[:n], &mesage_buffer)
		mesage_buffer.From = user_name
		if mesage_buffer.To == "server" {
			srv.handle_server_cmd(conn, mesage_buffer)
		} else {
			srv.handle_message(conn, mesage_buffer)
		}
	}
}

func (srv *Server_t) handle_login(conn net.Conn, reader *bufio.Reader, buffer []byte, username *string) error {
	user_info := User{}
	n, err := util.ReadFullJson(reader, buffer)
	login_json := buffer[:n]
	if err != nil {
		log.Print(err)
		return err
	}
	server_error := Server_message{}
	err = json.Unmarshal(login_json, &user_info)
	if err != nil {
		log.Print(err)
		return err
	}
	srv.user.mtx.Lock()
	defer srv.user.mtx.Unlock()
	pass, user_exists := srv.user.users[user_info.Username]
	if user_exists {
		pass_match_err := bcrypt.CompareHashAndPassword(pass[:], []byte(user_info.Password))
		if pass_match_err != nil {
			server_error = Server_message{
				Code: mperror.Invalid_Pass_code,
				Time: time.Now().Format(time.ANSIC),
			}
			json_err, err := json.Marshal(server_error)
			if err != nil {
				log.Print(err)
			}
			conn.Write(json_err)
			return &mperror.Login_User_Invaild_Pass{}
		}
	} else {
		hashed_pass, err := bcrypt.GenerateFromPassword([]byte(user_info.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Print(err)
			return err
		}
		srv.user.users[user_info.Username] = hashed_pass
	}
	server_error = Server_message{
		Code: mperror.Successful_Login_Code,
		Time: time.Now().Format(time.ANSIC),
	}
	json_msg, err := json.Marshal(server_error)
	if err != nil {
		log.Print(err)
		return err
	}
	conn.Write(json_msg)
	srv.active.mtx.Lock()
	defer srv.active.mtx.Unlock()
	srv.active.active[user_info.Username] = conn
	*username = user_info.Username
	return nil
}

func (srv *Server_t) handle_server_cmd(conn net.Conn, cmmd_msg Message) {
	switch cmmd_msg.Contents {
	case servercmds.Show_Active_User:
		srv.cmmd_show_active_users(conn)
	}
}

func (srv *Server_t) handle_message(conn net.Conn, mesage Message) error {
	srv.active.mtx.Lock()
	server_reply := Server_message{}
	defer srv.active.mtx.Unlock()
	dest_conn, user_active := srv.active.active[mesage.To]
	if !user_active {
		server_reply := Server_message{
			Code: mperror.User_Not_active_Code,
			Time: time.Now().Format(time.ANSIC),
		}
		error_json, err := json.Marshal(server_reply)
		if err != nil {
			log.Print(err)
		}
		conn.Write(error_json)
		return nil
	}
	message_json, err := json.Marshal(mesage)
	if err != nil {
		log.Print(err)
		return err
	}
	_, err = dest_conn.Write(message_json)
	server_reply = Server_message{
		Code: mperror.Successful_Login_Code,
		Time: time.Now().Format(time.ANSIC),
	}
	if err != nil {
		log.Print(err)

		server_reply := Server_message{
			Code: mperror.Failure_Message_Code,
			Time: time.Now().Format(time.ANSIC),
		}
		error_json, err := json.Marshal(server_reply)
		if err != nil {
			log.Print(err)
		}
		conn.Write(error_json)
		return err
	}
	server_reply_json, err := json.Marshal(server_reply)
	if err != nil {
		log.Print(err)
		return err
	}
	conn.Write(server_reply_json)
	return nil
}

// commands section
func (srv *Server_t) cmmd_show_active_users(conn net.Conn) {
	srv.active.mtx.Lock()
	defer srv.active.mtx.Unlock()
	var builder strings.Builder
	builder.Grow(500)
	for key := range srv.active.active {
		builder.WriteString(key + "\n")
	}
	server_response := Server_message{
		Code:     servercmds.Show_Active_User_code,
		Contents: builder.String(),
		Time:     time.Now().Format(time.ANSIC),
	}
	server_response_marsh, err := json.Marshal(server_response)
	if err != nil {
		log.Print(err)
		return
	}
	conn.Write(server_response_marsh)
}
