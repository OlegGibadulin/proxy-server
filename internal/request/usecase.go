package request

import (
	"net/http"

	"github.com/OlegGibadulin/proxy-server/internal/models"
)

type RequestUseCase interface {
	StoreRequest(request *http.Request) (*models.Request, error)
	GetAllRequests() ([]*models.Request, error)
	GetRequest(requestID uint64) (*models.Request, error)
}
