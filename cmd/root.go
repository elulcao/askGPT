package cmd

import (
	"fmt"
	"os"

	"github.com/elulcao/askGPT/pkg/ask"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string
var vpr = viper.New()

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "askGPT",
	Short: "Ask to the GPT model",
	Long: `Ask to the GPT model and get the answer to your question.
Running this command will start a prompt where you can type your question.
When you are done, press Enter to get the answer.`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}
		endpoint, err := cmd.Flags().GetString("endpoint")
		if err != nil {
			return err
		}

		if token == "" || endpoint == "" {
			return fmt.Errorf("token and endpoint are required")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		endpoint, _ := cmd.Flags().GetString("endpoint")

		err := ask.GPT(token, endpoint)
		if err != nil && err.Error() != "user cancelled the operation" {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.askGPT.yaml)")
	rootCmd.PersistentFlags().StringP("token", "t", "", "OPenAI API token")
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "OpenAI endpoint")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		vpr.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".askGPT" (without extension).
		vpr.AddConfigPath(".")
		vpr.AddConfigPath(home)
		vpr.SetConfigType("yaml")
		vpr.SetConfigName(".askGPT")
	}

	vpr.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := vpr.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", vpr.ConfigFileUsed())
		viperToCobraFlags(rootCmd, vpr)
	}
}

// viperToCobraFlags copies the values from viper to cobra flags.
func viperToCobraFlags(cmd *cobra.Command, vpr *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if vpr.IsSet(f.Name) {
			_ = cmd.Flags().Set(f.Name, vpr.GetString(f.Name))
		}
	})
}
