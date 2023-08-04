package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

const unknownErr = "unknown error"

var (
	ErrorDefaultAccountNotExist = errors.New("Default account of client is not exist ")
	ErrorProposalIDNotFound     = errors.New("Proposal ID not found ")
)

// ErrResponse define the information of the error response
type ErrResponse struct {
	ErrResponse ErrResponseJson
	StatusCode  int
}

type ErrResponseJson struct {
	Code    int32
	Message string
}

// Error returns the error msg
func (r ErrResponse) Error() string {
	return fmt.Sprintf("statusCode %v : code : %s  (Message: %s)",
		r.StatusCode, r.ErrResponse.Code, r.ErrResponse.Message)
}

// ConstructErrResponse  checks the response is an error response
func ConstructErrResponse(r *http.Response, bucketName, objectName string) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	/*
		if r == nil {
			return ErrResponse{
				StatusCode: r.StatusCode,
				ErrResponseJson{int32(404), "unknown err"},
			}
		}
	*/
	errResp := ErrResponse{}
	errResp.StatusCode = r.StatusCode

	// read err body of max 10M size
	const maxBodySize = 10 * 1024 * 1024
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		json := ErrResponseJson{Code: 404, Message: err.Error()}
		return ErrResponse{
			StatusCode:  r.StatusCode,
			ErrResponse: json,
		}
	}

	var resp ErrResponseJson
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Error().Msg("unmarshal err:" + err.Error())
	}
	errResp.ErrResponse = resp
	return errResp
}

// ToInvalidArgumentResp returns invalid argument response.
func ToInvalidArgumentResp(message string) error {
	return ErrResponse{
		StatusCode: http.StatusBadRequest,
		Code:       "InvalidArgument",
		Message:    message,
	}
}
