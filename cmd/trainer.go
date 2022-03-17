package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/kaushikc92/cagliostro/pkg/trainer"
)

var trainerCmd = &cobra.Command{
	Use:   "trainer",
	Short: "Interactive Trainer",
	Long: `Interactive trainer to practice moves`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please specify a board position in FEN")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := trainer.Interactive(args[0])
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(trainerCmd)
}
