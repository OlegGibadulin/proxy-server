package network

import (
	"bufio"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HTTPSHandler struct {
	connRequest *http.Request
	servRequest *http.Request
}

func NewHTTPSHandler(connRequest *http.Request) *HTTPSHandler {
	return &HTTPSHandler{
		connRequest: connRequest,
	}
}

func (hh *HTTPSHandler) Handle(writer http.ResponseWriter) (*http.Request, error) {
	hijackedConn, err := hh.hijackConnection(writer)
	if err != nil {
		return nil, err
	}
	defer hijackedConn.Close()

	url, _ := url.Parse(hh.connRequest.RequestURI)
	if err != nil {
		return nil, err
	}

	cert, err := GenerateCertificate(url)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}, ServerName: url.Scheme}

	clientConn, err := hh.getClientConnection(hijackedConn, config)
	if err != nil {
		return nil, err
	}
	defer clientConn.Close()

	serverConn, err := hh.getSererConnection(config)
	if err != nil {
		return nil, err
	}

	err = hh.sendRequest(clientConn, serverConn)
	if err != nil {
		return nil, err
	}
	return hh.servRequest, nil
}

func (hh *HTTPSHandler) hijackConnection(writer http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := writer.(http.Hijacker)
	if !ok {
		return nil, errors.New("Error in converting writer to http.Hijacker")
	}

	hijackedConn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}

	establishResp := []byte("HTTP/1.1 200 Connection established\r\n\r\n")
	if _, err = hijackedConn.Write(establishResp); err != nil {
		hijackedConn.Close()
		return nil, err
	}
	return hijackedConn, nil
}

func (hh *HTTPSHandler) getClientConnection(hijackedConn net.Conn, config *tls.Config) (*tls.Conn, error) {
	clientConn := tls.Server(hijackedConn, config)

	if err := clientConn.Handshake(); err != nil {
		clientConn.Close()
		return nil, err
	}

	connRequest, err := http.ReadRequest(bufio.NewReader(clientConn))
	if err != nil {
		clientConn.Close()
		return nil, err
	}
	hh.servRequest = connRequest

	return clientConn, nil
}

func (hh *HTTPSHandler) getSererConnection(config *tls.Config) (*tls.Conn, error) {
	serverConn, err := tls.Dial("tcp", hh.connRequest.Host, config)
	if err != nil {
		return nil, err
	}
	return serverConn, nil
}

func (hh *HTTPSHandler) sendRequest(clientConn *tls.Conn, serverConn *tls.Conn) error {
	dumpedRequest, err := httputil.DumpRequest(hh.servRequest, true)
	if err != nil {
		return err
	}
	if _, err = serverConn.Write(dumpedRequest); err != nil {
		return err
	}

	serverReader := bufio.NewReader(serverConn)
	resp, err := http.ReadResponse(serverReader, hh.servRequest)
	if err != nil {
		return err
	}

	decodedResp, err := DecodeResponse(resp)
	if err != nil {
		return err
	}
	if _, err = clientConn.Write(decodedResp); err != nil {
		return err
	}
	return nil
}
