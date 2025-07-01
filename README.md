# Command Tree
Group similar commands in a labled tree. Example uses cases:
- SSH/SFTP connections to different servers
- VPNs for different environments
- Hard to remember commands with similar functions, like most `tar` commands

## How to use

Define your command tree in a yaml file

```yaml
name: "Title"
branches:
  - name: "First" # This branch will execute its command and exit
    command: "echo First"
  - name: "Second" # This branch will execute its command and present the next choices
    command: "echo Second"
    branches:
    - name: "X" 
      command: "echo Second -> X"
  - name: "Third" # This branch will only present the next choices
    branches:
    - name: "one"
      command: "echo Third -> one"
    - name: "two"
      command: "echo Third -> two"
    - name: "three"
      branches:
      - name: "A"
        command: "echo Third -> two -> A"
      - name: "B"
        command: "echo Third -> two -> B"
      - name: "C" # This branch is invalid, does not have a command or child branches

```


Run the program specifying the configuration's location

```bash
command-tree /path/to/conf
```