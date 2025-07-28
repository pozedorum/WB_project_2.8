// Package builtins реализует наобр стандартных команд и утилит
package builtins

import (
	"io"

	"task15/internal/core"
)

type BuiltinCommand interface {
	Name() string                                                                         // Возвращает имя команды
	Execute(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error // Выполняет команду
}

type Registry struct {
	commands map[string]BuiltinCommand
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]BuiltinCommand),
	}
}

func NewRegistryWithDefaults() *Registry {
	r := NewRegistry()
	r.Register(NewCdUtil())
	r.Register(NewEchoUtil())
	r.Register(NewKillUtil())
	r.Register(NewPsUtil())
	r.Register(NewPwdUtil())
	return r
}

func (r *Registry) Register(newCom BuiltinCommand) {
	r.commands[newCom.Name()] = newCom
}

func (r *Registry) GetCommand(name string) (BuiltinCommand, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// =====================================
// Функции для тестов

type genericCommand struct {
	name    string
	execute func(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error
}

func NewGenericCommand(
	name string,
	executeFunc func(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error,
) BuiltinCommand {
	return &genericCommand{
		name:    name,
		execute: executeFunc,
	}
}

func (g *genericCommand) Name() string {
	return g.name
}

func (g *genericCommand) Execute(args []string, env core.Environment, stdin io.Reader, stdout io.Writer) error {
	return g.execute(args, env, stdin, stdout)
}
