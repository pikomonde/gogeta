//go:build !(darwin || freebsd || linux || windows) && js && !android && !ios
// +build !darwin,!freebsd,!linux,!windows
// +build js
// +build !android
// +build !ios

package storage

import (
	"syscall/js"
)

func DefaultStoragePath(appName string) string {
	return "game_data::" + appName
}

func Save(path string, data []byte) error {
	js.Global().Get("localStorage").Call("setItem", path, string(data))
	return nil
}

func Load(path string) ([]byte, error) {
	data := js.Global().Get("localStorage").Call("getItem", "path")
	return []byte(data.String()), nil
}
