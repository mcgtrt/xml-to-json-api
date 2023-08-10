package api

import "net/http"

type Error struct {
	Status int    `json:"status"`
	Err    string `json:"err"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(status int, err string) Error {
	return Error{
		Status: status,
		Err:    err,
	}
}

func ErrBadRequest() Error {
	return Error{
		Status: http.StatusBadRequest,
		Err:    "bad request",
	}
}

func ErrInternalServerError() Error {
	return Error{
		Status: http.StatusInternalServerError,
		Err:    "internal server error",
	}
}

func ErrNoDocuments() Error {
	return Error{
		Status: http.StatusNoContent,
		Err:    "nothing found",
	}
}
