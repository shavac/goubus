package goubus

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type UbusFileList struct {
	Entries []UbusFileListData
}

type UbusFileListData struct {
	Name string
	Type string
}

type UbusFileStat struct {
	Path  string
	Type  string
	Size  int
	Mode  int
	Atime int
	Mtime int
	Ctime int
	Inode int
	Uid   int
	Gid   int
}

type UbusFile struct {
	Data string
}

func (u *ubus) FileExec(id int, command string, params []string) (ubusExec, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return ubusExec{}, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"file", 
				"exec", 
				{ 
					"command": "` + command + `",
					"params": [
						"` + strings.Join(params, "\",\n\" ") + `"
					]
				} 
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		if strings.Compare(err.Error(), "Object not found") == 0 {
			return ubusExec{}, errors.New("File module not found, try 'opkg update && opkg install rpcd-mod-file && service rpcd restart'")
		}
		return ubusExec{}, err
	}
	ubusData := ubusExec{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return ubusExec{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}

func (u *ubus) FileWrite(id int, path string, data string, append bool, mode int, base64 bool) error {
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
				"file", 
				"write", 
				{ 
					"path": "` + path + `",
					"data": "` + data + `",
					"append": "` + strconv.FormatBool(append) + `",
					"mode": "` + strconv.Itoa(mode) + `",
					"base64": "` + strconv.FormatBool(base64) + `"
				} 
			] 
		}`)
	_, err := u.JsonRequest(jsonStr)
	if err != nil {
		if strings.Compare(err.Error(), "Object not found") == 0 {
			return errors.New("File module not found, try 'opkg update && opkg install rpcd-mod-file && service rpcd restart'")
		}
		return err
	}
	return nil
}

func (u *ubus) FileRead(id int, path string) (UbusFile, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return UbusFile{}, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"file", 
				"read", 
				{ 
					"path": "` + path + `"
				} 
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		if strings.Compare(err.Error(), "Object not found") == 0 {
			return UbusFile{}, errors.New("File module not found, try 'opkg update && opkg install rpcd-mod-file && service rpcd restart'")
		}
		return UbusFile{}, err
	}
	ubusData := UbusFile{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return UbusFile{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}

func (u *ubus) FileStat(id int, path string) (UbusFileStat, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return UbusFileStat{}, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"file", 
				"stat", 
				{ 
					"path": "` + path + `"
				} 
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		if strings.Compare(err.Error(), "Object not found") == 0 {
			return UbusFileStat{}, errors.New("File module not found, try 'opkg update && opkg install rpcd-mod-file && service rpcd restart'")
		}
		return UbusFileStat{}, err
	}
	ubusData := UbusFileStat{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return UbusFileStat{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}

func (u *ubus) FileList(id int, path string) (UbusFileList, error) {
	errLogin := u.Logined()
	if errLogin != nil {
		return UbusFileList{}, errLogin
	}
	var jsonStr = []byte(`
		{ 
			"jsonrpc": "2.0", 
			"id": ` + strconv.Itoa(id) + `, 
			"method": "call", 
			"params": [ 
				"` + u.authData.UbusRPCSession + `", 
				"file", 
				"list", 
				{ 
					"path": "` + path + `"
				} 
			] 
		}`)
	call, err := u.JsonRequest(jsonStr)
	if err != nil {
		if strings.Compare(err.Error(), "Object not found") == 0 {
			return UbusFileList{}, errors.New("File module not found, try 'opkg update && opkg install rpcd-mod-file && service rpcd restart'")
		}
		return UbusFileList{}, err
	}
	ubusData := UbusFileList{}
	ubusDataByte, err := json.Marshal(call.Result.([]interface{})[1])
	if err != nil {
		return UbusFileList{}, errors.New("Data error")
	}
	json.Unmarshal(ubusDataByte, &ubusData)
	return ubusData, nil
}
