package variant_plugin

import (
	"context"

	"github.com/nfwGytautas/appy-cli/utils"
)

func (cfg *Config) Start(ctx context.Context) error {
	utils.Console.DebugLn("Starting Plugin...")
	return nil
}
