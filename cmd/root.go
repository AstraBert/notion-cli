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

var maxRetries int
var retryTime int

var readCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Get a Notion Page providing its ID.",
	Long:    "Get the content of a Notion Page by providing its ID.",
	Run: func(cmd *cobra.Command, args []string) {
		pageId := args[0]
		if pageId == "" {
			fmt.Println("\x1b[1;31mYou must provide a positional argument for `page_id`")
			os.Exit(1)
		}
		notionClient, err := internals.NewNotionClientFromDefaults()
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while initializing the Notion client: %s\n", err.Error())
			os.Exit(1)
		}
		app := internals.NewNotion(notionClient)
		content, err := app.Read(pageId, maxRetries, retryTime)
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
var writeMaxRetries int
var writeRetryTime int

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
		pageId, err := app.Write(writeContent, title, parentId, parentLiteral, writeMaxRetries, writeRetryTime)
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while writing the Notion page: %s\n", err.Error())
			os.Exit(2)
		}
		fmt.Println(pageId)
	},
}

var appendContent string
var appendMaxRetries int
var appendRetryTime int

var appendCmd = &cobra.Command{
	Use:     "append",
	Aliases: []string{"a"},
	Short:   "Append markdown content at the end of a Notion page.",
	Long:    "Append markdown content at the end of a Notion page, by providing the page ID and the content. Returns the ID of the modified page.",
	Run: func(cmd *cobra.Command, args []string) {
		pageId := args[0]
		if pageId == "" {
			fmt.Println("\x1b[1;31mYou must provide a positional argument for `page_id`")
			os.Exit(1)
		}
		if appendContent == "" {
			fmt.Println("\x1b[1;31mYou must provide an argument for `--content`/`-c`")
			os.Exit(1)
		}
		notionClient, err := internals.NewNotionClientFromDefaults()
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while initializing the Notion client: %s\n", err.Error())
			os.Exit(1)
		}
		app := internals.NewNotion(notionClient)
		returnedId, err := app.Append(pageId, appendContent, appendMaxRetries, appendRetryTime)
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while appending content to the Notion page: %s\n", err.Error())
			os.Exit(2)
		}
		fmt.Println(returnedId)
		os.Exit(0)
	},
}

var searchSortStrategy string
var searchPageSize int
var searchMaxRetries int
var searchRetryTime int

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Search pages (by title) within your Notion workspace.",
	Long:    "Search pages (by title) within your Notion workspace, providing a query, and, optionally, parameters like page size and sorting direction. Returns an array of page IDs.",
	Run: func(cmd *cobra.Command, args []string) {
		searchQuery := args[0]
		if searchQuery == "" {
			fmt.Println("\x1b[1;31mYou must provide a positional argument for `query`")
			os.Exit(1)
		}
		var sortStrategy internals.SortStrategyLiteral
		switch searchSortStrategy {
		case string(internals.AscendingSortStrategy):
			sortStrategy = internals.AscendingSortStrategy
		case string(internals.DescendingSortStrategy):
			sortStrategy = internals.DescendingSortStrategy
		default:
			fmt.Printf("\x1b[1;31mInvalid argument for `--sort`/`-s`: %s. Allowed arguments are: 'ascending', 'descending'\n", parentType)
			os.Exit(1)
		}
		notionClient, err := internals.NewNotionClientFromDefaults()
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while initializing the Notion client: %s\n", err.Error())
			os.Exit(1)
		}
		app := internals.NewNotion(notionClient)

		returnedIds, err := app.Search(searchQuery, "", sortStrategy, searchPageSize, searchMaxRetries, searchRetryTime)
		if err != nil {
			fmt.Printf("\x1b[1;31mAn error occurred while appending content to the Notion page: %s\n", err.Error())
			os.Exit(2)
		}
		for _, returnedId := range returnedIds {
			fmt.Println(returnedId)
		}
		os.Exit(0)
	},
}

func init() {
	writeCmd.Flags().StringVarP(&parentType, "parent-type", "p", "page", "Type of parent ('database' or 'page'). Defaults to 'page'.")
	writeCmd.Flags().StringVarP(&parentId, "parent-id", "i", "", "ID of the parent element. Required.")
	writeCmd.Flags().StringVarP(&writeContent, "content", "c", "", "Markdown content to write. Required.")
	writeCmd.Flags().StringVarP(&title, "title", "t", "", "Title for the page. Defaults to an empty string.")
	writeCmd.Flags().IntVarP(&writeMaxRetries, "max-retries", "m", internals.MaxRetries, "Maximum number of retries for failed API calls. Defaults to 3.")
	writeCmd.Flags().IntVarP(&writeRetryTime, "retry-interval", "r", internals.DefaultRetryTime, "Retry interval (in seconds) for failed API calls. Defaults to 1 second.")
	readCmd.Flags().IntVarP(&maxRetries, "max-retries", "m", internals.MaxRetries, "Maximum number of retries for failed API calls. Defaults to 3.")
	readCmd.Flags().IntVarP(&retryTime, "retry-interval", "r", internals.DefaultRetryTime, "Retry interval (in seconds) for failed API calls. Defaults to 1 second.")
	appendCmd.Flags().StringVarP(&appendContent, "content", "c", "", "Markdown content to append. Required.")
	appendCmd.Flags().IntVarP(&appendMaxRetries, "max-retries", "m", internals.MaxRetries, "Maximum number of retries for failed API calls. Defaults to 3.")
	appendCmd.Flags().IntVarP(&appendRetryTime, "retry-interval", "r", internals.DefaultRetryTime, "Retry interval (in seconds) for failed API calls. Defaults to 1 second.")
	searchCmd.Flags().StringVarP(&searchSortStrategy, "sort", "s", "descending", "Order to follow when sorting by last edited. Allowed values: 'ascending', 'descending'. Defaults to 'descending'.")
	searchCmd.Flags().IntVarP(&searchPageSize, "page-size", "p", -1, "Page size for paginated API responses. Defaults to -1 (unspecified page size)")
	searchCmd.Flags().IntVarP(&searchMaxRetries, "max-retries", "m", internals.MaxRetries, "Maximum number of retries for failed API calls. Defaults to 3.")
	searchCmd.Flags().IntVarP(&searchRetryTime, "retry-interval", "r", internals.DefaultRetryTime, "Retry interval (in seconds) for failed API calls. Defaults to 1 second.")

	_ = writeCmd.MarkFlagRequired("parent-id")
	_ = writeCmd.MarkFlagRequired("content")
	_ = appendCmd.MarkFlagRequired("content")

	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(appendCmd)
	rootCmd.AddCommand(searchCmd)
}
