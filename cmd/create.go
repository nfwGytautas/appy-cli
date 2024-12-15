package cmd

import (
	"github.com/nfwGytautas/appy-cli/logic"
	"github.com/spf13/cobra"
)

// createCmd represents the scaffold command
var createCmd = &cobra.Command{
	Use:       "create [endpoint]",
	Short:     "",
	Long:      ``,
	ValidArgs: []string{"endpoint"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		logic.Create(args)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scaffoldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scaffoldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
