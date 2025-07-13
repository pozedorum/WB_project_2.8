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

func ParceOptions() (*FlagStruct, []string) {
	var fs FlagStruct

	flagset := flag.NewFlagSet("myflags", flag.ContinueOnError)
	flagset.SetNormalizeFunc(func(f *flag.FlagSet, name string) flag.NormalizedName {
		return flag.NormalizedName(name)
	})

	fs.KFlag = flagset.Int("k", 0, "sort by column number N")                                 //
	fs.NFlag = flagset.Bool("n", false, "try to interpret strings as numbers and sort by it") //
	fs.RFlag = flagset.Bool("r", false, "sort in reverse order")                              //
	fs.UFlag = flagset.Bool("u", false, "output only sorted unique string")
	fs.MFlag = flagset.Bool("m", false, "sort by month")          //
	fs.BFlag = flagset.Bool("b", false, "ignore trailing blanks") //
	fs.CFlag = flagset.Bool("c", false, "check if data is sorted")
	fs.HFlag = flagset.Bool("h", false, "sort by numerical value, taking into account suffixes (for example, K = kilobyte, M = megabyte).") //

	err := flagset.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return &fs, flagset.Args()
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
