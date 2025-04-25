package errorhandling

import (
	"fmt"
	"os"
)

func ReportAndExit(line int, where, message string) {
	fmt.Printf("line %d Error at %s: %s \n", line, where, message)
	os.Exit(1)
}
