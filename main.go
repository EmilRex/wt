package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/emilrex/wt/internal/cmd"
)

var version = "dev"

const usage = `wt - Manage isolated git worktrees for parallel Claude Code sessions

Usage:
  wt <command> [arguments]

Commands:
  new [name] [-b branch]  Create a new worktree session and launch Claude Code
  fg <session-name>       Resume an existing session (foreground)
  ls                      List all active sessions
  rm <session-name>       Remove a session
  rm -a|--all             Remove all sessions
  cd <session-name>       Open a shell in a session's worktree

Examples:
  wt new                       # New session with auto-generated name
  wt new auth-feature          # New session named 'auth-feature'
  wt new hotfix -b main        # New session from main branch
  wt fg auth-feature           # Resume the auth-feature session
  wt ls                        # List all sessions
  wt rm auth-feature           # Remove specific session
  wt rm --all                  # Remove all sessions
  wt cd auth-feature           # Open shell in session directory
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "new":
		runNew(os.Args[2:])
	case "fg":
		runFg(os.Args[2:])
	case "ls":
		runLs()
	case "rm":
		runRm(os.Args[2:])
	case "cd":
		runCd(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Print(usage)
	case "-v", "--version", "version":
		fmt.Println(version)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		fmt.Print(usage)
		os.Exit(1)
	}
}

func runNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	branch := fs.String("b", "", "Source branch to create worktree from")
	fs.StringVar(branch, "branch", "", "Source branch to create worktree from")
	_ = fs.Parse(args) // ExitOnError handles errors

	opts := cmd.NewOptions{
		SourceBranch: *branch,
	}

	// First non-flag argument is the session name
	if fs.NArg() > 0 {
		opts.Name = fs.Arg(0)
	}

	if err := cmd.RunNew(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runFg(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: session name required")
		fmt.Fprintln(os.Stderr, "Usage: wt fg <session-name>")
		os.Exit(1)
	}

	if err := cmd.RunFg(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runLs() {
	if err := cmd.RunLs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runRm(args []string) {
	fs := flag.NewFlagSet("rm", flag.ExitOnError)
	all := fs.Bool("a", false, "Remove all sessions")
	fs.BoolVar(all, "all", false, "Remove all sessions")
	_ = fs.Parse(args) // ExitOnError handles errors

	opts := cmd.RmOptions{
		All: *all,
	}

	if fs.NArg() > 0 {
		opts.SessionName = fs.Arg(0)
	}

	if !opts.All && opts.SessionName == "" {
		fmt.Fprintln(os.Stderr, "Error: session name required (or use --all)")
		fmt.Fprintln(os.Stderr, "Usage: wt rm <session-name>")
		fmt.Fprintln(os.Stderr, "       wt rm -a|--all")
		os.Exit(1)
	}

	if err := cmd.RunRm(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCd(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: session name required")
		fmt.Fprintln(os.Stderr, "Usage: wt cd <session-name>")
		os.Exit(1)
	}

	if err := cmd.RunCd(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
