package cli

import (
	"fmt"
	"path/filepath"

	"github.com/jackh54/mcinit/internal/server"
	"github.com/jackh54/mcinit/internal/utils"
)

// defaultConfigName is the config file looked up when --config is not given.
const defaultConfigName = "mcinit.json"

// resolveConfig returns the absolute config path and server folder for the current command.
// With --config, the server folder is the folder that holds the config file and the file must
// exist; without it, both come from the current working directory.
func resolveConfig() (cfgPath, serverDir string, err error) {
	if cfgFile == "" {
		dir, err := utils.AbsolutePath(".")
		if err != nil {
			return "", "", fmt.Errorf("failed to resolve current directory: %w", err)
		}
		return filepath.Join(dir, defaultConfigName), dir, nil
	}

	abs, err := utils.AbsolutePath(cfgFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve --config path: %w", err)
	}
	if !utils.PathExists(abs) {
		return "", "", fmt.Errorf("config file not found: %s", abs)
	}
	return abs, filepath.Dir(abs), nil
}

// newManager builds a server manager from --config (or ./mcinit.json).
func newManager() (*server.Manager, error) {
	cfgPath, _, err := resolveConfig()
	if err != nil {
		return nil, err
	}
	mgr, err := server.NewManagerForConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create server manager: %w", err)
	}
	return mgr, nil
}
