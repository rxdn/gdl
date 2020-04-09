package command

import "github.com/rxdn/gdl/gateway"

type Command struct {
	Name    string
	Aliases []string
	Handler func(CommandContext)
}

func NewCommand(name string, aliases []string, handler func(CommandContext)) Command {
	return Command{
		Name:    name,
		Aliases: aliases,
		Handler: handler,
	}
}

func (c *Command) Register(sm *gateway.ShardManager) {

}
