package goubus

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/shavac/httpunix"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

type ubusRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// ubusResult represents a response from JSON-RPC
type ubusResult []interface{}

func (res ubusResult) Code() int {
	if len(res) == 0 {
		return -1
	}
	code, ok := res[0].(float64)
	if !ok {
		return -2
	}
	return int(code)
}

func (res ubusResult) ToString() string {
	if len(res) < 2 {
		return ""
	}
	s := fmt.Sprint(res[1])
	return s
}

func (res ubusResult) ToBytes() []byte {
	return []byte(res.ToString())
}

// ubus represents information to JSON-RPC Interaction with router
type ubus struct {
	endpoint string
	authData
	id int
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

func (u *ubus) RPCRequest(method, ubusObj, ubusMethod string, args map[string]interface{}) (ubusResult, error) {
	jsonReq := u.buildReqestJson(method, ubusObj, ubusMethod, args)
	//slog.Debug(string(jsonReq))
	body, err := u.request(jsonReq)
	if err != nil || body == nil {
		return nil, err
	}
	//slog.Debug(string(body))
	resp := struct {
		JsonRPC string     `json:"jsonrpc"`
		ID      int        `json:"id"`
		Result  ubusResult `json:"result"`
	}{}
	json.Unmarshal(body, &resp)
	//Function Error
	if resp.Result.Code() != UbusStatusOK {
		return nil, UbusError(resp.Result.Code())
	}
	return resp.Result, nil
}

func (u *ubus) RPCCall(ubusObj, ubusMethod string, args map[string]interface{}) (ubusResult, error) {
	return u.RPCRequest("call", ubusObj, ubusMethod, args)
}

func httpRequest(url string, jsonStr []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error %s on (%s)", resp.Status, url)
	}
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func socketRequest(filepath string, jsonStr []byte) ([]byte, error) {
	tr := &httpunix.Transport{
		DialTimeout:           500 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	tr.RegisterLocation("myservice", filepath)

	var client = http.Client{
		Transport: tr,
	}

	resp, err := client.Post("http+unix://myservice", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error(err.Error())
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error %s", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}
