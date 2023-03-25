package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
	"os"
	"wpull/internal/printline"
	"wpull/internal/snok"
	"wpull/internal/stringtoint"
)

var (
	noColor            bool
	debug              bool
	threads            int
	overwrite          bool
	versionFlag        bool
	chunkSizeStr       string
	username, password string
	version            = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wsnok",
	Short: "Download a file from a URL with multiple threads",
	Long: `Downloads a file
with multiple threads. The amount of threads can be simply set with the -t flag.`,
	Run: wPull,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&noColor, "no-color", "n", false, "Disable color output")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show the current version")
	rootCmd.Flags().IntVarP(&threads, "threads", "t", 10, "How many threads to use")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite if file already exists")
	rootCmd.Flags().StringVarP(&chunkSizeStr, "chunk-size", "c", "1M", "Set the default chunk size")
	rootCmd.Flags().StringVarP(&username, "http-user", "u", "", "Set Username for Basic Auth")
	rootCmd.Flags().StringVarP(&password, "http-password", "p", "", "Set Password for Basic Auth")
}

func wPull(cmd *cobra.Command, args []string) {

	if versionFlag {
		fmt.Printf("Wsnok %s\n", version)
		os.Exit(0)
	}

	if len(args) == 0 {
		showHelpText("%s: missing URL\n")
	}

	chunkSize, err := stringtoint.Parse(chunkSizeStr)
	if err != nil {
		showHelpText("%s: Unknown chunk-size string. please use 1M or 512B etc...\n")
	}

	snok.Debug = debug
	snok.Threads = threads
	snok.ChunkSize = chunkSize
	snok.OverWrite = overwrite
	snok.Username = username
	snok.Password = password
	for _, url := range args {
		err := snok.Snok(url)
		if err != nil {
			printline.Printf(false, "Could not download file: %s (%s)\n", url, err.Error())
		}
	}

	os.Exit(0)

}

func showColor() bool {
	return !noColor
}

func showHelpText(error string) {
	args := os.Args
	filename := args[0]
	if showColor() {
		filename = fmt.Sprintf("%s", aurora.Bold(filename))
	}
	printline.Printf(false, error, filename)
	printline.Printf(false, "Usage: %s [OPTION]... [URL]...\n\n", filename)
	printline.Printf(false, "Try '%s --help' for more options.\n", filename)
	os.Exit(1)
}
