package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

const (
	COLCMD = "cmd"
	COLERR = "error"
	COLRESET = "reset"
	COLSTDERR = "stderr"
	COLSTDOUT = "stdout"
)

func main() {
	var commands stringSlice
	var keyFilename, sshUsername, hosts string
	var colour, verbose bool
	flag.StringVar(&keyFilename, "key", "", "Name of private key file")
	flag.StringVar(&sshUsername, "user", "", "User name for ssh connections")
	flag.StringVar(&hosts, "hosts", "", "List of hosts for ssh connections")
	flag.BoolVar(&colour, "colour", false, "Produce colour output")
	flag.BoolVar(&verbose, "verbose", false, "Produce verbose output")
	flag.Var(&commands, "command", "Command(s) to run on hosts")
	flag.Parse()

	col := make(map[string]string)
	if colour {
		col[COLCMD] = "\033[36;1m" // cyan bold
		col[COLERR] = "\033[31;1;7m" // red bold inverse
		col[COLRESET] = "\033[0m" // reset
		col[COLSTDERR] = "\033[33;1m" // yellow bold
		col[COLSTDOUT] = "\033[32;1m" // green bold
	} else {
		col[COLCMD] = ""
		col[COLERR] = ""
		col[COLRESET] = ""
		col[COLSTDERR] = ""
		col[COLSTDOUT] = ""
	}

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

	mcon := NewMultipleConnection()

	hostsList, err := expand.Expand(hosts)
	if err != nil {
		log.Fatalf("unable to parse hosts list: %v", err)
	}
	for _, host := range hostsList {
		mcon.Add(&config, "tcp", host)
	}

	code := 0
	for _, cmd := range commands {
		if verbose {
			fmt.Printf("%sCommand%s: %s\n", col[COLCMD], col[COLRESET], cmd)
		}
		mcon.Command(cmd)
		rlist := mcon.Response()
		for _, r := range rlist {
			if r.StdOut != "" {
				nl := "\n"
				if r.StdOut[len(r.StdOut)-1:] == "\n" {
					nl = ""
				}
				fmt.Printf("%s %s[out]%s: %s%s", r.HostPort, col[COLSTDOUT], col[COLRESET], r.StdOut, nl)
			}
			if r.StdErr != "" {
				nl := "\n"
				if r.StdErr[len(r.StdErr)-1:] == "\n" {
					nl = ""
				}
				fmt.Printf("%s %s[err]%s: %s%s", r.HostPort, col[COLSTDERR], col[COLRESET], r.StdErr, nl)
			}
			if r.Err != nil {
				fmt.Printf("%s %s[Error]%s: %v\n", r.HostPort, col[COLERR], col[COLRESET], r.Err)
				code = 1
			}
		}
		if code > 0 {
			break
		}
	}
	os.Exit(code)
}
