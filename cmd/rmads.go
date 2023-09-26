package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tartale/remove-ads/pkg/config"
	"github.com/tartale/remove-ads/pkg/rmads"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rmads",
	Short: "Remove ads from video files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rmads.RemoveAds(cmd.Context())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(func() {
		config.InitConfig(cfgFile)

		// Only validate config values in running programs
		err := config.Values.Validate()
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(os.Stdout, "config loaded", config.Values)
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "preview",
		Short: "Shows a preview of the cut points, then exits",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rmads.Preview(cmd.Context())
		},
	})

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.rmads.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.Values.SkipFilePath, "skip", "s", "", "Skip file; contains the list of timestamps that indicate location of ads")
	rootCmd.PersistentFlags().StringVarP(&config.Values.InputFilePath, "input", "i", "", "Input video path; must be a path to an existing file")
	rootCmd.PersistentFlags().StringVarP(&config.Values.OutputFilePath, "output", "o", "", "Output path; defaults to stdout")

	// Copy those flags into root command
	rootCmd.PersistentFlags().AddFlagSet(pflag.CommandLine)
	viper.BindPFlags(rootCmd.PersistentFlags())

}
