package request

import (
	"github.com/OlegGibadulin/proxy-server/internal/models"
)

type RequestRepository interface {
	InsertRequest(request *models.Request) error
	SelectByID(requestID uint64) (*models.Request, error)
	SelectLast() (*models.Request, error)
	SelectAll() ([]*models.Request, error)
}
