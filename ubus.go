package goubus

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

type UbusParamMap map[string]interface{}

type ubusParams struct {
	UbusRPCSession,
	UbusObj,
	UbusMethod string
	UbusParamMap `json:",omitempty"`
}

type ubusRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	ID      int64         `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,ommitempty"`
}

// UBus represents information to JSON-RPC Interaction with router
type UBus struct {
	endpoint string
	authData
	id int64
	//jsonrpc func(method, object, ubusMethod string, args ...string) ubusResponse
	request func([]byte) ([]byte, error)
}

func NewUbus(endp string) (*UBus, error) {
	if len(endp) == 0 {
		endp = DefaultSocketPath
	}
	u, err := url.Parse(endp)
	if err != nil {
		return nil, err
	}
	ub := &UBus{endpoint: endp, id: 1, authData: authData{UbusRPCSession: EmptySession}}
	if u.Scheme == "http" || u.Scheme == "https" {
		ub.request = func(jsonStr []byte) ([]byte, error) {
			return httpRequest(endp, jsonStr)
		}
	} else if u.Scheme == "" {
		ub.request = func(jsonStr []byte) ([]byte, error) {
			return socketRequest(endp, jsonStr)
		}
	}
	return ub, nil
}

func (u *UBus) buildReqestJson(method, ubusObj, ubusMethod string, args map[string]interface{}) []byte {
	if args == nil {
		args = map[string]interface{}{}
	}
	params := []interface{}{u.UbusRPCSession, ubusObj, ubusMethod, args}
	req := &ubusRequest{
		JsonRPC: "2.0",
		ID:      u.id,
		Method:  method,
		Params:  params,
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return []byte{}
	}
	return jsonReq
}

func (u *UBus) RPCRequest(method, ubusObj, ubusMethod string, args map[string]interface{}) (string, error) {
	jsonReq := u.buildReqestJson(method, ubusObj, ubusMethod, args)
	//slog.Debug(string(jsonReq))
	body, err := u.request(jsonReq)
	if err != nil || body == nil {
		return "", err
	}
	bd := gjson.ParseBytes(body)
	id := bd.Get("id").Int()
	//Function Error
	if id != u.id {
		return "", SysErrorIDMismatch
	}
	u.id++
	if err := bd.Get("error"); err.Exists() {
		code := err.Get("code").Int()
		return "", UbusError(int(code))
	}
	res := gjson.GetBytes(body, "result").Raw
	return res, nil
}

func (u *UBus) Call(ubusObj, ubusMethod string, args map[string]interface{}) (string, error) {
	return u.RPCRequest("call", ubusObj, ubusMethod, args)
}

func httpRequest(url string, jsonStr []byte) ([]byte, error) {
	log.Debug("URL:", url, "REQ:", string(jsonStr), "\n")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, httpError(resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	log.Debug("URL:", url, "RES:", string(body), "\n")
	return body, nil
}

func socketRequest(filepath string, jsonStr []byte) ([]byte, error) {
	return nil, SysErrorNotImplemented
}
