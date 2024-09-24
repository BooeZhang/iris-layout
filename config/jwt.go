package config

import (
	"time"
)

// Jwt JWT配置项.
type Jwt struct {
	Key        string        `json:"key"         mapstructure:"key"`
	Expired    time.Duration `json:"expired"     mapstructure:"expired"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}
