package main

import (
	"github.com/google/go-github/v27/github"
	"github.com/m-mizutani/deepalert"
)

type GithubSettings githubSettings

func PublishReport(report deepalert.Report, settings GithubSettings) (*github.Issue, error) {
	return publishReport(report, githubSettings(settings))
}

var (
	NewGithubClient = newGithubClient
	ReportToBody    = reportToBody
)
