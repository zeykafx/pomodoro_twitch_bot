package twitch_api_wrapper

import (
	"bufio"
	"net"
	"net/textproto"
)

type Client struct {
	Token    string
	Nick     string
	conn     net.Conn
	writer   *textproto.Writer
	reader   *textproto.Reader
	handlers map[string][]func(*Command) bool
}

func (client *Client) Connect(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	client.conn = conn
	client.writer = textproto.NewWriter(bufio.NewWriter(conn))
	client.reader = textproto.NewReader(bufio.NewReader(conn))
	return nil
}

func (client *Client) Auth() error {
	err := client.writer.PrintfLine("PASS %s", client.Token)
	if err != nil {
		return err
	}
	err = client.writer.PrintfLine("NICK %s", client.Nick)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) AddHandler(command string, f func(*Command) bool) {
	if client.handlers == nil {
		client.handlers = make(map[string][]func(*Command) bool)
	}
	handlers, ok := client.handlers[command]
	if !ok {
		client.handlers[command] = []func(*Command) bool{f}
	} else {
		client.handlers[command] = append(handlers, f)
	}
}

func (client *Client) Send(command *Command) error {
	err := client.writer.PrintfLine(command.Build())
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) CapReq(cap string) error {
	err := client.Send(&Command{
		Command: "CAP",
		Args:    []string{"REQ"},
		Suffix:  cap,
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Join(channel string) error {
	err := client.Send(&Command{
		Command: "JOIN",
		Args:    []string{channel},
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Close() {
	err := client.conn.Close()
	if err != nil {
		return
	}
}

func (client *Client) Handle() error {
	for {
		packet, err := client.reader.ReadLine()
		if err != nil {
			return err
		}

		command := ParsePacket(packet)
		handlers := client.handlers[command.Command]
		for _, handler := range handlers {
			if !handler(command) {
				return nil
			}
		}
	}
}
