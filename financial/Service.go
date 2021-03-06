package exactonline_bq

import (
	eo "github.com/leapforce-libraries/go_exactonline_new"
	ef "github.com/leapforce-libraries/go_exactonline_new/financial"
)

type Service struct {
	exactOnlineService *eo.Service
}

func NewService(exactService *eo.Service) *Service {
	return &Service{exactService}
}

func (service *Service) FinancialService() *ef.Service {
	return service.exactOnlineService.FinancialService
}
