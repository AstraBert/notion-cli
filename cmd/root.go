package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AstraBert/notion-cli/internals"
	"github.com/spf13/cobra"
)

var showHelp bool
var pageId string

var rootCmd = &cobra.Command{
	Use:   "notion-cli",
	Short: "notion-cli is a CLI tool to read and write Notion pages.",
	Long:  "notion-cli is a CLI tool that reads and writes Notion pages connecting to the Notion API.",
	Run: func(cmd *cobra.Command, args []string) {
		if showHelp {
			_ = cmd.Help()
			return
		}
		client, err := internals.NewNotionClientFromDefaults()
		if err != nil {
			log.Printf("An error occurred: %s\n", err.Error())
			os.Exit(1)
		}
		app := internals.NewNotion(client)
		if pageId != "" {
			result, err := app.Read(pageId)
			if err != nil {
				log.Printf("An error occurred: %s\n", err.Error())
				os.Exit(2)
			}
			fmt.Println(result)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing bsgit '%s'\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showHelp, "help", "h", false, "Show the help message and exit.")
	rootCmd.Flags().StringVarP(&pageId, "page", "p", "", "ID of the page to read")
}
