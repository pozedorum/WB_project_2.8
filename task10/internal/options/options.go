package options

import (
	"flag"
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

func GetFlags() *FlagStruct {
	var fs FlagStruct
	fs.KFlag = flag.Int("k", 0, "sort by column number N")
	fs.NFlag = flag.Bool("n", false, "")
	fs.RFlag = flag.Bool("r", false, "")
	fs.UFlag = flag.Bool("u", false, "")
	fs.MFlag = flag.Bool("m", false, "")
	fs.BFlag = flag.Bool("b", false, "")
	fs.CFlag = flag.Bool("c", false, "")
	fs.HFlag = flag.Bool("h", false, "")
	flag.Parse()
	return &fs
}
