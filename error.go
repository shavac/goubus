package goubus

import (
	"fmt"
)

type ubusError struct {
	Code    int
	Message string
}

func (uErr ubusError) Error() string {
	return fmt.Sprintf("UBUS ERROR CODE(%d): %s", uErr.Code, uErr.Message)
}

var (
	UbusErrorPermissionDenied = ubusError{UbusStatusPermissionDenied, "Permission Denied"}
	UbusErrorInvalidCommand   = ubusError{UbusStatusInvalidCommand, "Invalid Command"}
	UbusErrorInvalidArgument  = ubusError{UbusStatusInvalidArgument, "Invalid Argument"}
)

func UbusError(code int) error {
	if code == UbusStatusOK {
		return nil
	}
	errMap := map[int]error{
		UbusStatusInvalidCommand:   UbusErrorInvalidCommand,
		UbusStatusInvalidArgument:  UbusErrorInvalidArgument,
		UbusStatusMethodNotFound:   nil,
		UbusStatusNotFound:         nil,
		UbusStatusNoData:           nil,
		UbusStatusPermissionDenied: UbusErrorPermissionDenied,
		UbusStatusTimeout:          nil,
		UbusStatusNotSupported:     nil,
		UbusStatusUnknownError:     nil,
		UbusStatusConnectionFailed: nil,
		UbusStatusLast:             nil,
	}
	err, ok := errMap[code]
	if !ok || err == nil {
		err = ubusError{code, "Not implemented error code"}
	}
	return err
}
