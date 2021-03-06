package exactonline_bq

import (
	eo "github.com/leapforce-libraries/go_exactonline_new"
	el "github.com/leapforce-libraries/go_exactonline_new/logistics"
)

type Service struct {
	exactOnlineService *eo.Service
}

func NewService(exactOnlineService *eo.Service) *Service {
	return &Service{exactOnlineService}
}

func (service *Service) LogisticsService() *el.Service {
	return service.exactOnlineService.LogisticsService
}
