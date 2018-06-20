package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/the-maldridge/NoobFarm2/internal/moderator"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&moderator.ListQuotesCmd{}, "Moderation")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
