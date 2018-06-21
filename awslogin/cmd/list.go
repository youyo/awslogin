package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youyo/awslogin"
)

func init() {}

func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List profiles",
		Run: func(cmd *cobra.Command, args []string) {
			c := awslogin.NewConfig()
			if err := c.SetData(); err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}
			aa := c.AvailableArn()
			list := strings.Join(aa[:], "\n")
			cmd.Println(list)
		},
	}
	cobra.OnInitialize(initConfig)
	return cmd
}
