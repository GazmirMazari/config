package config

import (
	"net/http"
	"strconv"
	"time"
)

type configFlag bool

type ClientConfig struct {
	Timeout            string `yaml:"Timeout"`
	IdleConnTimeout    string `yaml:"IdleConnTimeout"`
	MaxIdleConsPerHost string `yaml:"MaxIdleConsPerHost"`
	MaxConsPerHost     string `yaml:"MaxConsPerHost"`
	DisableCompression configFlag
	InsecureSkipVerify configFlag
}

func httpClient(cc ClientConfig) *http.Client {
	disableCompression := false
	timeout := 15

	if cc.DisableCompression == True {
		disableCompression = true
	}

	timeoutValue := toInt(cc.Timeout)
	if timeoutValue != 0 {
		timeout = timeoutValue
	}

	maxIdleConsPerHost, _ := strconv.Atoi(cc.MaxIdleConsPerHost)
	maxConsPerHost, _ := strconv.Atoi(cc.MaxConsPerHost)
	idleConnTimeout, _ := strconv.Atoi(cc.IdleConnTimeout)

	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
			MaxIdleConnsPerHost: maxIdleConsPerHost,
			MaxConnsPerHost:     maxConsPerHost,
			DisableCompression:  disableCompression,
		},
	}
}

func toInt(s string) int {
	value, _ := strconv.Atoi(s)
	return value
}
