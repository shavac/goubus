package goubus

import (
	"fmt"
	"testing"
)

func Test_ubus_buildJson(t *testing.T) {
	type fields struct {
		endpoint string
		authData authData
		id       int64
		reqFunc  func([]byte) ([]byte, error)
	}
	type args struct {
		cmd        string
		obj        string
		ubusmethod string
		args       map[string]string
	}
	u, _ := NewUbus("")
	jsonReq := u.buildJson(
		"call",
		"session",
		"login",
		map[string]interface{}{
			"username": "root",
			"password": "pass",
		})
	fmt.Println(string(jsonReq))
}

func Test_ubus_Login(t *testing.T) {
	u, _ := NewUbus("http://192.168.1.1/ubus")
	ok, err := u.Login("root", "P@5sw0rd")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ok, u.authData.UbusRPCSession)
}
