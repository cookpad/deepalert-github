package main

import (
	"github.com/google/go-github/v27/github"
	"github.com/m-mizutani/deepalert"
)

type GithubSettings githubSettings

func Publish(report deepalert.Report, settings GithubSettings) (*github.Issue, error) {
	return publish(report, githubSettings(settings))
}

var (
	NewGithubClient = newGithubClient
	ReportToBody    = reportToBody
)
