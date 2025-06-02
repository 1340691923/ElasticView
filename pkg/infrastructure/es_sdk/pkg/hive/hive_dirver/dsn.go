package hive_dirver

import (
	gohive "github.com/philhuan/gohive-driver"
	"os/user"
	"strings"
)

type DSNConfig gohive.Config

func (c *DSNConfig) Complete() *DSNConfig {
	c.ColumnsWithoutTableName = true
	if c.User == "" {
		_user, _ := user.Current()
		if _user != nil {
			c.User = strings.Replace(_user.Name, " ", "", -1)
		}
		if c.User == "" {
			c.User = "default_user"
		}
	}
	// password may not matter but can't be empty
	if c.Passwd == "" {
		c.Passwd = "x"
	}
	if c.Auth == "" {
		c.Auth = "PLAIN"
	}
	return c
}

func ParseDSN(dsn string) (*DSNConfig, error) {
	config, err := gohive.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return (*DSNConfig)(config), nil
}

func (c *DSNConfig) FormatDSN() string {
	return (*gohive.Config)(c).FormatDSN()
}
