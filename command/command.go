package command

type Command struct {
	Name        string
	Aliases     []string
	Async       bool
	Handler     func(CommandContext)
	SubCommands []Command
}

func NewCommand(name string, aliases []string, async bool, handler func(CommandContext)) Command {
	return Command{
		Name:    name,
		Aliases: aliases,
		Async:   async,
		Handler: handler,
	}
}

func (c *Command) RegisterSubCommand(sub Command) {
	c.SubCommands = append(c.SubCommands, sub)
}
