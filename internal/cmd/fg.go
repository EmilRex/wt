package cmd

import (
	"fmt"

	"github.com/emilrex/wt/internal/session"
)

// RunFg resumes an existing session by launching Claude Code with --continue
func RunFg(sessionName string) error {
	sess, err := session.Find(sessionName)
	if err != nil {
		return err
	}

	fmt.Printf("Resuming session '%s'...\n", sess.Name)
	fmt.Printf("  Branch: %s\n", sess.Branch)
	fmt.Printf("  Path: %s\n", sess.Path)
	fmt.Println()

	// Launch Claude Code with --continue using shell for alias support
	return launchClaude(sess.Path, "", true)
}
