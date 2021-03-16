package network

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func CopyHeaders(from, to http.Header) {
	for header, values := range from {
		to.Set(header, strings.Join(values, ","))
	}
}

func ConvertBodyToString(HTTPBody io.ReadCloser) (string, error) {
	body := new(bytes.Buffer)
	if _, err := io.Copy(body, HTTPBody); err != nil {
		return "", err
	}
	return body.String(), nil
}

func getHeadersStr(response *http.Response) string {
	var headersStr string
	for header, values := range response.Header {
		for _, value := range values {
			headersStr += fmt.Sprintf("%s:%s\n", header, value)
		}
	}
	return headersStr
}

func DecodeResponse(resp *http.Response) ([]byte, error) {
	headers := getHeadersStr(resp)
	respTop := fmt.Sprintf("%s\n%s\n%s\n", resp.Status, resp.Proto, headers)

	body := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, _ = gzip.NewReader(body)
	}
	defer body.Close()

	bodyContent, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	wholeResp := append([]byte(respTop), bodyContent...)
	return wholeResp, nil
}
