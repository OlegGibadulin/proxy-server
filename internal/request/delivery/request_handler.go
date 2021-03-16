package delivery

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/OlegGibadulin/proxy-server/internal/helpers/network"
	"github.com/OlegGibadulin/proxy-server/internal/models"
	"github.com/OlegGibadulin/proxy-server/internal/request"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type RequestHandler struct {
	requestUseCase request.RequestUseCase
}

func NewRequestHandler(requestUseCase request.RequestUseCase) *RequestHandler {
	return &RequestHandler{
		requestUseCase: requestUseCase,
	}
}

func (rh *RequestHandler) Configure(e *echo.Echo) {
	e.GET("/requests", rh.GetAllRequests())
	e.GET("/requests/:id", rh.GetRequest())
	e.GET("/repeat/:id", rh.RepeatRequest())
	e.GET("/scan/:id", rh.ScanRequest())
}

func (rh *RequestHandler) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodConnect {
		HTTPSHandler := network.NewHTTPSHandler(request)

		reqToStore, err := HTTPSHandler.Handle(writer)
		if err != nil {
			logrus.Error(err)
			return
		}

		if _, err := rh.requestUseCase.StoreRequest(reqToStore); err != nil {
			logrus.Error(err)
			return
		}
	} else {
		if _, err := rh.requestUseCase.StoreRequest(request); err != nil {
			logrus.Error(err)
			return
		}

		resp, err := network.HandleHTTPRequest(request)
		if err != nil {
			logrus.Error(err)
			return
		}

		if _, err = io.Copy(writer, strings.NewReader(resp)); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func (rh *RequestHandler) GetAllRequests() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		requests, err := rh.requestUseCase.GetAllRequests()
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}
		return cntx.JSON(http.StatusOK, requests)
	}
}

func (rh *RequestHandler) GetRequest() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		requestID, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		request, err := rh.requestUseCase.GetRequest(requestID)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}
		return cntx.JSON(http.StatusOK, request)
	}
}

func (rh *RequestHandler) RepeatRequest() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		requestID, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		requestModel, err := rh.requestUseCase.GetRequest(requestID)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		request, err := models.ConvertModelToRequest(requestModel)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		resp, err := network.HandleHTTPRequest(request)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}
		return cntx.String(http.StatusOK, resp)
	}
}

func (rh *RequestHandler) ScanRequest() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		requestID, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		requestModel, err := rh.requestUseCase.GetRequest(requestID)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		request, err := models.ConvertModelToRequest(requestModel)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		realResp, err := network.HandleHTTPRequest(request)
		if err != nil {
			logrus.Error(err)
			return cntx.String(http.StatusInternalServerError, err.Error())
		}

		for _, sym := range [2]string{`'`, `"`} {
			newRequest := request

			if requestModel.Method == http.MethodGet {
				params := request.URL.Query()
				q := request.URL.Query()
				for key, values := range params {
					for _, value := range values {
						q.Add(key, fmt.Sprintf("%s%s", value, sym))
					}
				}
				newRequest.URL.RawQuery = q.Encode()
			}
			if requestModel.Method == http.MethodPost {
				if err := request.ParseForm(); err != nil {
					logrus.Error(err)
					return cntx.String(http.StatusInternalServerError, err.Error())
				}
				for key, values := range request.PostForm {
					for _, value := range values {
						newRequest.PostForm.Set(key, fmt.Sprintf("%s%s", value, sym))
					}
				}
			}

			curResp, err := network.HandleHTTPRequest(newRequest)
			if err != nil {
				logrus.Error(err)
				return cntx.String(http.StatusInternalServerError, err.Error())
			}

			if len(curResp) != len(realResp) {
				return cntx.String(http.StatusOK, "Request contains sql injection")
			}
		}

		return cntx.String(http.StatusOK, "Request does not contain sql injection")
	}
}
