package exceptions

import "errors"

// Restaurant errors
var ErrRestaurantAlreadyExists = errors.New("restaurant: restaurant already exists")
var ErrRestaurantIdIsRequired = errors.New("restaurant: restaurant_id is required and must be at least 3 characters")
var ErrTagIsRequired = errors.New("restaurant: id is required and must be at least 3 characters")
var ErrUrlIsNotValid = errors.New("restaurant: url is not a valid URL")
var ErrTagIsNotValid = errors.New("restaurant: id is not valid")
var ErrRestaurantNotFound = errors.New("restaurant: restaurant not found")
var ErrBadData = errors.New("restaurant: unprocessable json")
var ErrBadRequest = errors.New("restaurant: can't update without valid field")
var ErrMissingField = errors.New("restaurant: can't create without valid field")
var ErrInternalServer = errors.New("restaurant: internal server error")
var ErrOrderNotFound = errors.New("restaurant: order not found")

// Bind errors
var ErrBindDataOnCreateRestaurant = errors.New("restaurant: error on bind restaurant request when creating restaurant")
var ErrBindDataOnUpdateRestaurant = errors.New("restaurant: error on bind restaurant request when updating restaurant")

// DB errors
var ErrCreateRestaurantInDB = errors.New("restaurant: error creating restaurant in the database")
var ErrUpdateRestaurantInDB = errors.New("restaurant: error updating restaurant in the database")
var ErrDeleteRestaurantInDB = errors.New("restaurant: error deleting restaurant in the database")
var ErrGetRestaurantInDB = errors.New("restaurant: error getting restaurant in the database")
var ErrListRestaurantsInDB = errors.New("restaurant: error listing restaurants in the database")
