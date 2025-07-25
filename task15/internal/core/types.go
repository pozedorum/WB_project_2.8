package core

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
