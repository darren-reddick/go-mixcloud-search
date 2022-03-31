/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/darren-reddick/go-mixcloud-search/mixcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger(flags *pflag.FlagSet) error {
	var err error
	// initialize the logger config and update level dependent
	// on the debug flag being set
	cfg := zap.NewProductionConfig()

	debug, err := flags.GetBool("debug")

	if err != nil {
		return err
	}

	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, err = cfg.Build()

	return err
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for mixes on mixcloud by term",
	Long:  `Search for mixes on mixcloud by term.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := initLogger(cmd.Flags())
		if err != nil {
			fmt.Printf("Error initializing logger: %s", err)
			return
		}

		include, _ := cmd.Flags().GetStringSlice("include")
		exclude, _ := cmd.Flags().GetStringSlice("exclude")
		limit, _ := cmd.Flags().GetInt("limit")

		filter, err := mixcloud.NewFilter(
			include,
			exclude,
		)

		if err != nil {
			logger.Panic(err.Error())
		}

		term, _ := cmd.Flags().GetString("term")

		mc, err := mixcloud.NewMixSearch(term, filter, &http.Client{}, mixcloud.NewStore(limit), logger)

		if err != nil {
			logger.Error(err.Error())
			return
		}

		logger.Debug(fmt.Sprintf("Running search for term %s", term))
		err = mc.GetAllParallel()

		if err != nil {
			logger.Error(err.Error())
			return
		}

		mc.WriteJsonToFile()

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	searchCmd.Flags().StringP("term", "t", "", "Search term")
	searchCmd.MarkFlagRequired("term")

	searchCmd.Flags().StringSliceP("include", "i", []string{}, "Filter to include entry")
	searchCmd.Flags().StringSliceP("exclude", "e", []string{}, "Filter to exclude entry")

	searchCmd.Flags().IntP("limit", "l", 0, "Limit number of results")

	searchCmd.Flags().BoolP("debug", "d", false, "Enable debug")

}
