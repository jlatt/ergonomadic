package irc

import (
	"bufio"
	"github.com/jlatt/ergonomadic/irc/debug"
	"io"
	"log"
	"net"
	"strings"
)

const (
	R   = '→'
	W   = '←'
	EOF = ""
)

type Socket struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewSocket(conn net.Conn, commands chan<- editableCommand) *Socket {
	socket := &Socket{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}

	go socket.readLines(commands)

	return socket
}

func (socket *Socket) String() string {
	return socket.conn.RemoteAddr().String()
}

func (socket *Socket) Close() {
	socket.conn.Close()
	if debug.Net {
		log.Printf("%s closed", socket)
	}
}

func (socket *Socket) readLines(commands chan<- editableCommand) {
	commands <- &ProxyCommand{
		hostname: AddrLookupHostname(socket.conn.RemoteAddr()),
	}

	for {
		line, err := socket.reader.ReadString('\n')
		if socket.isError(err, R) {
			break
		}
		line = strings.TrimRight(line, CRLF)
		if len(line) == 0 {
			continue
		}
		if debug.Net {
			log.Printf("%s → %s", socket, line)
		}

		msg, err := ParseCommand(line)
		if err != nil {
			// TODO error messaging to client
			continue
		}
		commands <- msg
	}

	commands <- &QuitCommand{
		message: "connection closed",
	}
}

func (socket *Socket) Write(line string) (err error) {
	if _, err = socket.writer.WriteString(line); socket.isError(err, W) {
		return
	}

	if _, err = socket.writer.WriteString(CRLF); socket.isError(err, W) {
		return
	}

	if err = socket.writer.Flush(); socket.isError(err, W) {
		return
	}

	if debug.Net {
		log.Printf("%s ← %s", socket, line)
	}
	return
}

func (socket *Socket) isError(err error, dir rune) bool {
	if err != nil {
		if debug.Net && (err != io.EOF) {
			log.Printf("%s %c error: %s", socket, dir, err)
		}
		return true
	}
	return false
}
