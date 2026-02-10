package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const staleDays = 30

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ".cache", "cht", "topics"), nil
}

// ensureCache returns the path to the cached topic list,
// fetching it if missing. Prints a staleness hint but never blocks.
func ensureCache() (string, error) {
	path, err := cachePath()
	if err != nil {
		return "", err
	}

	info, statErr := os.Stat(path)
	if statErr != nil {
		return path, fetchTopicList(path)
	}

	age := time.Since(info.ModTime())
	if age > staleDays*24*time.Hour {
		fmt.Fprintf(os.Stderr, "cht: topic cache is %d days old. Run cht --refresh\n", int(age.Hours()/24))
	}
	return path, nil
}

func refreshCache() error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	return fetchTopicList(path)
}

func fetchTopicList(dest string) error {
	req, err := http.NewRequest("GET", "https://cheat.sh/:list", nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "curl/8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch topic list: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cheat.sh/:list returned %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create cache file: %w", err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		_ = f.Close()
		return fmt.Errorf("write cache: %w", err)
	}

	// Close checks flush errors on write files.
	if err := f.Close(); err != nil {
		return fmt.Errorf("close cache file: %w", err)
	}
	return nil
}
