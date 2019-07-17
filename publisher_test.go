package main_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m-mizutani/deepalert"
	main "github.com/m-mizutani/deepalert-github"
)

type testConfig struct {
	GithubToken    string
	GithubEndpoint string
	GithubRepo     string
}

var config testConfig

func init() {
	confPath := "test.json"
	raw, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal("Fail to read config file", err)
	}

	if err := json.Unmarshal(raw, &config); err != nil {
		log.Fatal("Fail to unmarshal config file", err)
	}
}

func TestPublish(t *testing.T) {
	settings := main.GithubSettings{
		GithubToken:    config.GithubToken,
		GithubEndpoint: config.GithubEndpoint,
		GithubRepo:     config.GithubRepo,
	}

	var report deepalert.Report
	err := main.PublishReport(report, settings)
	assert.NoError(t, err)
}
