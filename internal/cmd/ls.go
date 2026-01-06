package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/emilrex/wt/internal/session"
)

// RunLs displays all active sessions for the current repository
func RunLs() error {
	sessions, err := session.List()
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("No active sessions")
		return nil
	}

	// Replace home directory with ~ for display
	home, _ := os.UserHomeDir()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Session\tBranch\tPath")
	fmt.Fprintln(w, "-------\t------\t----")

	for _, s := range sessions {
		displayPath := s.Path
		if home != "" && strings.HasPrefix(s.Path, home) {
			displayPath = "~" + strings.TrimPrefix(s.Path, home)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.Name, s.Branch, displayPath)
	}

	return w.Flush()
}
