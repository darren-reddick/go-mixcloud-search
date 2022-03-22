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

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Search listen history for a user",
	Long:  `Search listen history for a user.`,
	Run: func(cmd *cobra.Command, args []string) {
		include, _ := cmd.Flags().GetStringSlice("include")
		exclude, _ := cmd.Flags().GetStringSlice("exclude")

		filter, err := mixcloud.NewFilter(
			include,
			exclude,
		)

		if err != nil {
			panic(err)
		}

		user, _ := cmd.Flags().GetString("user")

		mc, err := mixcloud.NewHistorySearch(user, filter, &http.Client{}, mixcloud.NewStore(0))

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
	rootCmd.AddCommand(historyCmd)

	historyCmd.Flags().StringP("user", "u", "", "User name to search")
	historyCmd.MarkFlagRequired("user")

	historyCmd.Flags().StringSliceP("include", "i", []string{}, "Filter to include entry")
	historyCmd.Flags().StringSliceP("exclude", "e", []string{}, "Filter to exclude entry")
}
