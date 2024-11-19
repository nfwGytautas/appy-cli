/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/nfwGytautas/appy-cli/logic"
	"github.com/spf13/cobra"
)

// scaffoldCmd represents the scaffold command
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		logic.Scaffold()
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scaffoldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scaffoldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
