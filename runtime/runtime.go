// Package runtime provides the thin GoCraft plugin runtime client.
//
// A GoCraft plugin is a standalone binary that connects to the server over a
// Unix socket. The server sends events; the plugin returns verdicts. Commands
// are invoked by the server and answered synchronously.
//
// This package will be replaced by the official GoCraft Go SDK once published.
package runtime

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

// Event is a plugin event dispatched by the server.
type Event struct {
	Type   string          `json:"type"`
	Fields map[string]any  `json:"fields"`
	Perms  map[string]bool `json:"perms"`
}

// IntField returns the integer value of a named event field, or 0.
func (e *Event) IntField(name string) int {
	v, _ := e.Fields[name]
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	}
	return 0
}

// StringField returns the string value of a named event field, or "".
func (e *Event) StringField(name string) string {
	v, _ := e.Fields[name]
	s, _ := v.(string)
	return s
}

// HasPermission reports whether the acting player holds perm.
func (e *Event) HasPermission(perm string) bool {
	return e.Perms[perm]
}

// Verdict is the plugin's response to an event.
type Verdict struct {
	Cancelled bool `json:"cancelled"`
}

// Allow returns a non-cancelling verdict.
func Allow() Verdict { return Verdict{} }

// Deny returns a cancelling verdict.
func Deny() Verdict { return Verdict{Cancelled: true} }

// Invocation is a command invocation dispatched by the server.
type Invocation struct {
	Command string   `json:"command"`
	Sender  string   `json:"sender"`
	Args    []string `json:"args"`
	conn    net.Conn
}

// Reply sends a text response back to the invoking player.
func (inv *Invocation) Reply(msg string) error {
	return sendFrame(inv.conn, &envelope{Type: "command.reply", Payload: map[string]any{
		"message": msg,
	}})
}

// Handler is implemented by the plugin to handle events and commands.
type Handler interface {
	OnEvent(e *Event) Verdict
	OnCommand(inv *Invocation) error
}

// Run connects to the GoCraft runtime socket and drives the event loop.
// It blocks until the connection is closed or an unrecoverable error occurs.
func Run(socketPath string, h Handler) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("connecting to runtime socket: %w", err)
	}
	defer conn.Close()

	if err := handshake(conn); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	for {
		env, err := recvFrame(conn)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("receiving frame: %w", err)
		}

		switch env.Type {
		case "event":
			var e Event
			if err := remarshal(env.Payload, &e); err != nil {
				continue
			}
			v := h.OnEvent(&e)
			if err := sendFrame(conn, &envelope{Type: "verdict", Payload: v}); err != nil {
				return err
			}
		case "command":
			var inv Invocation
			if err := remarshal(env.Payload, &inv); err != nil {
				continue
			}
			inv.conn = conn
			if err := h.OnCommand(&inv); err != nil {
				_ = inv.Reply(fmt.Sprintf("Error: %v", err))
			}
		}
	}
}

type envelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func handshake(conn net.Conn) error {
	return sendFrame(conn, &envelope{Type: "hello", Payload: map[string]any{"api": 1}})
}

func sendFrame(conn net.Conn, env *envelope) error {
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	var length [4]byte
	binary.BigEndian.PutUint32(length[:], uint32(len(data)))
	if _, err := conn.Write(length[:]); err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

func recvFrame(conn net.Conn) (*envelope, error) {
	var length [4]byte
	if _, err := io.ReadFull(conn, length[:]); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(length[:])
	buf := make([]byte, n)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	var env envelope
	if err := json.Unmarshal(buf, &env); err != nil {
		return nil, err
	}
	return &env, nil
}

func remarshal(src, dst any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
