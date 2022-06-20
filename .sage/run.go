package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"go.einride.tech/sage/sg"
)

func Run(ctx context.Context) error {
	sg.Logger(ctx).Println("running...")

	cmd := exec.Command("go", "run", "./bin/futbot/main.go")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("failed to kill process: %w", err)
		}
	}()

	return cmd.Run()
}
