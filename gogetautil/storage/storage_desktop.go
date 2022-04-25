//go:build (darwin || freebsd || linux || windows) && !js && !android && !ios
// +build darwin freebsd linux windows
// +build !js
// +build !android
// +build !ios

package storage

import (
	"os"
	"path/filepath"
)

func DefaultStoragePath(appName string) string {
	return "game_data_" + appName + ".dat"
}

func Save(path string, data []byte) error {
	return os.WriteFile(filepath.FromSlash(path), data, 0644)
}

func Load(path string) ([]byte, error) {
	return os.ReadFile(filepath.FromSlash(path))
}
