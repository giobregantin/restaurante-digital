package exceptions

import "fmt"

type Error struct {
	CustomErr error
	Err       error
}

func (e Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s - %s", e.CustomErr.Error(), e.Err.Error())
	}
	return e.CustomErr.Error()
}

func New(customError error, err error) *Error {
	return &Error{
		CustomErr: customError,
		Err:       err,
	}
}