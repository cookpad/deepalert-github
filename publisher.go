package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v27/github"
	"github.com/m-mizutani/deepalert"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type githubSettings struct {
	GithubToken      string `json:"github_token"`
	GithubEndpoint   string `json:"github_endpoint"`
	GithubRepo       string `json:"github_repo"`
	GithubAppID      string `json:"github_app_id"`
	GithubInstallID  string `json:"github_install_id"`
	GithubPrivateKey string `json:"github_private_key"`
}

func (x githubSettings) hasAppSettings() bool {
	return (x.GithubAppID != "" && x.GithubInstallID != "" && x.GithubPrivateKey != "")
}

func (x githubSettings) newClient() (*github.Client, error) {
	if x.hasAppSettings() {
		appID, err := strconv.ParseInt(x.GithubAppID, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to parse appID: %s", x.GithubAppID)
		}

		installID, err := strconv.ParseInt(x.GithubInstallID, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to parse InstallID: %s", x.GithubInstallID)
		}

		privateKey, err := base64.StdEncoding.DecodeString(x.GithubPrivateKey)
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to decode privateKey as base64, len:%d", len(x.GithubPrivateKey))
		}

		return newGithubAppClient(x.GithubEndpoint, int(appID), int(installID), privateKey)
	} else {
		return newGithubClient(x.GithubEndpoint, x.GithubToken)
	}
}

func newGithubClient(endpoint, token string) (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if !strings.HasSuffix(endpoint, "/") {
		endpoint = endpoint + "/"
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to parse GtihubEndpoint: %s", endpoint)
	}
	client.BaseURL = url

	return client, nil
}

func newGithubAppClient(endpoint string, appID int, installID int, privateKey []byte) (*github.Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.New(tr, appID, installID, privateKey)
	if !strings.HasSuffix(endpoint, "/") {
		endpoint = endpoint + "/"
	}

	if err != nil {
		return nil, errors.Wrap(err, "Fail to create GH client")
	}
	itr.BaseURL = endpoint + "app"

	client := github.NewClient(&http.Client{Transport: itr})

	if endpoint != "" {
		url, err := url.Parse(endpoint)
		if err != nil {
			return nil, errors.Wrapf(err, "Invalid URL format: %s", endpoint)
		}

		client.BaseURL = url
	}

	Logger.WithField("client", client).Trace("Github Client is created")

	return client, nil
}

func reportToTitle(report deepalert.Report) string {
	return fmt.Sprintf("[%s] %s: %s", report.Alerts[0].Detector, report.Alerts[0].RuleName, report.Alerts[0].Description)
}

func publishReport(report deepalert.Report, settings githubSettings) (*github.Issue, error) {
	Logger.WithField("report", report).Info("Publishing report")

	client, err := settings.newClient()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to create ")
	}

	title := reportToTitle(report)
	buf, err := reportToBody(report)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to build body")
	}
	body := buf.String()

	ctx := context.Background()
	issueReq := github.IssueRequest{
		Title: github.String(title),
		Body:  github.String(body),
	}
	arr := strings.Split(settings.GithubRepo, "/")
	if len(arr) != 2 {
		return nil, fmt.Errorf("%s is not repository format, must be {owner}/{repo_name}", settings.GithubRepo)
	}

	issue, resp, err := client.Issues.Create(ctx, arr[0], arr[1], &issueReq)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to create issue")
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("Fail to create issue because response code is not 201: %d", resp.StatusCode)
	}

	return issue, nil
}
