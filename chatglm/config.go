package chatglm

import "errors"

type Config struct {
	APIKey string
	Debug  bool
}

// Validate validate config
//   验证配置是否合法
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("APIKey is required")
	}
	return nil
}
