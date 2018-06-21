package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/youyo/awslogin"
)

var (
	Name            string
	Version         string
	CommitHash      string
	BuildTime       string
	GoVersion       string
	app             string
	profile         string
	readFromEnv     bool
	durationSeconds int
)

const (
	Peco          string = "peco"
	PecoGithubUrl string = "https://github.com/peco/peco"
)

func init() {}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "awslogin",
		Short: "Login to the AWS management console.",
		Long:  `Using AssumeRole, accept IAMRole and login to the AWS management console.`,
		Run: func(cmd *cobra.Command, args []string) {
			c := awslogin.NewConfig()
			if err := c.SetData(); err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}
			if err := c.SelectProfile(profile, readFromEnv); err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}
			if err := c.FetchArn(); err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}
			if durationSeconds != 3600 {
				c.SetDurationSeconds(durationSeconds)
			}

			al := awslogin.NewAwslogin(c)
			al.BuildAssumeRoleProvider()
			if err := al.GetCredentials(); err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}
			if err := al.GetFederatedSession(); err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}
			al.BuildSigninTokenRequestURL()
			if err := al.RequestSigninToken(); err != nil {
				cmd.SetOutput(os.Stderr)
				cmd.Println(err)
				os.Exit(1)
			}
			al.BuildSigninURL()
			signinUrl := al.GetSigninUrl()

			// Open browser
			switch app {
			case "":
				if err := awslogin.Browsing(signinUrl); err != nil {
					cmd.SetOutput(os.Stderr)
					cmd.Println(err)
					os.Exit(1)
				}
			default:
				if err := awslogin.BrowsingSpecificApp(signinUrl, app); err != nil {
					cmd.SetOutput(os.Stderr)
					cmd.Println(err)
					os.Exit(1)
				}
			}
		},
	}
	cobra.OnInitialize(initConfig)
	cmd.Flags().StringVarP(&app, "app", "a", "", "Opens with the specified application.")
	cmd.Flags().StringVarP(&profile, "profile", "p", "", "Use a specific profile.")
	cmd.Flags().BoolVarP(&readFromEnv, "read-from-env", "e", false, "Use a specific profile read from the environment. [$AWS_PROFILE]")
	cmd.Flags().IntVarP(&durationSeconds, "duration-seconds", "d", 3600, "Request a session duration seconds. 900 - 43200")
	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdVersion())
	return cmd
}
