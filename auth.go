package goubus

import (
	"encoding/json"
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
	ad := authData{}
	if err := json.Unmarshal(res.toBytes(), &ad); err != nil {
		return nil, err
	}
	u.authData = ad
	return &ad, nil
}

// Logined check if login RPC Session id has expired
func (u *ubus) Logined() error {
	if u.authData.UbusRPCSession == EmptySession {
		return UbusErrorPermissionDenied
	}
	return nil
}
