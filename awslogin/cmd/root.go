package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youyo/awslogin"
	"github.com/youyo/awssh"
)

var Version string

var rootCmd = &cobra.Command{
	Use:          "awslogin",
	Short:        "Login to the AWS management console.",
	Version:      Version,
	PreRunE:      awssh.PreRun,
	RunE:         awslogin.Run,
	SilenceUsage: true,
}

// Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("select-profile", "S", false, "interactive select profile")
	rootCmd.Flags().BoolP("cache", "c", false, "enable cache a credentials.")
	rootCmd.Flags().StringP("profile", "p", "default", "use a specific profile from your credential file.")
	rootCmd.Flags().StringP("browser", "b", "", "Opens with the specified browse application")
	rootCmd.Flags().BoolP("output-url", "O", false, "output signin url")
	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	viper.BindEnv("profile", "AWS_PROFILE")
}
