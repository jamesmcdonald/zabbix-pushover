package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type PushoverMessage struct {
	token   string
	user    string
	title   string
	message string
	devices string
}

func (msg *PushoverMessage) Send() error {
	data := url.Values{}
	data.Set("token", msg.token)
	data.Set("user", msg.user)
	data.Set("message", msg.message)
	if len(msg.title) > 0 {
		data.Set("title", msg.title)
	}
	if len(msg.devices) > 0 {
		data.Set("device", msg.devices)
	}
	_, err := http.PostForm("https://api.pushover.net/1/messages.json", data)
	return err
}

// The API token for the Pushover application
var token string

func init() {
	// TODO Make this a configurable path
	const (
		configfile = "/etc/zabbix/pushover.conf"
	)
	conf, err := os.Open(configfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: can't read config: %s\n", os.Args[0], err)
		os.Exit(1)
	}
	defer conf.Close()
	data, err := ioutil.ReadAll(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error reading config: %s\n", os.Args[0], err)
		os.Exit(1)
	}
	token = strings.TrimSpace(string(data))
}

func main() {
	usage := fmt.Sprintf("usage: %s <user>[:device[,device...]] <title> <body>\n", os.Args[0])
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}

	parts := strings.Split(os.Args[1], ":")
	if len(parts) > 2 {
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}
	user := parts[0]
	var devices string
	if len(parts) == 2 {
		devices = parts[1]
	}

	p := PushoverMessage{
		token,
		user,
		os.Args[2],
		os.Args[3],
		devices,
	}
	p.Send()
}
