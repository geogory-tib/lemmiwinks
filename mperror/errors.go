package mperror

import (
	"fmt"
	"time"
)

const (
	Invalid_Pass_code       = 10
	Successful_Login_Code   = 7
	User_Not_active_Code    = 11
	Successful_Message_code = 8
	Failure_Message_Code    = 12
)

type Login_User_Invaild_Pass struct {
	message string
	code    int
}
type User_Not_active struct {
	message string
	code    int
}

func (err *Login_User_Invaild_Pass) Error() string {
	err.message = fmt.Sprintf("%s User attempted to login with invaild password", time.Now().Format(time.ANSIC))
	err.code = Invalid_Pass_code
	return fmt.Sprintf("Error %d : %s", err.code, err.message)
}

func (err *User_Not_active) Error() string {
	err.message = fmt.Sprintf("%s User attempted to send message to unactive user", time.Now().Format(time.ANSIC))
	err.code = Invalid_Pass_code
	return fmt.Sprintf("Error %d : %s", err.code, err.message)
}
