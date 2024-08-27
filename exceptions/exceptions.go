package exceptions

import "errors"

// User errors
var ErrUserAlreadyExists = errors.New("user: user already exists")
var ErrUserIdIsRequired = errors.New("user: user_id is required and must be at least 3 characters")
var ErrTagIsRequired = errors.New("user: tag is required and must be at least 3 characters")
var ErrUrlIsNotValid = errors.New("user: url is not a valid URL")
var ErrIdIsNotValid = errors.New("item: id is not valid")
var ErrUserNotFound = errors.New("user: user not found")
var ErrBadData = errors.New("user: unprocessable json")
var ErrBadRequest = errors.New("user: can't update without valid field")
var ErrMissingField = errors.New("user: can't create without valid field")
var ErrInternalServer = errors.New("user: internal server error")
var ErrItemIdIsRequired = errors.New("user: item_id is required and must be at least 3 characters")

// Bind errors
var ErrBindDataOnCreateUser = errors.New("user: error on bind user request when creating user")
var ErrBindDataOnUpdateUser = errors.New("user: error on bind user request when updating user")

// DB errors
var ErrCreateUserInDB = errors.New("user: error creating user in the database")
var ErrUpdateUserInDB = errors.New("user: error updating user in the database")
var ErrDeleteUserInDB = errors.New("user: error deleting user in the database")
var ErrGetUserInDB = errors.New("user: error getting user in the database")
var ErrListUsersInDB = errors.New("user: error listing users in the database")
