package exceptions

import (
	"errors"
	"net/http"
)

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandleException(err error) ErrResponse {
	customErr, ok := err.(*Error)
	if !ok {
		return ErrResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	switch err != nil {
	case
		errors.Is(customErr.CustomErr, ErrUserIdIsRequired),
		errors.Is(customErr.CustomErr, ErrTagIsRequired),
		errors.Is(customErr.CustomErr, ErrUrlIsNotValid),
		errors.Is(customErr.CustomErr, ErrIdIsNotValid),
		errors.Is(customErr.CustomErr, ErrBadRequest),
		errors.Is(customErr.CustomErr, ErrItemIdIsRequired):
		return ErrResponse{
			Code:    http.StatusBadRequest,
			Message: customErr.CustomErr.Error(),
		}
	case
		errors.Is(customErr.CustomErr, ErrCreateUserInDB),
		errors.Is(customErr.CustomErr, ErrGetUserInDB),
		errors.Is(customErr.CustomErr, ErrListUsersInDB),
		errors.Is(customErr.CustomErr, ErrUpdateUserInDB),
		errors.Is(customErr.CustomErr, ErrDeleteUserInDB),
		errors.Is(customErr.CustomErr, ErrBindDataOnCreateUser),
		errors.Is(customErr.CustomErr, ErrBindDataOnUpdateUser):
		return ErrResponse{
			Code:    http.StatusInternalServerError,
			Message: customErr.CustomErr.Error(),
		}
	case
		errors.Is(customErr.CustomErr, ErrBadData):
		return ErrResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: customErr.CustomErr.Error(),
		}
	case
		errors.Is(customErr.CustomErr, ErrUserNotFound):
		return ErrResponse{
			Code:    http.StatusNotFound,
			Message: customErr.CustomErr.Error(),
		}
	default:
		return ErrResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
}
