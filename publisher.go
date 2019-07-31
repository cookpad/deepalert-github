package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v27/github"
	"github.com/m-mizutani/deepalert"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

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

func reportToTitle(report deepalert.Report) string {
	return fmt.Sprintf("[%s] %s: %s", report.Alerts[0].Detector, report.Alerts[0].RuleName, report.Alerts[0].Description)
}

func publishReport(report deepalert.Report, settings githubSettings) (*github.Issue, error) {
	Logger.WithField("report", report).Info("Publishing report")

	client, err := newGithubClient(settings.GithubEndpoint, settings.GithubToken)
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
