A utility to help you manage multiple computers ("nodes") from one. The node where you managem them is called "root" aand the nodes you manage are called "children".
## Running a root node
```
./fakework root 8000
```
starts a root node on port 8000 on all ip addresses available
```
./fakework root 127.0.0.1 8000
```
starts a root node on 127.0.0.1:8000
## Running a child node
```
./fakework child 127.0.0.1 8000
```
starts a child node connected to the root node listening at 127.0.0.1:8000
# Usage
```
Usage:
  fakeroot [command]

Available Commands:
  add         add work
  child       run a child node
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  log         view root daemon logs
  root        run a root node
  show        show something

Flags:
  -h, --help   help for fakeroot

Use "fakeroot [command] --help" for more information about a command.
```
