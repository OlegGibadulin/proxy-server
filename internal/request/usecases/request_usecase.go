package usecases

import (
	"net/http"

	"github.com/OlegGibadulin/proxy-server/internal/models"
	"github.com/OlegGibadulin/proxy-server/internal/request"
)

type RequestUseCase struct {
	requestRepo request.RequestRepository
}

func NewRequestUseCase(repo request.RequestRepository) request.RequestUseCase {
	return &RequestUseCase{
		requestRepo: repo,
	}
}

func (ru *RequestUseCase) StoreRequest(request *http.Request) (*models.Request, error) {
	reqModel, err := models.ConvertRequestToModel(request)
	if err != nil {
		return nil, err
	}
	if err := ru.requestRepo.InsertRequest(reqModel); err != nil {
		return nil, err
	}
	return reqModel, nil
}

func (ru *RequestUseCase) GetAllRequests() ([]*models.Request, error) {
	requests, err := ru.requestRepo.SelectAll()
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (ru *RequestUseCase) GetRequest(requestID uint64) (*models.Request, error) {
	request, err := ru.requestRepo.SelectByID(requestID)
	if err != nil {
		return nil, err
	}
	return request, nil
}
