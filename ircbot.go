package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type IRCConnection struct {
	connection net.Conn
}

func (irc *IRCConnection) Connect(config Config) {
	var err error
	irc.connection, err = net.Dial("tcp", config.Adress)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(irc.connection, "NICK %s\r\n", config.Nick)
	fmt.Fprintf(irc.connection, "USER %s 0 * %s\r\n", config.Nick, config.Nick)
}

func (irc *IRCConnection) Receive(ch chan bool) {
	scanner := bufio.NewScanner(bufio.NewReader(irc.connection))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		eolRune := []rune{'\r', '\n'}
		eol := []byte(string(eolRune))
		for i := 0; i < len(data)-1; i++ {
			if bytes.Equal(data[i:i+2], eol) {
				return i + 1, data[:i+1], nil
			}
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("%s\n", text)
	}
	ch <- false
}

func main() {
	var config Config
	config.Parse("serverconfig.yml")

	var server IRCConnection
	server.Connect(config)

	ch := make(chan bool)

	go server.Receive(ch)

	for {
		if !<-ch {
			break
		}
	}
}
