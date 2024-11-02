package printer

import (
	"fmt"
	"os"
)

func Stdoutln(format string, args ...any) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func Stderrln(format string, args ...any) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
}
