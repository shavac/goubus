package goubus

import (
	"encoding/json"
	"errors"
	"strconv"
)

type UbusUciConfigs struct {
	Configs []string
}

type UbusUciRequest struct {
	Config  string            `json:"config"`
	Section string            `json:"section,omitempty"`
	Option  string            `json:"option,omitempty"`
	Type    string            `json:"type,omitempty"`
	Match   string            `json:"match,omitempty"`
	Values  map[string]string `json:"values,omitempty"`
}

type UbusUciResponse struct {
	Values interface{}
}

func (u *ubus) UciGetConfigs(id int) (UbusUciConfigs, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return UbusUciConfigs{}, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"uci", 
				"configs", 
				{} 
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		return UbusUciConfigs{}, err
	}
	ubusData := UbusUciConfigs{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return UbusUciConfigs{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}

func (u *ubus) UciGetConfig(id int, request UbusUciRequest) (UbusUciResponse, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return UbusUciResponse{}, errLogin
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return UbusUciResponse{}, errors.New("Error Parsing UCI Request Data")
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"uci", 
				"get", 
				` + string(jsonData) + ` 
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		return UbusUciResponse{}, err
	}
	ubusData := UbusUciResponse{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return UbusUciResponse{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}

func (u *ubus) UciSetConfig(id int, request UbusUciRequest) error {
	errLogin := u.Logined()
	if errLogin != nil {
		return errLogin
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return errors.New("Error Parsing UCI Request Data")
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"uci", 
				"set", 
				` + string(jsonData) + ` 
			] 
		}`)
	_, err = u.JsonRequest(jsonStr)
	if err != nil {
		return err
	}
	return nil
}

func (u *ubus) UciChanges(id int) (map[string]map[string][][]string, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return nil, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"uci", 
				"changes", 
				{}
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		return nil, err
	}
	// fmt.Println(call)
	var ubusData map[string]map[string][][]string
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return nil, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}

func (u *ubus) UciCommit(id int) error {
	errLogin := u.Logined()
	if errLogin != nil {
		return errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"uci", 
				"commit", 
				{}
			] 
		}`)
	_, err := u.JsonRequest(jsonStr)
	if err != nil {
		return err
	}
	return nil
}

func (u *ubus) UciReloadConfig(id int) error {
	errLogin := u.Logined()
	if errLogin != nil {
		return errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"uci", 
				"reload_config",
				{}
			] 
		}`)
	_, err := u.JsonRequest(jsonStr)
	if err != nil {
		return err
	}
	return nil
}
