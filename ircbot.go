package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
)

// IRCClient contains connection related variables
type IRCClient struct {
	connection  net.Conn
	config      Config
	established bool
}

// Connect establishes connection to IRC server
func (irc *IRCClient) Connect() {
	var err error
	irc.connection, err = net.Dial("tcp", irc.config.Adress)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(irc.connection, "NICK %s\r\n", irc.config.Nick)
	fmt.Fprintf(irc.connection, "USER %s 0 * %s\r\n", irc.config.Nick, irc.config.Nick)
	fmt.Fprintf(irc.connection, "PRIVMSG NickServ IDENTIFY %s\r\n", irc.config.Pass)

}

// Receive reads incoming stream and passes messages
// onto the parser
func (irc *IRCClient) Receive(messages chan string) {
	defer close(messages)
	scanner := bufio.NewScanner(bufio.NewReader(irc.connection))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		eolRune := []rune{'\r', '\n'}
		eol := []byte(string(eolRune))
		for i := 0; i < len(data)-1; i++ {
			if bytes.Equal(data[i:i+2], eol) {
				return i + 2, data[:i+1], nil
			}
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	for scanner.Scan() {
		text := scanner.Text()
		messages <- text
	}
}

// Parse incoming messages and direct them to the
// appropriate handler
func (irc *IRCClient) Parse(messages chan string) {
	for c := range messages {
		words := strings.Fields(c)
		var hostname string
		if c[0] == ':' {
			hostname = words[0][1:]
			words = words[1:]
		}

		if val, ok := handleMap[words[0]]; ok {
			val(irc, words)
		} else {
			fmt.Printf("[%s] %s\n", hostname, strings.Join(words, " "))
		}
	}
}

func main() {
	var server IRCClient
	server.config.Parse("serverconfig.yml")
	server.Connect()

	messages := make(chan string, 512)

	var wg sync.WaitGroup
	wg.Add(2)
	go server.Receive(messages)
	go server.Parse(messages)

	for {
		if server.established {
			break
		}
	}

	for _, channel := range server.config.Channels {
		fmt.Printf("Trying to join %s\n", channel)
		fmt.Fprintf(server.connection, "JOIN %s\r\n", channel)
	}

	wg.Wait()
}
