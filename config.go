package main

import (
	"os"
	"strconv"
)

type config struct {
	Port int64
	Logs bool
}

func (c *config) parse() error {
	const (
		portEnvVarName = "PORT"
		logsEnvVarName = "LOGS"
		portDefault    = 80
	)

	envPort := os.Getenv(portEnvVarName)
	if envPort != "" {
		var err error
		if c.Port, err = strconv.ParseInt(envPort, 10, 64); err != nil {
			return err
		}
	} else {
		c.Port = portDefault
	}

	envLogs := os.Getenv(logsEnvVarName)
	if envLogs != "" {
		var err error
		if c.Logs, err = strconv.ParseBool(envLogs); err != nil {
			return err
		}

		return nil
	}

	return nil
}
