package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/YourTechBud/inferix/server"
)

var (
	// Used for flags
	configDriver           string
	configPath             string
	defaultConfigPath      string
	createDefaultWorkspace bool

	rootCmd = &cobra.Command{
		Use:   "inferix",
		Short: "Inferix is a OpenAI compatible backend to build Generative AI applications.",
		Run: func(cmd *cobra.Command, args []string) {
			// Create the server
			router, err := server.New(server.Options{
				ConfigDriver:           server.ConfigDriverType(configDriver),
				ConfigPath:             configPath,
				DefaultConfigPath:      defaultConfigPath,
				CreateDefaultWorkspace: createDefaultWorkspace,
			})
			if err != nil {
				panic(err)
			}

			// Start the server
			if err := router.Start(context.TODO()); err != nil {
				panic(err)
			}
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
	rootCmd.Flags().StringVar(&configDriver, "config-driver", "file", "The configuration driver to use.")
	rootCmd.Flags().StringVar(&configPath, "config-path", "inferix.yaml", "Path to your configuration.")
	rootCmd.Flags().StringVar(&defaultConfigPath, "default-config-path", "", "Path to your default configuration.")
	rootCmd.Flags().BoolVar(&createDefaultWorkspace, "create-default-workspace", true, "Create a default workspace.")
}
