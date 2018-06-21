package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/youyo/awslogin"
)

func init() {}

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			ver := awslogin.NewVersions().
				WithVersion(awslogin.Version).
				WithCommitHash(CommitHash).
				WithBuildTime(BuildTime).
				WithGoVersion(GoVersion).
				WithGithubOwner("youyo").
				WithGithubRepository(Name)
			ver.SetGithubTag()
			res, err := ver.FetchVersionData(ver.GithubTag, ver.Version)
			if err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
			}
			ver.SetVersionCheckData(res)
			message := ver.OutputVersionMessage()
			cmd.Print(message)
		},
	}
	cobra.OnInitialize(initConfig)
	return cmd
}
