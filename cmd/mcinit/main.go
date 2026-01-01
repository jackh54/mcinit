package main

import (
	"os"

	"github.com/jackh54/mcinit/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
