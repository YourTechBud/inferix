package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/YourTechBud/inferix/server"
	"github.com/YourTechBud/inferix/utils/hash"
)

var (
	rootCmd = &cobra.Command{
		Use:   "inferix",
		Short: "Inferix is an OpenAI compatible backend to build Generative AI applications.",
		Run: func(cmd *cobra.Command, args []string) {
			// Extract the configuration
			flags := RootFlags{}
			cobra.CheckErr(viper.Unmarshal(&flags))

			// Initialise the hasher
			cobra.CheckErr(hash.InitialiseHasher(flags.Hash))

			// Create the server
			router, err := server.New(flags.Options)
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

type RootFlags struct {
	server.Options `mapstructure:",squash"`

	Hash hash.Options `mapstructure:"hash"`
}

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
	rootCmd.Flags().String("config-driver", "libsql", "The configuration driver to use.")
	rootCmd.Flags().String("config-path", "./inferix", "Path to your configuration.")
	rootCmd.Flags().String("default-config-path", "", "Path to your default configuration.")
	rootCmd.Flags().Bool("create-default-workspace", true, "Create a default workspace.")

	rootCmd.Flags().String("hash.algo", "bcrypt", "The hasher to use.")
	rootCmd.Flags().Int("hash.bcrypt.cost", 12, "The cost for the bcrypt hasher.")

	rootCmd.Flags().Bool("auth.enabled", false, "Enable basic authentication.")
	rootCmd.Flags().String("auth.user", "admin", "The username to use for authentication.")
	rootCmd.Flags().String("auth.pass", "1234", "The password to use for authentication.")

	// Bind the flags with viper
	viper.BindPFlag("config-driver", rootCmd.Flags().Lookup("config-driver"))
	viper.BindPFlag("config-path", rootCmd.Flags().Lookup("config-path"))
	viper.BindPFlag("default-config-path", rootCmd.Flags().Lookup("default-config-path"))
	viper.BindPFlag("create-default-workspace", rootCmd.Flags().Lookup("create-default-workspace"))

	viper.BindPFlag("hash.algo", rootCmd.Flags().Lookup("hash.algo"))
	viper.BindPFlag("hash.bcrypt.cost", rootCmd.Flags().Lookup("hash.bcrypt.cost"))

	viper.BindPFlag("auth.enabled", rootCmd.Flags().Lookup("auth.enabled"))
	viper.BindPFlag("auth.user", rootCmd.Flags().Lookup("auth.user"))
	viper.BindPFlag("auth.pass", rootCmd.Flags().Lookup("auth.pass"))
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
