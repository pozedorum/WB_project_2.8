package parcer

type Command struct {
	Name      string
	Args      []string
	Redirects []Redirect // >, <, >>
	PipeTo    *Command   // |
	AndNext   *Command   // &&
	OrNext    *Command   // ||
}

type Redirect struct {
	Type string // ">", "<", ">>", "<<"
	File string // имя файла
	IsFd bool   // если true, File — это дескриптор (например, "2>&1")
}

const (
	// Основные операторы управления потоком
	Pipe       = "|"  // Конвейер
	And        = "&&" // Логическое И
	Or         = "||" // Логическое ИЛИ
	Semicolon  = ";"  // Разделитель команд
	Background = "&"  // Фоновое выполнение

	// Операторы перенаправления ввода/вывода
	RedirectOut     = ">"  // Перезапись файла
	RedirectIn      = "<"  // Чтение из файла
	RedirectAppend  = ">>" // Дописывание в файл
	RedirectHereDoc = "<<" // Here-document
	RedirectFdOut   = ">&" // Перенаправление дескриптора (вывод)
	RedirectFdIn    = "<&" // Перенаправление дескриптора (ввод)

	// Группировка команд
	SubshellOpen  = "(" // Начало подсекции
	SubshellClose = ")" // Конец подсекции
	BlockOpen     = "{" // Начало блока
	BlockClose    = "}" // Конец блока

	// Специальные операторы
	Not     = "!" // Логическое НЕ
	Comment = "#" // Комментарий
)

// Порядок от высшего к низшему
var operatorPrecedence = map[string]int{
	SubshellOpen:  5,
	SubshellClose: 5,
	Pipe:          4,
	RedirectOut:   3,
	RedirectIn:    3,
	And:           2,
	Or:            2,
	Semicolon:     1,
	Background:    1,
}

// Управляющие операторы
var controlOperators = map[string]bool{
	Pipe:       true,
	And:        true,
	Or:         true,
	Semicolon:  true,
	Background: true,
}

// Операторы перенаправления
var redirectOperators = map[string]bool{
	RedirectOut:     true,
	RedirectIn:      true,
	RedirectAppend:  true,
	RedirectHereDoc: true,
}
