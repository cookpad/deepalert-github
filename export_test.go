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
	Handler         = handler
)

func InjectPublish(recv func(deepalert.Report)) {
	publish = func(report deepalert.Report, settings githubSettings) (*github.Issue, error) {
		recv(report)
		return nil, nil
	}
}
func FixPublish() { publish = publishToGithub }

func InjectGetSecretValue(f func(secretArn string, values interface{}) error) {
	getSecretValues = f
}
func FixGetSecretValue() { getSecretValues = awsGetSecretValues }
