package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jnodorp/typist/pkg/typist"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var typeCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "type",
	Short: "Emulate typing standard input (or the contents of a file) to standard output",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			in, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("cannot read %s: %w", args[0], err)
			}

			defer in.Close()
			cmd.SetIn(in)
		}

		a := viper.GetFloat64(flagAccuracy)
		wpm := viper.GetInt(flagWPM)
		t, err := typist.New(wpm, a)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(cmd.InOrStdin())
		for scanner.Scan() {
			t.Type(cmd.OutOrStdout(), scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("cannot read input: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(typeCmd)
}
