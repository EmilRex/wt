package cmd

import (
	"fmt"

	"github.com/emilrex/wt/internal/session"
)

// RmOptions contains options for the rm command
type RmOptions struct {
	SessionName string
	All         bool
}

// RunRm removes one or more sessions
func RunRm(opts RmOptions) error {
	if opts.All {
		fmt.Println("Removing all sessions...")
		return session.RemoveAll()
	}

	if opts.SessionName == "" {
		return fmt.Errorf("session name required (or use --all)")
	}

	return session.Remove(opts.SessionName)
}
