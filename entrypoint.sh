#!/bin/sh
set -e

# Accepts 'args' input from action.yml
go run ./.action-tmp/generate_todo_md.go "$INPUT_ARGS"
