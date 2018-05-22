package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youyo/awslogin"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := execList(os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

func execList(w io.Writer) (err error) {
	cfg, err := awslogin.NewConfig()
	if err != nil {
		return
	}
	aa := cfg.AvailableArn()
	list := strings.Join(aa[:], "\n")
	fmt.Fprintln(w, list)
	return
}

func init() {
	RootCmd.AddCommand(listCmd)
}
