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

type ubusRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	ID      int64         `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// ubusResult represents a response from JSON-RPC
type ubusResult []interface{}

func (res ubusResult) toString() string {
	return string(res.toBytes())
}

func (res ubusResult) toBytes() []byte {
	b, _ := json.Marshal(res)
	return b
}

// ubus represents information to JSON-RPC Interaction with router
type ubus struct {
	endpoint string
	authData
	id int64
	//jsonrpc func(method, object, ubusMethod string, args ...string) ubusResponse
	request func([]byte) ([]byte, error)
}

func NewUbus(endp string) (*ubus, error) {
	if len(endp) == 0 {
		endp = DefaultSocketPath
	}
	u, err := url.Parse(endp)
	if err != nil {
		return nil, err
	}
	ub := &ubus{endpoint: endp, id: 1, authData: authData{UbusRPCSession: EmptySession}}
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

func (u *ubus) buildReqestJson(method, ubusObj, ubusMethod string, args map[string]interface{}) []byte {
	req := &ubusRequest{
		JsonRPC: "2.0",
		ID:      u.id,
		Method:  method,
		Params: []interface{}{
			u.authData.UbusRPCSession,
			ubusObj,
			ubusMethod,
			args,
		},
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return []byte{}
	}
	return jsonReq
}

func (u *ubus) RPCRequest(method, ubusObj, ubusMethod string, args map[string]interface{}) (string, error) {
	jsonReq := u.buildReqestJson(method, ubusObj, ubusMethod, args)
	//slog.Debug(string(jsonReq))
	body, err := u.request(jsonReq)
	if err != nil || body == nil {
		return "", err
	}
	//slog.Debug(string(body))
	/*
		resp := struct {
			JsonRPC string     `json:"jsonrpc"`
			ID      int        `json:"id"`
			Result  ubusResult `json:"result"`
			Error   struct {
				Code    int
				Message string
			} `json:"error"`
		}{}
		json.Unmarshal(body, &resp)
	*/
	id := gjson.GetBytes(body, "id").Int()
	//Function Error
	if id != u.id {
		return "", SysErrorIDMismatch
	}
	u.id++
	res := gjson.GetBytes(body, "result").Raw
	return res, nil
}

func (u *ubus) Call(ubusObj, ubusMethod string, args map[string]interface{}) (string, error) {
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
