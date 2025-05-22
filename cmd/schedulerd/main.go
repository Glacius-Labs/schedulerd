package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/glacius-labs/schedulerd/cmd/schedulerd/cli"
	"github.com/spf13/cobra"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	root := &cobra.Command{
		Use:   "schedulerd",
		Short: "Stateless workload scheduler powered by rule-based filtering and scoring",
	}

	root.AddCommand(cli.EvalCmd(logger))

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
