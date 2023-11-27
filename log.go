package goubus

import (
	"encoding/json"
	"errors"
	"strconv"
)

type UbusLog struct {
	Log []UbusLogData
}

type UbusLogData struct {
	Msg      string
	ID       int
	Priority int
	Source   int
	Time     int
}

func (u *ubus) LogWrite(id int, event string) error {
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
				"log", 
				"write", 
				{ 
					"event": "` + event + `"
				} 
			] 
		}`)
	_, err := u.JsonRequest(jsonStr)
	if err != nil {
		return err
	}
	return nil
}

func (u *ubus) LogRead(id int, lines int, stream bool, oneshot bool) (UbusLog, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return UbusLog{}, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"log", 
				"read", 
				{ 
					"lines": ` + strconv.Itoa(lines) + `,
					"stream": ` + strconv.FormatBool(stream) + `,
					"oneshot":` + strconv.FormatBool(oneshot) + `
				} 
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		return UbusLog{}, err
	}
	ubusData := UbusLog{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return UbusLog{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}
