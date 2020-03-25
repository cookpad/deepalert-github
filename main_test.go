package main_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"

	"github.com/m-mizutani/deepalert"
	. "github.com/m-mizutani/deepalert-github"
)

func loadJsonData(fpath string, data interface{}) {
	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(raw, data); err != nil {
		log.Fatal(err)
	}
}

func TestEventData(t *testing.T) {
	fpath := os.Getenv("TEST_JSON_FILE")
	if fpath == "" {
		t.Skip("No TEST_JSON_FILE")
	}

	var event events.SNSEvent
	loadJsonData(fpath, &event)

	var report *deepalert.Report
	defer FixPublish()
	InjectPublish(func(r deepalert.Report) { report = &r })

	InjectGetSecretValue(func(secretArn string, values interface{}) error { return nil })
	defer FixGetSecretValue()

	err := Handler("dummyARN", event)
	require.NoError(t, err)
	assert.NotEqual(t, "", report.ID)
}
