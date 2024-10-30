package server

import (
	"fmt"
)

func initialiseConfigDriver(opts Options) (ConfigDriver, error) {
	switch opts.ConfigDriver {
	case ConfigDriverType_File:
		return NewFileConfigDriver(opts)
	case ConfigDriverType_LibSQL:
		return NewLibSQLConfigDriver(opts)
	default:
		return nil, fmt.Errorf("unknown storage driver: %s", opts.ConfigDriver)
	}
}
