package goubus

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

// authData represents the Data response from auth module
type authData struct {
	UbusRPCSession string `json:"ubus_rpc_session"`
	Timeout        int
	Expires        int
	ACLS           acls `json:"acls"`
	Data           map[string]interface{}
}

// acls represents the ACL from user on Authentication
type acls struct {
	AccessGroup map[string][]string `json:"access-group"`
	Ubus        map[string][]string
	Uci         map[string][]string
}

// Login Call JSON-RPC method to Router Authentication
func (u *ubus) Login(username, password string) (*authData, error) {
	u.authData.UbusRPCSession = EmptySession
	res, err := u.Call(
		"session",
		"login",
		map[string]interface{}{
			"username": username,
			"password": password,
		})
	if err != nil {
		return nil, err
	}
	resArray := gjson.Parse(res).Array()
	//fmt.Println(resArray)
	if len(resArray) < 2 {
		return nil, UbusErrorUnknown
	}
	rcode := resArray[0].Int()
	if rcode != 0 {
		return nil, UbusError(int(rcode))
	}
	ad := &authData{}
	if err := json.Unmarshal([]byte(resArray[1].Raw), ad); err != nil {
		return nil, err
	}
	return ad, nil
}

// Logined check if login RPC Session id has expired
func (u *ubus) Logined() error {
	if u.authData.UbusRPCSession == EmptySession {
		return UbusErrorPermissionDenied
	}
	return nil
}
