package scaffolds

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/shared"
)

type ScaffoldFn func(cfg *config.AppyConfig) error

var scaffoldTypes = map[string]ScaffoldFn{
	shared.ScaffoldHMD: scaffoldHMD,
}

func Scaffold(module string, scaffoldType string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg := config.AppyConfig{
		Module:  module,
		Type:    scaffoldType,
		Project: filepath.Base(dir),
		Version: shared.Version,
	}

	scaffold, exists := scaffoldTypes[scaffoldType]
	if !exists {
		return errors.New("invalid scaffold type " + scaffoldType)
	}

	err = scaffold(&cfg)
	if err != nil {
		return err
	}

	err = cfg.Reconfigure()
	if err != nil {
		return err
	}

	err = cfg.Save()
	if err != nil {
		return err
	}

	return nil
}
