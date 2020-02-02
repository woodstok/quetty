#!/usr/bin/env bash

if ! command -v quetty &> /dev/null ; then
	tmux display-message "quetty not found in path. Make sure it is installed in path"
	exit 1
fi

quetty -init -word
