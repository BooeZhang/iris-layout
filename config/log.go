package config

type Log struct {
	Formatter    string `json:"format" mapstructure:"format"`
	EnableCaller bool   `json:"enable-caller" mapstructure:"enable-caller"`
}
