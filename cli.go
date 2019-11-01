package awslogin

import (
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youyo/awslogin/browse"
	"github.com/youyo/awsprofile"
)

// Run
func Run(cmd *cobra.Command, args []string) (err error) {
	profile := viper.GetString("profile")
	durationSeconds := viper.GetInt("duration_seconds")
	browser := viper.GetString("browser")
	outputUrl := viper.GetBool("output-url")

	awsSession := NewAwsSession(profile, time.Duration(durationSeconds)*time.Second)

	temporaryCredentials, err := BuildTemporaryCredentials(awsSession)
	if err != nil {
		return err
	}

	requestUrl := BuildSigninTokenRequestURL(temporaryCredentials, strconv.Itoa(durationSeconds))

	signinToken, err := RequestSigninToken(requestUrl)
	if err != nil {
		return err
	}

	signinUrl := BuildSigninURL(signinToken, *awsSession.Config.Region)

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

	configs := awsprofile.NewConfigs()
	file, err := awsprofile.GetConfigsPath()
	if err != nil {
		return err
	}

	if err := configs.Parse(file); err != nil {
		return err
	}

	if selectProfile {
		profiles, err := configs.ProfileNames()
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

	for _, config := range *configs {
		if config.ProfileName == viper.GetString("profile") {

			durationSeconds := config.GetDurationSeconds()

			if durationSeconds != 0 {
				viper.Set("duration_seconds", durationSeconds)
			} else {
				viper.Set("duration_seconds", 3600)
			}

			break
		}
	}

	return nil
}
