package main

import (
	"fmt"

	"github.com/m-mizutani/deepalert"
	"github.com/m-mizutani/deepalert-github/md"
)

func buildUserInspections(users map[string][]deepalert.ReportUser,
	attrs map[string]*deepalert.Attribute) (nodes []md.Node) {

	if len(users) == 0 {
		return
	}

	for hv, contents := range users {
		attr, _ := attrs[hv]
		merged := mergeReportUser(contents)

		nodes = append(nodes, &md.Heading{
			Level:   2,
			Content: md.ToLiteral(fmt.Sprintf("User: %s", attr.Value)),
		})

		nodes = append(nodes, buildActivitiesSection(merged.Activities)...)
	}

	return nil
}

func mergeReportUser(contents []deepalert.ReportUser) (merged deepalert.ReportUser) {
	for _, c := range contents {
		merged.Activities = append(merged.Activities, c.Activities...)
	}

	return
}
