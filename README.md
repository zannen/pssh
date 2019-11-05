# pssh

`pssh` is a utility for parallel execution of commands over ssh.

It parallelises the initial connection, and the command(s) specified in each
`-command` argument. If _any_ of the commands exit with non-zero code, no
further commands are run.

Consider the difference between:

* `-command xx -command yy`

    This will parallelise `xx`, check that exit codes from all nodes were zero,
    then parallelise `yy`, then check that exit codes from all nodes were zero.
* `-command 'xx ; yy'`

    This will parallelise `xx ; yy`, thus ignoring the exit code of `xx` and
    only acting on the exit codes of `yy`.
* `-command 'xx && yy'`

    Standard shell rules apply. e.g. (if using `bash`) `yy` will only be run if
    `xx` exited with code zero.

## Run

```bash
pssh \
    -user remote_username \
    -key ~/.ssh/id_rsa \
    -hosts 'node0[1-3].somecluster.example.com:22' \
    -command hostname \
    -command 'echo "More commands"'
```

## Arguments

* `-command` (string(s)) See below.
* `-colour` (boolean) Produce colour output
* `-hosts`(string) List of hosts for ssh connections
    * Use "`a;b;c`" to specify three hosts, `a`, `b` and `c`
    * Use "`node[0-2,5].example.com`" to specify `node0.example.com`, `node1.example.com`, `node2.example.com` and `node5.example.com`
    * Use "`{london,paris,berlin}.example.com`" to specify `london.example.com`, `paris.example.com` and `berlin.example.com`
    * The above can be combined to a limited extent, e.g. use "`aaa;pre{bbb[1-3],ccc[5,8]}post;zzz`" to specify:
        * `aaa`
        * `prebbb1post`, `prebbb2post`, `prebbb3post`
        * `preccc5post`, `preccc8post`
        * `zzz`
* `-key` (string) Name of private key file, e.g. `$HOME/.ssh/id_rsa`
* `-verbose` (boolean) Produce verbose output
* `-user` (string) User name for ssh connections

### Commands

* Check for file contents, and abort if contents do not exist on _all_ hosts:

    ```bash
    pssh -user username -key /path/to/id_rsa \
        -hosts 'node0[0-9].cluster.example.com:22' \
        -command 'grep localhost /etc/hosts' \
        -command 'echo "Found localhost in /etc/hosts"'
    ```
* Check for file contents, and continue even if contents do not exist on
    _all_ hosts:

    ```bash
    pssh -user username -key /path/to/id_rsa \
        -hosts 'node0[0-9].cluster.example.com:22' \
        -command 'grep localhost /etc/hosts ; /bin/true' \
        -command 'echo "localhost may or may not be in /etc/hosts"'
    ```
* Complex scriptlet
    ```bash
    pssh -user username -key /path/to/id_rsa \
        -hosts 'node0[0-9].cluster.example.com:22' \
        -command 'echo "Your scriptlet here" | grep script ; echo "Multiple commands are ok."'
    ```
