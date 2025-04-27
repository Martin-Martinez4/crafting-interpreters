package errorhandling

import (
	"fmt"
	"os"
)

func ReportError(line int, where, message string) {
	fmt.Printf("line %d Error at %s: %s \n", line, where, message)
}

func ReportAndExit(line int, where, message string) {
	ReportError(line, where, message)
	os.Exit(1)
}
