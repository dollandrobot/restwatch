package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"runtime"
)

type UserOptions struct {
	Port              int    `json:"port"`
	MaxMessagesToKeep int    `json:"maxMessagesToKeep"`
	DefaultEndpoint   string `json:"defaultEndpoint"`
	JumpToLatest      bool   `json:"jumpToLatest"`
}

func loadUserOptions() (UserOptions, error) {
	opts := UserOptions{
		Port:              2999,
		MaxMessagesToKeep: 100,
		DefaultEndpoint:   "/messages",
		JumpToLatest:      true,
	}

	configDir, err := localConfigDir()
	if err != nil {
		return opts, err
	}

	if b, err := os.ReadFile(filepath.Join(configDir, "settings.json")); err == nil {
		if err := json.Unmarshal(b, &opts); err != nil {
			return opts, err
		}
	} else if !os.IsNotExist(err) {
		return opts, err
	}

	return opts, nil
}

func saveUserOptions(opts UserOptions) error {
	b, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	configDir, err := localConfigDir()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, "settings.json"), b, os.ModePerm)
}

func localConfigDir() (string, error) {
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}
	case "darwin":
		dir = os.Getenv("HOME") + "/Library/Application Support"
	default: // Unix
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			dir = os.Getenv("XDG_CONFIG_HOME")
		} else {
			dir = filepath.Join(os.Getenv("HOME"), ".config")
		}
	}

	localConfig := filepath.Join(dir, filepath.Join("RestWatch"))
	os.MkdirAll(localConfig, os.ModePerm)
	return localConfig, nil
}
