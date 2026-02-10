package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func fetchAndDisplay(topic string) error {
	req, err := http.NewRequest("GET", "https://cheat.sh/"+topic+"?T", nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "curl/8.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no results for %q", topic)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cheat.sh returned %d", resp.StatusCode)
	}

	// ?T strips ANSI codes so bat handles highlighting.
	bat := exec.Command("bat", "--plain", "--language=sh", "--paging=always")
	bat.Stdin = resp.Body
	bat.Stdout = os.Stdout
	bat.Stderr = os.Stderr

	if err := bat.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("bat exited %d", exitErr.ExitCode())
		}
		return fmt.Errorf("bat: %w", err)
	}
	return nil
}

func requireCmd(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s not found. Install: brew install %s", name, name)
	}
	return nil
}
