package main

import (
	"fmt"
	"strings"
)

var handleMap = map[string]func(client *IRCClient, words []string){
	"PING": handlePing,
	"004":  handle004,
	"353":  handleUserlist,
}

func handlePing(client *IRCClient, words []string) {
	fmt.Println(strings.Join(words, " "))
	responseText := fmt.Sprintf("PONG %s\r\n", words[1])
	fmt.Println(responseText)
	client.connection.Write([]byte(responseText))
	return
}

func handleUserlist(client *IRCClient, words []string) {
	return
}

func handle004(client *IRCClient, words []string) {
	client.established = true
}
