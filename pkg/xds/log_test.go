package xds

import (
	"log"
)

func ExampleLoggerFuncs() {
	logger := log.Logger{}

	xdsLogger := LoggerFuncs{
		DebugFunc: logger.Printf,
		InfoFunc:  logger.Printf,
		WarnFunc:  logger.Printf,
		ErrorFunc: logger.Printf,
	}

	xdsLogger.Debugf("debug")
	xdsLogger.Infof("info")
	xdsLogger.Warnf("warn")
	xdsLogger.Errorf("error")
}
