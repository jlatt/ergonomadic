package irc

import (
	"bufio"
	"io"
	"net"
)

const (
	R = '→'
	W = '←'
)

type IRCConn struct {
	closed  bool
	conn    net.Conn
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func NewIRCConn(conn net.Conn) *IRCConn {
	return &IRCConn{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
		writer:  bufio.NewWriter(conn),
	}
}

func (socket *IRCConn) String() string {
	return socket.conn.RemoteAddr().String()
}

func (socket *IRCConn) Close() {
	if socket.closed {
		return
	}
	socket.closed = true
	socket.conn.Close()
	Log.debug.Printf("%s closed", socket)
}

func (socket *IRCConn) Read() (line string, err error) {
	if socket.closed {
		err = io.EOF
		return
	}

	for socket.scanner.Scan() {
		line = socket.scanner.Text()
		if len(line) == 0 {
			continue
		}
		Log.debug.Printf("%s → %s", socket, line)
		return
	}

	err = socket.scanner.Err()
	socket.isError(err, R)
	if err == nil {
		err = io.EOF
	}
	return
}

func (socket *IRCConn) Write(line string) (err error) {
	if socket.closed {
		err = io.EOF
		return
	}

	if _, err = socket.writer.WriteString(line); socket.isError(err, W) {
		return
	}

	if _, err = socket.writer.WriteString(CRLF); socket.isError(err, W) {
		return
	}

	if err = socket.writer.Flush(); socket.isError(err, W) {
		return
	}

	Log.debug.Printf("%s ← %s", socket, line)
	return
}

func (socket *IRCConn) isError(err error, dir rune) bool {
	if err != nil {
		if err != io.EOF {
			Log.debug.Printf("%s %c error: %s", socket, dir, err)
		}
		return true
	}
	return false
}
