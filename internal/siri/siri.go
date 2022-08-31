package siri

import (
	"fmt"
	"net/http"
)

type RemoteError struct {
	Err error
	Loc string
}

func (e *RemoteError) Error() string {
	return fmt.Sprintf("%s: %s", e.Loc, e.Err)
}

func SoapCall(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
