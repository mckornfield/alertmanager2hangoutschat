package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type alertStub struct {
}

func (a alertStub) Status() string {
	return "firing"
}

func (a alertStub) Labels() map[string]string {
	return map[string]string{
		"alertname": "This is my fight song",
		"site":      "this-is-a-site",
		"namespace": "this-is-a-namespace",
		"node":      "this-is-a-node",
		"instance":  "this-is-a-instance",
	}
}

type Pair struct {
	Name, Value string
}

func (a alertStub) Annotations() map[string]interface{} {
	return map[string]interface{}{
		"grafanaDashboardQueryParams": "you=true",
		"grafanaDashboard":            "kmlnfsdankjsdafnkj",
		"createIssueQueryParams":      "heck yes",
		"showMe":                      "the money",
		"SortedPairs": []Pair{
			{Name: "grafanaDashboardQueryParams", Value: "abc123dxfsdca"},
			{Name: "grafanaDashboard", Value: "you=true"},
			{Name: "aws cluster thingy", Value: "heck yes all the clusters"},
			{Name: "namespace", Value: "iris-idk"},
			{Name: "createIssueQueryParams", Value: "pid=10360&issuetype=3&reporter=pare-svc-acct&priority=4&customfield_13840=15245&customfield_12140=15844&customfield_10240=15846"},
		},
	}
}

func (a alertStub) Values() map[string]interface{} {
	return map[string]interface{}{
		"gchat": map[string]interface{}{
			"createIssueQueryParams": "pid=10360&issuetype=3&reporter=pare-svc-acct&priority=4&customfield_13840=15245&customfield_12140=15844&customfield_10240=15846&customfield_10441=INFRA-12320",
			"alertLabels": map[string]string{
				"test": "test123",
			},
			"grafanaHost": "grafana.org",
		},
		"prometheus": map[string]interface{}{
			"alertmanager": map[string]interface{}{
				"ingress": map[string]interface{}{
					"hosts": []string{
						"create-issue.dev.com",
					},
				},
			},
		},
	}
}
func (a alertStub) gchat() string {
	return ""
}

var templateString = flag.String("template-string", "", "template for the messages sent to hangouts chat")
var filePathString = flag.String("file-path", "", "path to the outer template")

var defaultFuncs = template.FuncMap{
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title":   strings.Title,
	"now":     time.Now,
	// Usage example for grafana timestamps representing 1hr ago:
	// {{ (now.Add (hour -1)).Unix }}000
	"hour": func(h int32) time.Duration {
		timeString := fmt.Sprintf("%dh", h)
		duration, _ := time.ParseDuration(timeString)
		return duration
	},
	"minute": func(m int32) time.Duration {
		timeString := fmt.Sprintf("%dm", m)
		duration, _ := time.ParseDuration(timeString)
		return duration
	},
	// join is equal to strings.Join but inverts the argument order
	// for easier pipelining in templates.
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
}

func generateTemplate(s string, data interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(defaultFuncs).Parse(s)
	if err != nil {
		return "", err
	}
	var to bytes.Buffer
	err = tmpl.Execute(&to, data)
	return to.String(), err
}

func main() {
	flag.Parse()
	fmt.Println(*filePathString)
	fmt.Println(*templateString)
	templString, err := ioutil.ReadFile(*filePathString)

	if err != nil {
		fmt.Println("File error", err)
	}
	alertData := alertStub{}
	if err != nil {
		fmt.Println("Error with parse json", err)
		return
	}
	fmt.Println(fmt.Printf("%v", alertData))
	res1, err1 := generateTemplate(string(templString), alertData)

	if err1 != nil {
		fmt.Println("Error with parse 1", err1)
		return
	}
	fmt.Println(res1)
	res2, err2 := generateTemplate(res1, alertData)

	if err2 != nil {
		fmt.Println("Error with parse 2", err2)
		return
	}
	fmt.Println(res2)
}
