package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/emilrex/wt/internal/session"
)

// RunCd opens an interactive shell in a session's worktree directory
func RunCd(sessionName string) error {
	sess, err := session.Find(sessionName)
	if err != nil {
		return err
	}

	fmt.Printf("Opening shell in session '%s' (%s)\n", sess.Name, sess.Path)
	fmt.Println("Type 'exit' to return to your original location")
	fmt.Println()

	// Get user's shell or fall back to /bin/bash
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell)
	cmd.Dir = sess.Path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment to show we're in a wt session
	cmd.Env = append(os.Environ(), fmt.Sprintf("WT_SESSION=%s", sess.Name))

	return cmd.Run()
}
