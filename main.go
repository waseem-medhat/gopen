package main

import (
    "os"
	"os/exec"
)

func main() {
    cmd := exec.Command("vi") // or absolute binary path
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        panic(err)
    }
}
