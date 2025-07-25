// Package builtins реализует наобр стандартных команд и утилит
package builtins

import "io"

type Environment interface {
	Getwd() (string, error)         // Получить текущую директорию
	Chdir(dir string) error         // Сменить директорию
	Getenv(key string) string       // Получить переменную окружения
	Setenv(key, value string) error // Установить переменную
	Environ() []string              // Получить все переменные окружения
}

type BuiltinCommand interface {
	Name() string                                              // Возвращает имя команды
	Execute(args []string, env Environment, w io.Writer) error // Выполняет команду
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
