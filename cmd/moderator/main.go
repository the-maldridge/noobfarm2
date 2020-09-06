package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/the-maldridge/noobfarm2/internal/moderator"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&moderator.ListQuotesCmd{}, "Moderation")
	subcommands.Register(&moderator.ShowQuotesCmd{}, "Moderation")
	subcommands.Register(&moderator.ApprovalCmd{}, "Moderation")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
