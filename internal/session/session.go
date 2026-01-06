package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/emilrex/wt/internal/git"
)

const (
	WorktreeBaseDir = ".wt"
	BranchPrefix    = "wt-"
)

// Session represents an isolated working environment
type Session struct {
	Name   string
	Branch string
	Path   string
}

// GetWorktreeBaseDir returns the base directory for all worktrees
func GetWorktreeBaseDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, WorktreeBaseDir), nil
}

// GenerateSessionName creates a timestamp-based session name
func GenerateSessionName() string {
	return time.Now().Format("20060102-150405")
}

// GetBranchName returns the branch name for a session
func GetBranchName(sessionName string) string {
	return BranchPrefix + sessionName
}

// GetSessionFromBranch extracts session name from branch name
func GetSessionFromBranch(branch string) string {
	return strings.TrimPrefix(branch, BranchPrefix)
}

// GetWorktreePath returns the worktree path for a session
func GetWorktreePath(repoName, sessionName string) (string, error) {
	baseDir, err := GetWorktreeBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, fmt.Sprintf("%s-%s", repoName, sessionName)), nil
}

// List returns all sessions for the current repository
func List() ([]Session, error) {
	repoName, err := git.GetRepoName()
	if err != nil {
		return nil, err
	}

	worktrees, err := git.ListWorktrees()
	if err != nil {
		return nil, err
	}

	baseDir, err := GetWorktreeBaseDir()
	if err != nil {
		return nil, err
	}

	var sessions []Session
	prefix := fmt.Sprintf("%s-", repoName)

	for _, wt := range worktrees {
		// Check if this worktree is in our base dir and matches our repo
		if !strings.HasPrefix(wt.Path, baseDir) {
			continue
		}

		dirName := filepath.Base(wt.Path)
		if !strings.HasPrefix(dirName, prefix) {
			continue
		}

		// Derive session name from directory, not branch
		// This makes sessions resilient to branch renames
		sessionName := strings.TrimPrefix(dirName, prefix)
		sessions = append(sessions, Session{
			Name:   sessionName,
			Branch: wt.Branch,
			Path:   wt.Path,
		})
	}

	return sessions, nil
}

// Find finds a session by name with partial matching support.
// If an exact match exists, it's returned. Otherwise, if exactly one
// session has a name starting with the query, that session is returned.
// If multiple sessions match, an error listing them is returned.
func Find(name string) (*Session, error) {
	sessions, err := List()
	if err != nil {
		return nil, err
	}

	// First, try exact match
	for _, s := range sessions {
		if s.Name == name {
			return &s, nil
		}
	}

	// Try prefix match
	var matches []Session
	for _, s := range sessions {
		if strings.HasPrefix(s.Name, name) {
			matches = append(matches, s)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("session '%s' not found", name)
	}

	if len(matches) == 1 {
		return &matches[0], nil
	}

	// Multiple matches - list them in the error
	var names []string
	for _, m := range matches {
		names = append(names, m.Name)
	}
	return nil, fmt.Errorf("'%s' matches multiple sessions: %s", name, strings.Join(names, ", "))
}

// Create creates a new session
func Create(name, sourceBranch string) (*Session, error) {
	repoName, err := git.GetRepoName()
	if err != nil {
		return nil, err
	}

	// Ensure base directory exists
	baseDir, err := GetWorktreeBaseDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create worktree base directory: %w", err)
	}

	branchName := GetBranchName(name)
	worktreePath, err := GetWorktreePath(repoName, name)
	if err != nil {
		return nil, err
	}

	// Check if worktree already exists
	if _, err := os.Stat(worktreePath); err == nil {
		return nil, fmt.Errorf("session '%s' already exists at %s", name, worktreePath)
	}

	// Fetch and fast-forward source branch
	fmt.Println("Fetching from origin...")
	if err := git.FetchOrigin(); err != nil {
		// Non-fatal: might not have a remote
		fmt.Printf("Warning: %v\n", err)
	}

	fmt.Printf("Updating %s...\n", sourceBranch)
	if err := git.FastForwardBranch(sourceBranch); err != nil {
		// Non-fatal: might not be fast-forwardable
		fmt.Printf("Warning: %v\n", err)
	}

	// Verify source branch has commits
	if !git.HasCommits(sourceBranch) {
		return nil, fmt.Errorf("source branch '%s' has no commits", sourceBranch)
	}

	// Create branch if it doesn't exist
	if !git.BranchExists(branchName) {
		fmt.Printf("Creating branch %s from %s...\n", branchName, sourceBranch)
		if err := git.CreateBranch(branchName, sourceBranch); err != nil {
			return nil, err
		}
	} else {
		fmt.Printf("Branch %s already exists, using existing branch\n", branchName)
	}

	// Create worktree
	fmt.Printf("Creating worktree at %s...\n", worktreePath)
	if err := git.AddWorktree(worktreePath, branchName); err != nil {
		// Clean up branch if we just created it
		_ = git.DeleteBranch(branchName)
		return nil, err
	}

	return &Session{
		Name:   name,
		Branch: branchName,
		Path:   worktreePath,
	}, nil
}

// Remove removes a session
func Remove(name string) error {
	session, err := Find(name)
	if err != nil {
		return err
	}

	fmt.Printf("Removing worktree %s...\n", session.Path)
	if err := git.RemoveWorktree(session.Path); err != nil {
		return err
	}

	fmt.Printf("Deleting branch %s...\n", session.Branch)
	if err := git.DeleteBranch(session.Branch); err != nil {
		// Non-fatal: branch might have been deleted already
		fmt.Printf("Warning: %v\n", err)
	}

	return nil
}

// RemoveAll removes all sessions for the current repository
func RemoveAll() error {
	sessions, err := List()
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions to remove")
		return nil
	}

	for _, s := range sessions {
		if err := Remove(s.Name); err != nil {
			fmt.Printf("Warning: failed to remove session '%s': %v\n", s.Name, err)
		}
	}

	return nil
}
