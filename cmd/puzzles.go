package cmd

import (
	"errors"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/kaushikc92/cagliostro/pkg/puzzles"
)

var puzzlesCmd = &cobra.Command{
	Use:   "puzzles",
	Short: "Puzzles",
	Long: `Interactive puzzle trainer to practice moves`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Please specify min and max rating")
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("Please enter a min and max rating")
		}
		if _, err := strconv.Atoi(args[1]); err != nil {
			return errors.New("Please enter a min and max rating")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		minRating, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		maxRating, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		err = puzzles.Interactive(minRating, maxRating)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(puzzlesCmd)
}

