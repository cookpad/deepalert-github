package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/m-mizutani/deepalert"
)

var (
	// Logger can be modified by test code
	Logger = logrus.New()
)

type githubSettings struct {
	GithubToken    string `json:"github_token"`
	GithubEndpoint string `json:"github_endpoint"`
	GithubRepo     string `json:"github_repo"`
}

func getSecretValues(secretArn string, values interface{}) error {
	// sample: arn:aws:secretsmanager:ap-northeast-1:1234567890:secret:mytest
	arn := strings.Split(secretArn, ":")
	if len(arn) != 7 {
		return errors.New(fmt.Sprintf("Invalid SecretsManager ARN format: %s", secretArn))
	}
	region := arn[3]

	ssn := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	mgr := secretsmanager.New(ssn)

	result, err := mgr.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	})

	if err != nil {
		return errors.Wrapf(err, "Fail to retrieve secret values: %s", secretArn)
	}

	err = json.Unmarshal([]byte(*result.SecretString), values)
	if err != nil {
		return errors.Wrapf(err, "Fail to parse secret values as JSON: %s", secretArn)
	}

	return nil
}

func handler(ctx context.Context, snsEvent events.SNSEvent) error {
	Logger.WithField("event", snsEvent).Info("Start handler")

	var secrets githubSettings
	if err := getSecretValues(os.Getenv("SECRET_ARN"), &secrets); err != nil {
		return err
	}

	for _, record := range snsEvent.Records {
		var report deepalert.Report
		if err := json.Unmarshal([]byte(record.SNS.Message), &report); err != nil {
			return err
		}

		Logger.WithField("report", report).Info("publishing report")
		if issue, err := publishReport(report, secrets); err == nil {
			Logger.WithField("url", issue.GetHTMLURL()).Info("Issue created")
		} else {
			return errors.Wrap(err, "Fail to publish report to Github")
		}
	}

	return nil
}

func main() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)

	lambda.Start(handler)
}
