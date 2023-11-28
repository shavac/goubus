package goubus

import (
	"github.com/tidwall/gjson"
)

func (u *ubus) List(ubusObj, ubusMethod string, args map[string]interface{}) (map[string]gjson.Result, error) {
	res, err := u.RPCRequest("list", ubusObj, ubusMethod, args)
	if err != nil {
		return nil, err
	}
	m := gjson.Parse(res).Map()
	return m, nil
}
