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
		errors.Is(customErr.CustomErr, ErrRestaurantIdIsRequired),
		errors.Is(customErr.CustomErr, ErrTagIsRequired),
		errors.Is(customErr.CustomErr, ErrUrlIsNotValid),
		errors.Is(customErr.CustomErr, ErrTagIsNotValid),
		errors.Is(customErr.CustomErr, ErrBadRequest):
		return ErrResponse{
			Code:    http.StatusBadRequest,
			Message: customErr.CustomErr.Error(),
		}
	case
		errors.Is(customErr.CustomErr, ErrCreateRestaurantInDB),
		errors.Is(customErr.CustomErr, ErrGetRestaurantInDB),
		errors.Is(customErr.CustomErr, ErrListRestaurantsInDB),
		errors.Is(customErr.CustomErr, ErrUpdateRestaurantInDB),
		errors.Is(customErr.CustomErr, ErrDeleteRestaurantInDB),
		errors.Is(customErr.CustomErr, ErrBindDataOnCreateRestaurant),
		errors.Is(customErr.CustomErr, ErrBindDataOnUpdateRestaurant):
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
		errors.Is(customErr.CustomErr, ErrRestaurantNotFound),
		errors.Is(customErr.CustomErr, ErrOrderNotFound):
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
