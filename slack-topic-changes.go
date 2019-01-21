package main

import (
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var count int
var debug bool

var msgRegex = regexp.MustCompile("set the channel's topic: (.+)$")

func init() {
	flag.IntVar(&count, "count", 100, "Number of results to return from Slack")
	flag.BoolVar(&debug, "debug", false, "Use debugging in Slack requests?")
}

func handleMsg(msg slack.SearchMessage, api *slack.Client, nameCache map[string]string) {

	if _, ok := nameCache[msg.User]; !ok {
		if debug {
			fmt.Printf("No entry for name %s in user cache, asking Slack\n", msg.User)
		}
		user, err := api.GetUserInfo(msg.User)
		if err != nil {
			fmt.Printf("Error looking up user %s: %s\n", msg.User, err)
			return
		}
		if debug {
			fmt.Printf("Adding name cache mapping for %s to user name %s\n", msg.User, user.Name)
		}
		nameCache[msg.User] = user.Name
	}
	userName := nameCache[msg.User]

	timeStrings := strings.Split(msg.Timestamp, ".")
	sec, err := strconv.ParseInt(timeStrings[0], 10, 64)
	if err != nil {
		fmt.Printf("Error converting timestamp %s (specifically %s) to int64: %s\n", msg.Timestamp, timeStrings[0], err)
		return
	}
	// For sake of convenience we outright ignore nsec here
	ts := time.Unix(sec, 0)

	regexMatches := msgRegex.FindSubmatch([]byte(msg.Text))
	if regexMatches != nil {
		fmt.Printf("%s: %s set topic to '%s'\n", ts, userName, regexMatches[1])
	} else {
		fmt.Printf("%s: %s emptied topic\n", ts, userName)
	}
}

func main() {

	flag.Parse()

	nameCache := make(map[string]string)

	// Originally I was creating a new HTTP client 'cuz I thought we'd have to do a manual OAuth process
	// in order to get a token... now I just do it for fun
	httpclient := &http.Client{}
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionHTTPClient(httpclient))
	slack.OptionDebug(debug)(api)

	params := slack.NewSearchParameters()
	params.Count = count
	queryString := fmt.Sprintf("\"set the channel's topic\" in:%s", os.Getenv("SLACK_CHANNEL"))
	msgs, err := api.SearchMessages(queryString, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, msg := range msgs.Matches {
		handleMsg(msg, api, nameCache)
	}
}
