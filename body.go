package main

import (
	"bytes"
	"fmt"

	"github.com/m-mizutani/deepalert"
	"github.com/m-mizutani/deepalert-github/md"
)

func attrToContents(attr deepalert.Attribute) md.Contents {
	return md.Contents{
		md.ToLiteral(fmt.Sprintf("%s (", attr.Key)),
		md.ToCode(string(attr.Type)),
		md.ToLiteral("): "),
		md.ToCode(attr.Value),
	}
}

func buildSummary(report deepalert.Report) []md.Node {
	attrList := &md.List{}
	attrMap := make(map[string]struct{})
	for _, alert := range report.Alerts {
		for _, attr := range alert.Attributes {
			hash := attr.Hash()
			if _, ok := attrMap[hash]; !ok {
				attrList.Items = append(attrList.Items, md.ListItem{
					Content: attrToContents(attr),
				})
				attrMap[hash] = struct{}{}
			}
		}
	}

	nodes := []md.Node{
		&md.Heading{
			Level:   1,
			Content: md.ToLiteral("Summary"),
		},
		&md.List{
			Items: []md.ListItem{
				{Content: md.Contents{
					md.ToLiteral("Severity: "),
					md.ToBold(string(report.Result.Severity)),
				}},
				{Content: md.Contents{
					md.ToLiteral("Reason: "),
					md.ToLiteral(report.Result.Reason),
				}},
				{Content: md.Contents{
					md.ToLiteral("Detected by "),
					md.ToCode(report.Alerts[0].Detector),
				}},
				{Content: md.Contents{
					md.ToLiteral("Rule: "),
					md.ToCode(report.Alerts[0].RuleName),
				}},
			},
		},
		&md.Heading{
			Level:   2,
			Content: md.ToLiteral("Attributes"),
		},
		attrList,
	}

	return nodes
}

func buildHostInspections(hosts map[string][]deepalert.ReportHost,
	attrs map[string]*deepalert.Attribute) (nodes []md.Node) {

	if len(hosts) == 0 {
		return
	}

	for hv, contents := range hosts {
		attr, _ := attrs[hv]
		nodes = append(nodes, &md.Heading{
			Level:   2,
			Content: md.ToLiteral(fmt.Sprintf("Host: %s", attr.Value)),
		})

		// list := md.List{}
		for _, c := range contents {
			if c.IPAddr != nil {
			}
			if c.Country != nil {

			}
			if c.ASOwner != nil {

			}
			if c.RelatedDomains != nil {

			}
			if c.RelatedURLs != nil {

			}
			if c.RelatedMalware != nil {

			}
			if c.Activities != nil {

			}
			if c.UserName != nil {

			}
			if c.Owner != nil {

			}
			if c.OS != nil {

			}
			if c.MACAddr != nil {

			}
			if c.HostName != nil {

			}
			if c.Software != nil {

			}
		}

	}

	return
}

func buildUserInspections(users map[string][]deepalert.ReportUser,
	attrs map[string]*deepalert.Attribute) (nodes []md.Node) {
	return nil
}

func buildBinaryInspections(binaries map[string][]deepalert.ReportBinary,
	attrs map[string]*deepalert.Attribute) (nodes []md.Node) {
	return nil
}

func buildInspections(report deepalert.Report) []md.Node {
	reportMap, err := report.ExtractContents()
	if err != nil {
		Logger.WithError(err).WithField("report", report).
			Error("Fail to extract contents from report")
		return nil
	}

	nodes := []md.Node{
		&md.Heading{
			Level:   1,
			Content: md.ToLiteral("Inspection Reports"),
		},
	}

	nodes = append(nodes, buildHostInspections(reportMap.Hosts, reportMap.Attributes)...)
	nodes = append(nodes, buildUserInspections(reportMap.Users, reportMap.Attributes)...)
	nodes = append(nodes, buildBinaryInspections(reportMap.Binaries, reportMap.Attributes)...)

	return nodes
}

func buildAlerts(report deepalert.Report) []md.Node {
	nodes := []md.Node{
		&md.Heading{
			Level:   1,
			Content: md.ToLiteral("Detail of Alerts"),
		},
	}

	for _, alert := range report.Alerts {
		nodes = append(nodes, []md.Node{
			&md.List{
				Items: []md.ListItem{
					{Content: md.Contents{
						md.ToLiteral("Description: "),
						md.ToLiteral(alert.Description),
					}},
					{Content: md.Contents{
						md.ToLiteral("Detected at: "),
						md.ToCode(alert.Timestamp.String()),
					}},
				},
			},
			&md.Heading{Content: md.ToLiteral("Attributes"), Level: 2},
		}...)

		attrList := &md.List{}
		for _, attr := range alert.Attributes {
			attrList.Items = append(attrList.Items, md.ListItem{
				Content: attrToContents(attr),
			})
		}

		nodes = append(nodes, attrList)
		nodes = append(nodes, &md.HorizontalRules{})
	}

	return nodes
}

func reportToBody(report deepalert.Report) (*bytes.Buffer, error) {
	doc := &md.Document{}
	doc.Extend(buildSummary(report))
	doc.Extend(buildInspections(report))
	doc.Extend(buildAlerts(report))

	buf := new(bytes.Buffer)
	if err := doc.Render(buf); err != nil {
		return nil, nil
	}

	return buf, nil
}
