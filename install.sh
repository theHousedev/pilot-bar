#!/usr/bin/env bash

go build -o ${HOME}/.local/bin/pilot-bar ./cmd/waybar/*
go build -o ${HOME}/.local/bin/pilot-bar-daemon ./cmd/daemon/*
