package ntpPrintTime

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func PrintNowTime() {
	nowTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ntp-server error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("now time is %s\n", nowTime.Format(time.RFC3339Nano))
}
