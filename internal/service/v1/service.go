package v1

import (
	"irir-layout/config"
)

type ServiceContext struct {
	Cfg *config.Config
}

func NewServiceContext() *ServiceContext {
	c := config.GetConfig()
	return &ServiceContext{Cfg: c}
}
