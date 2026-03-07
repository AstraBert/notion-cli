package cmd

import (
	"fmt"
	"os"

	"github.com/AstraBert/notion-cli/internals"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "notion-cli",
	Short: "notion-cli is a CLI tool to read and write Notion pages.",
	Long:  "notion-cli is a CLI tool that reads and writes Notion pages connecting to the Notion API.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing bsgit '%s'\n", err)
		os.Exit(1)
	}
}

var readCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Get a Notion Page providing its ID.",
	Long:    "Get the content of a Notion Page by providing its ID.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pageId := args[0]
		notionClient, err := internals.NewNotionClientFromDefaults()
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while initializing the Notion client: %s\n", err.Error())
			os.Exit(1)
		}
		app := internals.NewNotion(notionClient)
		content, err := app.Read(pageId)
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while reading the Notion page: %s\n", err.Error())
			os.Exit(2)
		}
		fmt.Println(content)
		os.Exit(0)
	},
}

var parentType string
var writeContent string
var parentId string
var title string

var writeCmd = &cobra.Command{
	Use:     "write",
	Aliases: []string{"w"},
	Short:   "Create a Notion page and return the page ID.",
	Long:    "Create a Notion page by providing its content, title, and parent element, and return the ID of the newly created page.",
	Run: func(cmd *cobra.Command, args []string) {
		if parentId == "" {
			fmt.Println("\x1b[1;31mYou must provide an argument for `--parent-id`/`-i`")
			os.Exit(1)
		} else if writeContent == "" {
			fmt.Println("\x1b[1;31mYou must provide an argument for `--content`/`-c`")
			os.Exit(1)
		}
		var parentLiteral internals.ParentLiteral
		switch parentType {
		case string(internals.DatabaseParentLiteral):
			parentLiteral = internals.DatabaseParentLiteral
		case string(internals.PageParentLiteral):
			parentLiteral = internals.PageParentLiteral
		default:
			fmt.Printf("\x1b[1;31mInvalid argument for `--parent-type`/`-t`: %s. Allowed arguments are: 'page', 'database'\n", parentType)
			os.Exit(1)
		}
		notionClient, err := internals.NewNotionClientFromDefaults()
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while initializing the Notion client: %s\n", err.Error())
			os.Exit(1)
		}
		app := internals.NewNotion(notionClient)
		pageId, err := app.Write(writeContent, title, parentId, parentLiteral)
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while writing the Notion page: %s\n", err.Error())
			os.Exit(2)
		}
		fmt.Println(pageId)
	},
}

func init() {
	writeCmd.Flags().StringVarP(&parentType, "parent-type", "p", "page", "Type of parent ('database' or 'page'). Defaults to 'page'.")
	writeCmd.Flags().StringVarP(&parentId, "parent-id", "i", "", "ID of the parent element. Required.")
	writeCmd.Flags().StringVarP(&writeContent, "content", "c", "", "Markdown content to write. Required.")
	writeCmd.Flags().StringVarP(&title, "title", "t", "", "Title for the page. Defaults to an empty string.")

	_ = writeCmd.MarkFlagRequired("parent-id")
	_ = writeCmd.MarkFlagRequired("content")

	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(writeCmd)
}
