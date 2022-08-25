package mattermost

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func AppErrorFromJSON(data io.Reader) *AppError {
	str := ""
	bytes, rerr := io.ReadAll(data)
	if rerr != nil {
		str = rerr.Error()
	} else {
		str = string(bytes)
	}

	decoder := json.NewDecoder(strings.NewReader(str))
	var er AppError
	err := decoder.Decode(&er)
	if err != nil {
		return NewAppError("AppErrorFromJSON", "model.utils.decode_json.app_error", nil, "body: "+str, http.StatusInternalServerError).Wrap(err)
	}
	return &er
}

type AppError struct {
	Id            string `json:"id"`
	Message       string `json:"message"`               // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"`        // Internal error string to help the developer
	RequestId     string `json:"request_id,omitempty"`  // The RequestId that's also set in the header
	StatusCode    int    `json:"status_code,omitempty"` // The http status code
	Where         string `json:"-"`                     // The function where it happened in the form of Struct.Func
	IsOAuth       bool   `json:"is_oauth,omitempty"`    // Whether the error is OAuth specific
	params        map[string]any
	wrapped       error
}

func (er *AppError) Error() string {
	var sb strings.Builder

	// render the error information
	sb.WriteString(er.Where)
	sb.WriteString(": ")
	sb.WriteString(er.Message)

	// only render the detailed error when it's present
	if er.DetailedError != "" {
		sb.WriteString(", ")
		sb.WriteString(er.DetailedError)
	}

	// render the wrapped error
	err := er.wrapped
	if err != nil {
		sb.WriteString(", ")
		sb.WriteString(err.Error())
	}

	return sb.String()
}
func NewAppError(where string, id string, params map[string]any, details string, status int) *AppError {
	ap := &AppError{
		Id:            id,
		params:        params,
		Message:       id,
		Where:         where,
		DetailedError: details,
		StatusCode:    status,
		IsOAuth:       false,
	}
	return ap
}
func (er *AppError) Unwrap() error {
	return er.wrapped
}

func (er *AppError) Wrap(err error) *AppError {
	er.wrapped = err
	return er
}
