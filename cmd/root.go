package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/YourTechBud/inferix/server"
)

var (
	// Used for flags
	configFilePath string

	rootCmd = &cobra.Command{
		Use:   "inferix",
		Short: "Inferix is a OpenAI compatible backend to build Generative AI applications.",
		Run: func(cmd *cobra.Command, args []string) {

			// Create the server
			router, err := server.New(configFilePath)
			if err != nil {
				panic(err)
			}

			fmt.Println("Starting server on port 4386")
			http.ListenAndServe(":4386", router)
		},
	}
)

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&configFilePath, "config", "inferix.yaml", "Path to your config file.")
}
