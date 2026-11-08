package main

import (
	"os"

	"github.com/jackh54/mcinit/internal/buildinfo"
	"github.com/jackh54/mcinit/internal/cli"
)

// version is stamped by GoReleaser (-X main.version=...).
var version = "dev"

func main() {
	buildinfo.Version = version
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
