package main

import (
	"github.com/m-mizutani/deepalert"
)

type GithubSettings githubSettings

func PublishReport(report deepalert.Report, settings GithubSettings) error {
	return publishReport(report, githubSettings(settings))
}
