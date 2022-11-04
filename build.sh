#!/bin/sh

go build -o google-font-localizer
GOOS=linux  GOARCH=amd64 go build -o google-font-localizer-linux-amd64
GOOS=linux  GOARCH=arm   go build -o google-font-localizer-linux-arm
GOOS=linux  GOARCH=arm64 go build -o google-font-localizer-linux-arm64
GOOS=darwin GOARCH=amd64 go build -o google-font-localizer-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o google-font-localizer-darwin-arm64
