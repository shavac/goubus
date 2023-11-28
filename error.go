package goubus

import (
	"errors"
	"fmt"
	"net/http"
)

type httpError int

func (herr httpError) Error() string {
	return http.StatusText(int(herr))
}

var (
	SysErrorIDMismatch     = errors.New("response id mismatch")
	SysErrorNotImplemented = errors.New("function not implemented")
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
	UbusErrorUnknown          = ubusError{UbusStatusUnknown, "Unknown Error"}
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
		UbusStatusUnknown:          UbusErrorUnknown,
		UbusStatusConnectionFailed: nil,
		UbusStatusLast:             nil,
	}
	err, ok := errMap[code]
	if !ok || err == nil {
		err = ubusError{code, "Not implemented error code"}
	}
	return err
}
