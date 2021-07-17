A utility to help you control multiple computers ("nodes") from one. The node where you control them is called "root" and the nodes you control are called "children".

To run a root node, you should start the `rootd` daemon first then you can talk to the daemon with the client (like how docker works).
## Installation
Grab the latest binary for your system from [releases](https://github.com/singurty/fakework/releases/) and you're good to go.
## Tutorial
```
$ ./fakework root 8000 &
```
starts a root daemon on port 8000 on all ip addresses available
```
./fakework root 127.0.0.1 8000
```
starts a root daemon on 127.0.0.1:8000
```
./fakework child 127.0.0.1 8000
```
starts a child node connected to the root node listening at 127.0.0.1:8000
```
$ ./fakework add whoami --each
successfully added
```
runs the command `whoami` on all connected children
```
$ ./fakework show workload
1. Command: whoami Status: work successfully executed Output: singurty
```
shows the current workload. shows output of the command if it has been executed.
```
$ ./fakework log -f
2021/07/08 21:58:19 starting rpc server
2021/07/08 21:58:19 listening for children
2021/07/08 21:58:19 polling workload
2021/07/08 21:58:26 new child connected: 127.0.0.1:35328
2021/07/08 21:58:33 work added: {0 0 whoami false <nil> <nil> }
```
shows live logs
# Usage
```
Usage:
  fakework [command]

Available Commands:
  add         add work
  child       run a child node
  help        Help about any command
  log         view root daemon logs
  root        run a root node
  show        show something

Flags:
  -h, --help   help for fakework

Use "fakework [command] --help" for more information about a command.
```
