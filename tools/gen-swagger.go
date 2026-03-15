//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("swag", "init",
		"--dir", "./internal/server/http",
		"--generalInfo", "server.go",
		"--output", "./docs/swagger/server",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
