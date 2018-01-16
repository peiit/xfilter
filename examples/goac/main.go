package main

import (
	"fmt"

	"github.com/sniperkit/xfilter/backend/goac"
)

func goacScanner() {

	ac := goac.NewAhoCorasick()
	ac.AddPatterns("google", []string{"google", "angular", "googlecloudplatform", "googlechrome", "golang", "gwtproject", "zxing", "v8"})
	ac.AddPatterns("facebook", []string{"facebook", "facebookarchive", "boltsframework"})
	ac.AddPatterns("postgres", []string{"postgres", "postgresql"})
	ac.AddPatterns("elasticsearch", []string{"elastic", "elasticsearch"})
	ac.AddPatterns("mongodb", []string{"mongodb", "mongo"})
	ac.AddPatterns("zeromq", []string{"zeromq", "zmq", "0mq"})
	ac.AddPatterns("kubernetes", []string{"kubernetes", "k8s"})
	ac.AddPatterns("boilerplate", []string{"boilerplate", "seed"})
	ac.AddPatterns("phantom", []string{"phantom", "phantomjs"})
	ac.AddPatterns("twitter", []string{"twbs", "twitter", "bower", "flightjs"})
	ac.AddPatterns("microsoft", []string{"microsoft", "dotnet", "aspnet", "exceptionless", "mono", "winjs"})
	ac.Build()

	results := ac.Scan(content)
	fmt.Println("Matches: ")
	for _, result := range results {
		fmt.Println("match=", string([]rune(content)[result.Start:result.End+1]), ", group=", result.Group, ", start=", result.Start, ", end=", result.End+1)
	}

}

func main() {
	goacScanner()
}
