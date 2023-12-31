package goubus

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

func (u *UBus) LogWrite(id int, event string) error {
	errLogin := u.Logined()
	if errLogin != nil {
		return errLogin
	}
	evtString := map[string]interface{}{"event": event}
	_, err := u.RPCRequest("call", "log", "write", evtString)
	if err != nil {
		return err
	}
	return nil
}

/*
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
	call, err := u.RequestJson(jsonStr)
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
*/
