/*
Wrapper over server
*/

package client

import (
	"net/url"
	"strings"
	"time"

	"github.com/Lasiar/pollsc/server"
)

var (
	add chan server.Client
)

// Message message for user
type Message struct {
	Text string
	ID   int
}

// Init init client
func Init() chan Message {
	var out chan server.Message

	outMessage := make(chan Message)
	add, _, out = server.Start()
	go middleChangeState(out, outMessage)
	return outMessage
}

func middleChangeState(out chan server.Message, outMessage chan Message) {
	for o := range out {
		outMessage <- Message{o.Text, o.Client.ClientID}
	}
}

// Processed main work client
func Processed(message string, id int) (string, error) {
	firstWord := message
	arguments := ""

	firstSpace := strings.Index(message, " ")

	if firstSpace > 0 {
		firstWord = message[:firstSpace]
		arguments = message[firstSpace:]
	}

	if len(message)+2 == firstSpace {
		return "", nil
	}

	if firstWord == "listen" {
		if err := addToListen(arguments, id); err != nil {
			return "", err
		}
		return "Added", nil
	}

	if firstWord == "get-all" {
		return server.GetInfo(id), nil
	}

	return "Не понял комнду", nil
}

func addToListen(args string, id int) error {
	var err error

	urls := strings.Split(strings.TrimSpace(args), ",")

	validUrls := make([]*url.URL, len(urls))

	for i, current := range urls {
		validUrls[i], err = url.Parse(current)

		if validUrls[i].Scheme == "" {
			return err
		}
		add <- server.Client{ClientID: id, URL: *validUrls[i], Timer: time.Duration(10 * time.Second)}
	}

	return nil
}
