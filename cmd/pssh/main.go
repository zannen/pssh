package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zannen/pssh"
	"github.com/zannen/pssh/expand"

	"golang.org/x/crypto/ssh"
)

type stringSlice []string

func (i *stringSlice) String() string {
	return "[strings]"
}

func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var commands stringSlice
	var keyFilename, sshUsername, hosts string
	flag.StringVar(&keyFilename, "key", "", "Name of private key file")
	flag.StringVar(&sshUsername, "user", "", "User name for ssh connections")
	flag.StringVar(&hosts, "hosts", "", "Hosts for ssh connections")
	flag.Var(&commands, "command", "Command(s) to run on hosts")
	flag.Parse()

	key, err := ioutil.ReadFile(keyFilename)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := ssh.ClientConfig{
		User:            sshUsername,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ssh.FixedHostKey(hostKey),
	}

	mcon := pssh.NewMultipleConnection()

	hostsList, err := expand.Expand(hosts)
	if err != nil {
		log.Fatalf("unable to parse hosts list: %v", err)
	}
	for _, host := range hostsList {
		mcon.Add(&config, "tcp", host)
	}

	code := 0
	for _, cmd := range commands {
		mcon.Command(cmd)
		rlist := mcon.Response()
		for _, r := range rlist {
			if r.StdOut != "" {
				nl := "\n"
				if r.StdOut[len(r.StdOut)-1:] == "\n" {
					nl = ""
				}
				fmt.Printf("%s [out]: %s%s", r.HostPort, r.StdOut, nl)
			}
			if r.StdErr != "" {
				nl := "\n"
				if r.StdErr[len(r.StdErr)-1:] == "\n" {
					nl = ""
				}
				fmt.Printf("%s [err]: %s%s", r.HostPort, r.StdErr, nl)
			}
			if r.Err != nil {
				fmt.Printf("%s [Error]: %v\n", r.HostPort, r.Err)
				code = 1
			}
		}
		if code > 0 {
			break
		}
	}
	os.Exit(code)
}
