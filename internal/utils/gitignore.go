package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AddToGitignore adds a path to .gitignore
func AddToGitignore(repoRoot, pathToIgnore string) error {
	gitignorePath := filepath.Join(repoRoot, ".gitignore")

	// Read existing .gitignore if it exists
	var lines []string
	if PathExists(gitignorePath) {
		file, err := os.Open(gitignorePath)
		if err != nil {
			return fmt.Errorf("failed to open .gitignore: %w", err)
		}
		defer func() { _ = file.Close() }()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read .gitignore: %w", err)
		}
	}

	// Check if path already in .gitignore
	normalizedPath := strings.TrimPrefix(pathToIgnore, "./")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == normalizedPath || trimmed == "./"+normalizedPath {
			// Already present
			return nil
		}
	}

	// Add path to .gitignore
	lines = append(lines, "", "# mcinit server", normalizedPath)

	// Write back to file
	file, err := os.Create(gitignorePath)
	if err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	defer func() { _ = file.Close() }()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}

// FindGitRoot finds the git repository root directory
func FindGitRoot(startPath string) (string, error) {
	absPath, err := AbsolutePath(startPath)
	if err != nil {
		return "", err
	}

	currentPath := absPath
	for {
		gitPath := filepath.Join(currentPath, ".git")
		if PathExists(gitPath) {
			return currentPath, nil
		}

		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			// Reached root without finding .git
			return "", fmt.Errorf("not a git repository")
		}
		currentPath = parent
	}
}

// IsInGitRepo checks if a path is inside a git repository
func IsInGitRepo(path string) bool {
	_, err := FindGitRoot(path)
	return err == nil
}

