package config

import (
	"fmt"
	"strings"
)

var (
	providers = make(map[string]NewConfigProviderFunc)
)

type NewConfigProviderFunc func() (ConfigurationProvider, error)

func RegisterConfigurationProvider(driverName string, fn NewConfigProviderFunc) {
	driverName = strings.TrimSpace(driverName)

	if len(driverName) == 0 {
		panic("driverName is empty")
	}

	if fn == nil {
		panic(driverName + "'s fn is nil")
	}

	_, exist := providers[driverName]

	if exist {
		panic(driverName + " already registered")
	}

	providers[driverName] = fn

	return
}

func NewConfigurationProvider(driverName string) (provider ConfigurationProvider, err error) {
	driverName = strings.TrimSpace(driverName)
	fn, exist := providers[driverName]

	if !exist {
		err = fmt.Errorf("%s's driver not exist", driverName)
		return
	}

	return fn()
}
