package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/emilrex/wt/internal/git"
	"github.com/emilrex/wt/internal/session"
)

// NewOptions contains options for the new command
type NewOptions struct {
	Name         string
	SourceBranch string
}

// RunNew creates a new worktree session and launches Claude Code
func RunNew(opts NewOptions) error {
	// Generate name if not provided
	name := opts.Name
	if name == "" {
		name = session.GenerateSessionName()
	}

	// Get source branch if not provided
	sourceBranch := opts.SourceBranch
	if sourceBranch == "" {
		var err error
		sourceBranch, err = git.GetCurrentBranch()
		if err != nil {
			return err
		}
	}

	// Get original repo root before creating session
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return err
	}

	// Create the session
	sess, err := session.Create(name, sourceBranch)
	if err != nil {
		return err
	}

	fmt.Printf("\nSession '%s' created successfully!\n", sess.Name)
	fmt.Printf("  Branch: %s\n", sess.Branch)
	fmt.Printf("  Path: %s\n", sess.Path)
	fmt.Println()

	// Launch Claude Code
	return launchClaude(sess.Path, repoRoot, false)
}

// launchClaude launches Claude Code in the specified directory
func launchClaude(worktreePath, repoRoot string, continueConversation bool) error {
	claudeArgs := "claude"

	if continueConversation {
		claudeArgs += " --continue"
	}

	// Add original repository as additional context
	if repoRoot != "" && repoRoot != worktreePath {
		claudeArgs += fmt.Sprintf(" --add-dir %q", repoRoot)
	}

	fmt.Printf("Launching Claude Code in %s...\n", worktreePath)

	// Use shell to run claude so that aliases work
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell, "-i", "-c", claudeArgs)
	cmd.Dir = worktreePath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
