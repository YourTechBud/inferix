package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/YourTechBud/inferix/server"
)

var (
	rootCmd = &cobra.Command{
		Use:   "inferix",
		Short: "Inferix is an OpenAI compatible backend to build Generative AI applications.",
		Run: func(cmd *cobra.Command, args []string) {
			// Extract the configuration
			serverOptions := server.Options{}
			cobra.CheckErr(viper.Unmarshal(&serverOptions))

			// Create the server
			router, err := server.New(serverOptions)
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
	cobra.OnInitialize(initConfig)

	// Setup the flags
	rootCmd.Flags().String("config-driver", "file", "The configuration driver to use.")
	rootCmd.Flags().String("config-path", "inferix.yaml", "Path to your configuration.")
	rootCmd.Flags().String("default-config-path", "", "Path to your default configuration.")
	rootCmd.Flags().Bool("create-default-workspace", true, "Create a default workspace.")

	// Bind the flags with viper
	viper.BindPFlag("config-driver", rootCmd.Flags().Lookup("config-driver"))
	viper.BindPFlag("config-path", rootCmd.Flags().Lookup("config-path"))
	viper.BindPFlag("default-config-path", rootCmd.Flags().Lookup("default-config-path"))
	viper.BindPFlag("create-default-workspace", rootCmd.Flags().Lookup("create-default-workspace"))
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
