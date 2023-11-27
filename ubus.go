package goubus

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shavac/httpunix"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

// ubus represents information to JSON-RPC Interaction with router
type ubus struct {
	endpoint string
	authData
	id int
	//jsonrpc func(method, object, ubusMethod string, args ...string) ubusResponse
	reqFunc func([]byte) ([]byte, error)
}

// authData represents the Data response from auth module
type authData struct {
	UbusRPCSession string `json:"ubus_rpc_session"`
	Timeout        int
	Expires        int
	ACLS           acls `json:"acls"`
	Data           map[string]string
}

// acls represents the ACL from user on Authentication
type acls struct {
	AccessGroup map[string][]string `json:"access-group"`
	Ubus        map[string][]string
	Uci         map[string][]string
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
		ub.reqFunc = func(jsonStr []byte) ([]byte, error) {
			return httpRequest(endp, jsonStr)
		}
	} else if u.Scheme == "" {
		ub.reqFunc = func(jsonStr []byte) ([]byte, error) {
			return socketRequest(endp, jsonStr)
		}
	}
	return ub, nil
}

// ubusResponse represents a response from JSON-RPC
type ubusResponse struct {
	JSONRPC          string
	ID               int
	Error            UbusResponseError
	Result           interface{}
	UbusResponseCode UbusResponseCode
}

type UbusResponseError struct {
	Code    int
	Message string
}

type ubusExec struct {
	Code   int
	Stdout string
}

type ubusParams struct {
	UbusRPCSession string
	ObjectName     string
	ubusMethod     string
	arguments      map[string]string
}

type ubusRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func (u *ubus) buildJson(cmd, obj, ubusMethod string, args map[string]interface{}) []byte {
	req := &ubusRequest{
		JsonRPC: "2.0",
		ID:      u.id,
		Method:  cmd,
		Params: []interface{}{
			u.authData.UbusRPCSession,
			obj,
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

func (u *ubus) JsonRPC(cmd, obj, ubusMethod string, args map[string]interface{}) (*ubusResponse, error) {
	jsonReq := u.buildJson(cmd, obj, ubusMethod, args)
	resp, err := u.JsonRequest(jsonReq)
	if err != nil || resp == nil {
		return nil, err
	}
	//Function Error
	if resp.Error.Code != 0 {
		if resp.Error.Code == UbusStatusPermissionDenied {
			return nil, errors.New("Access denied for this instance, read https://openwrt.org/docs/techref/ubus#acls ")
		}
		return nil, errors.New(resp.Error.Message)
	}
	//Workaround cause response code not contempled by unmarshal function
	resp.UbusResponseCode = UbusResponseCode(resp.Result.([]interface{})[0].(float64))
	if resp.UbusResponseCode != UbusStatusOK {
		return resp, fmt.Errorf("Ubus Status Failed: %d", resp.UbusResponseCode)
	}
	return resp, nil
}

// Login Call JSON-RPC method to Router Authentication
func (u *ubus) Login(username, password string) (bool, error) {
	u.authData.UbusRPCSession = EmptySession
	resp, err := u.JsonRPC("call",
		"session",
		"login",
		map[string]interface{}{
			"username": username,
			"password": password,
		})
	if err != nil {
		return false, err
	}
	ubusData := authData{}
	ubusDataByte, err := json.Marshal(resp.Result.([]interface{})[1])
	if err != nil {
		return false, errors.New("Error Parsing Login Data")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	u.authData = ubusData
	return true, nil
}

// Logined check if login RPC Session id has expired
func (u *ubus) Logined() error {
	if u.authData.UbusRPCSession == EmptySession {
		return errors.New("Not logined error")
	}
	return nil
}

// JsonRequest do a request to Json-RPC to get/set information
func (u *ubus) JsonRequest(jsonStr []byte) (*ubusResponse, error) {
	body, err := u.reqFunc(jsonStr)
	if err != nil {
		return nil, err
	}
	resp := &ubusResponse{}
	json.Unmarshal([]byte(body), &resp)
	//Function Error
	if resp.Error.Code != 0 {
		if strings.Compare(resp.Error.Message, "Access denied") == 0 {
			return nil, errors.New("Access denied for this instance, read https://openwrt.org/docs/techref/ubus#acls ")
		}
		return nil, errors.New(resp.Error.Message)
	}
	//Workaround cause response code not contempled by unmarshal function
	resp.UbusResponseCode = UbusResponseCode(resp.Result.([]interface{})[0].(float64))
	//Workaround to get UbusData cause the structure of this array has a problem with unmarshal
	if resp.UbusResponseCode == UbusStatusOK {
		return resp, nil
	}
	return nil, fmt.Errorf("Ubus Status Failed: %d", resp.UbusResponseCode)
}

func httpRequest(url string, jsonStr []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error %s", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}
