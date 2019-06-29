package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
)

type IRCConnection struct {
	connection  net.Conn
	established bool
}

func (irc *IRCConnection) Connect(config Config) {
	var err error
	irc.connection, err = net.Dial("tcp", config.Adress)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(irc.connection, "NICK %s\r\n", config.Nick)
	fmt.Fprintf(irc.connection, "USER %s 0 * %s\r\n", config.Nick, config.Nick)
	fmt.Fprintf(irc.connection, "PRIVMSG NickServ IDENTIFY %s\r\n", config.Pass)

}

func (irc *IRCConnection) Receive(EOC chan bool, messages chan string) {
	defer close(EOC)
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
	EOC <- true
}

func (irc *IRCConnection) Parse(messages chan string) {
	for c := range messages {
		words := strings.Fields(c)
		var hostname string
		if c[0] == ':' {
			hostname = words[0]
			words = words[1:]
		}
		if strings.EqualFold(words[0], "004") {
			irc.established = true
		}
		if strings.EqualFold(words[0], "PING") {
			fmt.Println(c)
			responseText := strings.Replace(c, "PING", "PONG", 1)
			irc.connection.Write([]byte(responseText))
		} else {
			fmt.Printf("[%s] %s \n", hostname, strings.Join(words, " "))
		}

	}
}

func main() {
	var config Config
	config.Parse("serverconfig.yml")

	var server IRCConnection
	server.Connect(config)

	eoc := make(chan bool)
	messages := make(chan string, 512)

	go server.Receive(eoc, messages)
	go server.Parse(messages)

	for {
		if server.established {
			break
		}
	}

	for _, channel := range config.Channels {
		fmt.Printf("Trying to join %s\n", channel)
		fmt.Fprintf(server.connection, "JOIN %s\r\n", channel)
	}

	for e := range eoc {
		if e {
			fmt.Println("Received exit")
		}
	}
}
