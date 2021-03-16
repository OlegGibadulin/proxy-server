package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Request struct {
	ID      uint64
	Method  string
	URL     string
	Host    string
	Headers string
	Body    string
}

func ConvertRequestToModel(HTTPRequest *http.Request) (*Request, error) {
	headers, err := json.Marshal(HTTPRequest.Header)
	if err != nil {
		return nil, err
	}

	logrus.Info(HTTPRequest)
	logrus.Info(HTTPRequest.Body)

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, HTTPRequest.Body); err != nil {
		return nil, err
	}

	url := HTTPRequest.RequestURI
	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("https://%s%s", HTTPRequest.Host, url)
	}

	request := &Request{
		Method:  HTTPRequest.Method,
		URL:     url,
		Host:    HTTPRequest.Host,
		Headers: string(headers),
		Body:    body.String(),
	}
	return request, nil
}

func ConvertModelToRequest(request *Request) (*http.Request, error) {
	HTTPRequest, err := http.NewRequest(request.Method, request.URL, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	headers := make(map[string][]string)
	if err := json.Unmarshal([]byte(request.Headers), &headers); err != nil {
		return nil, err
	}
	for header, values := range headers {
		HTTPRequest.Header.Set(header, strings.Join(values, ","))
	}
	return HTTPRequest, nil
}
