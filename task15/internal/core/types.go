package core

type Command struct {
	Name      string
	Args      []string
	Redirects []Redirect // >, <, >>
	PipeTo    *Command   // |
	AndNext   *Command   // &&
	OrNext    *Command   // ||
}

func (c *Command) IsEmpty() bool {
	return c.Name == "" &&
		len(c.Args) == 0 &&
		len(c.Redirects) == 0 &&
		c.PipeTo == nil &&
		c.AndNext == nil &&
		c.OrNext == nil
}

type Redirect struct {
	Type string // ">", "<", ">>", "<<"
	File string // имя файла
}

const (
	// Основные операторы управления потоком
	Pipe = "|"  // Конвейер
	And  = "&&" // Логическое И
	Or   = "||" // Логическое ИЛИ

	// Операторы перенаправления ввода/вывода
	RedirectOut = ">" // Перезапись файла
	RedirectIn  = "<" // Чтение из файла

)

// Порядок от высшего к низшему
var OperatorPrecedence = map[string]int{
	Pipe:        4,
	RedirectOut: 3,
	RedirectIn:  3,
	And:         2,
	Or:          2,
}

// Управляющие операторы
var ControlOperators = map[string]bool{
	Pipe: true,
	And:  true,
	Or:   true,
}

// Операторы перенаправления
var RedirectOperators = map[string]bool{
	RedirectOut: true,
	RedirectIn:  true,
}
