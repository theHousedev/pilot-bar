#!/usr/bin/env bash

go build -o /home/kh/.local/bin/pilot-bar ./cmd/waybar/*
go build -o /home/kh/.local/bin/pilot-bar-daemon ./cmd/daemon/*
