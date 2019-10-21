package awslogin

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youyo/awslogin/browse"
	"github.com/youyo/awsprofile"
)

const (
	CachePath string = "~/.config/awslogin/cache"
)

// Run
func Run(cmd *cobra.Command, args []string) (err error) {
	profile := viper.GetString("profile")
	browser := viper.GetString("browser")
	cache := viper.GetBool("cache")
	outputUrl := viper.GetBool("output-url")

	awsSession, err := NewAwsSession(profile, cache, CachePath)
	if err != nil {
		return err
	}

	temporaryCredentials, err := BuildTemporaryCredentials(awsSession)
	if err != nil {
		return err
	}

	requestUrl := BuildSigninTokenRequestURL(temporaryCredentials)

	signinToken, err := RequestSigninToken(requestUrl)
	if err != nil {
		return err
	}

	signinUrl := BuildSigninURL(signinToken)

	if outputUrl {
		cmd.Println(signinUrl)
		return nil
	}

	if browser != "" {
		browse.StartWith(signinUrl, browser)
	} else {
		browse.Start(signinUrl)
	}

	return nil
}

func PreRun(cmd *cobra.Command, args []string) (err error) {
	selectProfile := viper.GetBool("select-profile")
	if selectProfile {
		awsProfile := awsprofile.New()

		if err := awsProfile.Parse(); err != nil {
			return err
		}

		profiles, err := awsProfile.ProfileNames()
		if err != nil {
			return err
		}

		prompt := promptui.Select{
			Label: "Profiles",
			Templates: &promptui.SelectTemplates{
				Label:    `{{ . | green }}`,
				Active:   `{{ ">" | blue }} {{ . | red }}`,
				Inactive: `{{ . | cyan }}`,
				Selected: `{{ . | yellow }}`,
			},
			Items: profiles,
			Size:  25,
			Searcher: func(input string, index int) bool {
				item := profiles[index]
				profileName := strings.Replace(strings.ToLower(item), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)
				if strings.Contains(profileName, input) {
					return true
				}
				return false
			},
			StartInSearchMode: true,
		}

		index, _, err := prompt.Run()
		if err != nil {
			return err
		}

		viper.Set("profile", profiles[index])
	}

	return nil
}
