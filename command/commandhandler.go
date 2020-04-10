package command

import (
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/gateway/payloads/events"
	"strings"
)

type CommandHandler struct {
	shardManager *gateway.ShardManager
	prefixes     []string
	commands     []Command
}

func NewCommandHandler(sm *gateway.ShardManager, prefixes ...string) *CommandHandler {
	ch := &CommandHandler{
		shardManager: sm,
		prefixes:     prefixes,
	}

	sm.RegisterListeners(ch.commandListener)

	return ch
}

func (h *CommandHandler) RegisterCommand(cmd Command) {
	h.commands = append(h.commands, cmd)
}

func (h *CommandHandler) commandListener(s *gateway.Shard, e *events.MessageCreate) {
	var isCommand bool
	var usedPrefix string

	for _, prefix := range h.prefixes {
		if strings.HasPrefix(e.Content, prefix) {
			isCommand = true
			usedPrefix = prefix
			break
		}
	}

	if !isCommand {
		return
	}

	split := strings.Split(e.Content, " ")
	root := strings.TrimPrefix(split[0], usedPrefix)

	var args []string
	if len(split) > 1 {
		args = split[1:]
	}

	ctx := CommandContext{
		MessageCreate: e,
		Shard: s,
		Args: args,
	}

	for _, cmd := range h.commands {
		match := strings.ToLower(cmd.Name) == strings.ToLower(root)

		// check aliases
		if !match {
			for _, alias := range cmd.Aliases {
				if strings.ToLower(cmd.Name) == strings.ToLower(alias) {
					match = true
					break
				}
			}
		}

		if match {
			// check subcommands
			argloop:
			for i, arg := range args {
				for _, sub := range cmd.SubCommands {
					subMatch := strings.ToLower(arg) == strings.ToLower(sub.Name)

					// check aliases
					if !subMatch {
						for _, alias := range sub.Aliases {
							if strings.ToLower(arg) == strings.ToLower(alias) {
								subMatch = true
								break
							}
						}
					}

					if subMatch {
						cmd = sub
						args = args[i:]
						ctx.Args = args
					} else {
						break argloop
					}
				}
			}

			go cmd.Handler(ctx)
		}
	}
}
