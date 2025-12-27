package main

import (
	"fmt"
	"os"

	"github.com/msmhq/msm/cmd/msm/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
