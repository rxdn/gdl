package command

type Command struct {
	Name        string
	Aliases     []string
	Handler     func(CommandContext)
	SubCommands []Command
}

func NewCommand(name string, aliases []string, handler func(CommandContext)) Command {
	return Command{
		Name:    name,
		Aliases: aliases,
		Handler: handler,
	}
}

func (c *Command) RegisterSubCommand(sub Command) {
	c.SubCommands = append(c.SubCommands, sub)
}
