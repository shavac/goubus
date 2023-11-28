package goubus

func (u *ubus) List(ubusObj, ubusMethod string, args map[string]interface{}) (*ubusResult, error) {
	return u.RPCRequest("list", ubusObj, ubusMethod, args)
}
