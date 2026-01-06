package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	Head   string
	Branch string
}

// GetRepoRoot returns the root directory of the current git repository
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRepoName returns the name of the repository (directory name)
func GetRepoName() (string, error) {
	root, err := GetRepoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Base(root), nil
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// FetchOrigin fetches from origin
func FetchOrigin() error {
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch from origin: %w", err)
	}
	return nil
}

// FastForwardBranch attempts to fast-forward the specified branch to origin
func FastForwardBranch(branch string) error {
	// Check if remote branch exists
	cmd := exec.Command("git", "rev-parse", "--verify", "origin/"+branch)
	if err := cmd.Run(); err != nil {
		// Remote branch doesn't exist, skip fast-forward
		return nil
	}

	// Get current branch to restore later
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	// If we're already on the branch, just pull
	if currentBranch == branch {
		cmd = exec.Command("git", "pull", "--ff-only")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to fast-forward %s: %w", branch, err)
		}
		return nil
	}

	// Otherwise, update the branch ref directly
	cmd = exec.Command("git", "fetch", "origin", fmt.Sprintf("%s:%s", branch, branch))
	if err := cmd.Run(); err != nil {
		// Branch might not be fast-forwardable, that's ok
		return nil
	}
	return nil
}

// BranchExists checks if a branch exists
func BranchExists(branch string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", branch)
	return cmd.Run() == nil
}

// CreateBranch creates a new branch from the source branch
func CreateBranch(name, source string) error {
	cmd := exec.Command("git", "branch", name, source)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s", name, string(output))
	}
	return nil
}

// DeleteBranch deletes a branch forcefully
func DeleteBranch(branch string) error {
	cmd := exec.Command("git", "branch", "-D", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete branch %s: %s", branch, string(output))
	}
	return nil
}

// HasCommits checks if a branch has at least one commit
func HasCommits(branch string) bool {
	cmd := exec.Command("git", "rev-parse", branch)
	return cmd.Run() == nil
}

// AddWorktree creates a new worktree at the specified path
func AddWorktree(path, branch string) error {
	cmd := exec.Command("git", "worktree", "add", path, branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	return nil
}

// RemoveWorktree removes a worktree forcefully
func RemoveWorktree(path string) error {
	cmd := exec.Command("git", "worktree", "remove", "--force", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %s", string(output))
	}
	return nil
}

// ListWorktrees returns all worktrees in porcelain format
func ListWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	var worktrees []Worktree
	var current Worktree

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
			}
			current = Worktree{}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "HEAD ") {
			current.Head = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch refs/heads/")
		}
	}

	// Don't forget the last one
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees, nil
}
