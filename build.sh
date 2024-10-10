#!/bin/bash

# Build the interactive_commands binary
mkdir -p build/interactive_commands
GOOS=linux GOARCH=amd64 go build -o build/interactive_commands/bootstrap ./functions/interactive_commands.lambda/main.go
chmod +x build/interactive_commands/bootstrap


# Build the update_commands binary
mkdir -p build/update_commands
GOOS=linux GOARCH=amd64 go build -o build/update_commands/bootstrap ./functions/update_commands.lambda/main.go
chmod +x build/update_commands/bootstrap
