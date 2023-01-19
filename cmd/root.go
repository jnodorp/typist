package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagAccuracy = "accuracy"
	flagWPM      = "wpm"
)

const (
	defaultAccuracy = 0.97
	defaultWPM      = 75
)

var rootCmd = &cobra.Command{
	Use: "typist",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	viper.SetEnvPrefix(rootCmd.Name())
	cobra.OnInitialize(viper.AutomaticEnv)

	rootCmd.PersistentFlags().Float64P(flagAccuracy, "a", defaultAccuracy, "accuracy of the simulated keystrokes")
	rootCmd.PersistentFlags().IntP(flagWPM, "w", defaultWPM, "average words per Minute")

	viper.BindPFlags(rootCmd.PersistentFlags())
}
