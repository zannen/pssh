package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type CmdResponse struct {
	HostPort string
	StdOut   string
	StdErr   string
	Err      error
}

type Connection struct {
	cfg      *ssh.ClientConfig
	client   *ssh.Client
	hostPort string
	network  string // tcp
	// session  *ssh.Session

	command  chan string
	response chan CmdResponse
}

func NewConnection(cfg *ssh.ClientConfig, network, hostPort string) *Connection {
	c := Connection{
		cfg:      cfg,
		hostPort: hostPort,
		network:  network,
		command:  make(chan string),
		response: make(chan CmdResponse),
	}

	go c.loop()

	return &c
}

func (c *Connection) dial() error {
	client, err := ssh.Dial(c.network, c.hostPort, c.cfg)
	if err != nil {
		return errors.Wrap(err, "Failed to dial")
	}
	c.client = client
	return nil
}

func (c *Connection) loop() {
	for {
		select {
		case cmd := <-c.command:
			for c.client == nil {
				err := c.dial()
				if err != nil {
					fmt.Printf("%s: Failed to dial: %v\n", c.hostPort, err)
					time.Sleep(time.Second)
				}
			}
			session, err := c.client.NewSession()
			if err != nil {
				c.response <- CmdResponse{
					HostPort: c.hostPort,
					Err:      errors.Wrap(err, "Failed to create session"),
				}
			} else {
				var o, e bytes.Buffer
				session.Stdout = &o
				session.Stderr = &e
				err = session.Run(cmd)
				c.response <- CmdResponse{
					HostPort: c.hostPort,
					StdOut:   o.String(),
					StdErr:   e.String(),
					Err:      err,
				}
			}

		default:
			time.Sleep(100 * time.Millisecond)

		}
	}
}

func (c *Connection) Command(cmd string) {
	c.command <- cmd
}

func (c *Connection) Response() CmdResponse {
	return <-c.response
}
