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

## Command examples

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
