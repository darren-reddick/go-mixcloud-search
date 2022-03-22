/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/darren-reddick/go-mixcloud-search/mixcloud"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for mixes on mixcloud by term",
	Long:  `Search for mixes on mixcloud by term.`,
	Run: func(cmd *cobra.Command, args []string) {

		include, _ := cmd.Flags().GetStringSlice("include")
		exclude, _ := cmd.Flags().GetStringSlice("exclude")
		limit, _ := cmd.Flags().GetInt("limit")

		filter, err := mixcloud.NewFilter(
			include,
			exclude,
		)

		if err != nil {
			panic(err)
		}

		term, _ := cmd.Flags().GetString("term")

		mc, err := mixcloud.NewMixSearch(term, filter, &http.Client{}, mixcloud.NewStore(limit))

		if err != nil {
			fmt.Println(err)
			return
		}

		err = mc.GetAllAsync()

		if err != nil {
			fmt.Println(err)
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

}
