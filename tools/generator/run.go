package main

import (
	"fmt"
	"os"

	"github.com/findy-network/findy-agent-vault/tools/data"

	"github.com/findy-network/findy-agent-vault/server"

	"github.com/findy-network/findy-agent-vault/tools/faker"

	"github.com/spf13/cobra"
)

const graphiQLURL = "http://localhost:8085/"

var fakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "Generate fake data",
	Run: func(cmd *cobra.Command, args []string) {
		s := data.InitState()
		faker.Run(s.Connections().Objects(), s.Events, s.Messages)
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate test auth token",
	Run: func(cmd *cobra.Command, args []string) {
		t, err := server.CreateToken("Emmett")
		if err == nil {
			fmt.Printf("Generated token.\nCopy and paste following to graphiQL (%s) \"Headers\" section:\n", graphiQLURL)
			fmt.Printf("{\"Authorization\": \"Bearer %s\"}\n", t)
		} else {
			fmt.Printf("Error generating token: %s\n", err.Error())
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "Helper tool to generate test assets",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	rootCmd.AddCommand(fakeCmd)
	rootCmd.AddCommand(tokenCmd)

	Execute()
}
