package repository

import (
	"github.com/OlegGibadulin/proxy-server/internal/models"
	"github.com/OlegGibadulin/proxy-server/internal/request"
	"github.com/jinzhu/gorm"
)

type RequestPgRepository struct {
	db *gorm.DB
}

func NewRequestPgRepository(conn *gorm.DB) request.RequestRepository {
	return &RequestPgRepository{
		db: conn,
	}
}

func (rr *RequestPgRepository) InsertRequest(request *models.Request) error {
	if err := rr.db.Create(request).Error; err != nil {
		return err
	}
	return nil
}

func (rr *RequestPgRepository) SelectByID(requestID uint64) (*models.Request, error) {
	request := &models.Request{}
	if err := rr.db.Where("id=?", requestID).First(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rr *RequestPgRepository) SelectLast() (*models.Request, error) {
	request := &models.Request{}
	if err := rr.db.Limit(1).Order("id desc").First(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (rr *RequestPgRepository) SelectAll() ([]*models.Request, error) {
	var requests []*models.Request
	if err := rr.db.Order("id desc").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}
