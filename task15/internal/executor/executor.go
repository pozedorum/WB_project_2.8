package executor

import (
	"task15/internal/builtins"
	"task15/internal/core"
)

type Executor struct {
	builtins  *builtins.Registry
	env       core.Environment
	processes []ProcessInfo // Для управления
}

func (e *Executor) Execute(cmd *core.Command) error {
	// 1. Проверка на builtin
	// 2. Настройка редиректов
	// 3. Запуск процесса
	// 4. Обработка pipe/условий
}

func (e *Executor) setupRedirects(cmd *core.Command) error {
	// Открыть файлы, подменить stdin/stdout
}

func (e *Executor) runPipeline(cmd *core.Command) error {
	// Создать pipe, запустить команды цепочкой
}
