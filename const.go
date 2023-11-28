package goubus

import "time"

type UbusResponseCode = int

// Represents enum ubus_msg_status from https://git.openwrt.org/?p=project/ubus.git;a=blob;f=ubusmsg.h;h=398b126b6dc01833937749a110181ea0debb1476;hb=HEAD
const (
	UbusStatusOK               UbusResponseCode = 0
	UbusStatusInvalidCommand   UbusResponseCode = 1
	UbusStatusInvalidArgument  UbusResponseCode = 2
	UbusStatusMethodNotFound   UbusResponseCode = 3
	UbusStatusNotFound         UbusResponseCode = 4
	UbusStatusNoData           UbusResponseCode = 5
	UbusStatusPermissionDenied UbusResponseCode = 6
	UbusStatusTimeout          UbusResponseCode = 7
	UbusStatusNotSupported     UbusResponseCode = 8
	UbusStatusUnknown          UbusResponseCode = 9
	UbusStatusConnectionFailed UbusResponseCode = 10
	UbusStatusLast             UbusResponseCode = 11
)

const (
	EmptySession         = "00000000000000000000000000000000"
	DefaultSocketPath    = "/var/run/ubus/ubus.sock"
	DefaultInvokeTimeout = time.Second * 3
)
