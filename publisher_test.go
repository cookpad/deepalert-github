package main_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-mizutani/deepalert"
	main "github.com/m-mizutani/deepalert-github"
)

type testConfig struct {
	GithubToken      string
	GithubEndpoint   string
	GithubRepo       string
	GithubAppID      string
	GithubInstallID  string
	GithubPrivateKey string
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

	reportID := deepalert.ReportID(uuid.New().String())
	report := deepalert.Report{
		ID: reportID,
		Alerts: []deepalert.Alert{
			{
				Detector:    "blue",
				RuleName:    "orange",
				AlertKey:    "five",
				Description: "not sane",
				Timestamp:   time.Now(),
				Attributes: []deepalert.Attribute{
					{
						Type:    deepalert.TypeIPAddr,
						Key:     "source",
						Value:   "192.168.0.1",
						Context: []deepalert.AttrContext{deepalert.CtxRemote},
					},
				},
			},
		},
	}

	issue, err := main.Publish(report, settings)
	assert.NoError(t, err)
	assert.NotNil(t, issue)

	client, err := main.NewGithubClient(settings.GithubEndpoint, settings.GithubToken)
	require.NoError(t, err)
	ctx := context.Background()

	arr := strings.Split(settings.GithubRepo, "/")
	fetched, _, err := client.Issues.Get(ctx, arr[0], arr[1], *issue.Number)
	require.NoError(t, err)
	assert.Contains(t, *fetched.Title, "blue")
	assert.Contains(t, *fetched.Body, reportID, "192.168.0.1")
}
