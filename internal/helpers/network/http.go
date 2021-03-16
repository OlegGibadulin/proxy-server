package network

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func SendHTTPRequest(request *http.Request) (*http.Response, error) {
	newRequest, err := http.NewRequest(request.Method, request.URL.String(), request.Body)
	if err != nil {
		return nil, err
	}
	newRequest.Header = request.Header

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(newRequest)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func HandleHTTPRequest(request *http.Request) (string, error) {
	resp, err := SendHTTPRequest(request)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	decodedResponse, err := DecodeResponse(resp)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return string(decodedResponse), nil
}
