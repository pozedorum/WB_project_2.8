package options

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type FlagStruct struct {
	KFlag *int
	NFlag *bool
	RFlag *bool
	UFlag *bool
	MFlag *bool
	BFlag *bool
	CFlag *bool
	HFlag *bool
}

func ParseOptions() (*FlagStruct, []string) {
	var fs FlagStruct

	fs.KFlag = flag.IntP("k", "k", 1, "sort by column number N")
	fs.NFlag = flag.BoolP("n", "n", false, "try to interpret strings as numbers and sort by it")
	fs.RFlag = flag.BoolP("r", "r", false, "sort in reverse order")
	fs.UFlag = flag.BoolP("u", "u", false, "output only sorted unique string")
	fs.MFlag = flag.BoolP("M", "M", false, "sort by month")
	fs.BFlag = flag.BoolP("b", "b", false, "ignore trailing blanks")
	fs.CFlag = flag.BoolP("c", "c", false, "check if data is sorted")
	fs.HFlag = flag.BoolP("h", "h", false, "sort by numerical value, taking into account suffixes")

	// Переопределяем Usage для отображения только коротких флагов
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] input_file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Парсим аргументы
	flag.Parse()

	return &fs, flag.Args()
}

func (fs *FlagStruct) PrintFlags() {
	fmt.Println("flag k -", *(fs.KFlag))
	fmt.Println("flag n -", *(fs.NFlag))
	fmt.Println("flag r -", *(fs.RFlag))
	fmt.Println("flag u -", *(fs.UFlag))
	fmt.Println("flag m -", *(fs.MFlag))
	fmt.Println("flag b -", *(fs.BFlag))
	fmt.Println("flag c -", *(fs.CFlag))
	fmt.Println("flag h -", *(fs.HFlag))
}
