package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

type NativeBrowser struct {
}

func (b *NativeBrowser) OpenURL(url string) error {
	return nil
}

func (b *NativeBrowser) OpenFile(file string) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	return b.OpenURL("file://" + path)
}

func (b *NativeBrowser) run(provider string, args ...string) error {
	cmd := exec.Command(provider, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
