package main

import (
	"log"

	"github.com/glacius-labs/schedulerd/internal/cli"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "schedulerd",
		Short: "Stateless workload scheduler powered by rule-based filtering and scoring",
	}

	root.AddCommand(cli.EvalCmd())

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
