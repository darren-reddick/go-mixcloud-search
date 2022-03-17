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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		term, _ := cmd.Flags().GetString("term")

		mc, err := mixcloud.NewSearch(term, filter, &http.Client{}, mixcloud.NewStore())

		if err != nil {
			fmt.Println(err)
			return
		}

		//err := mc.GetAllSync()
		err = mc.GetAllAsync()

		if err != nil {
			fmt.Println(err)
			return
		}

		mc.WriteJsonToFile()

		//fmt.Printf("%+v\n", rez)

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

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
