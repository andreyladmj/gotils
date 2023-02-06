package display

import (
	"fmt"
	"github.com/andreyladmj/gotils/internal"
	"github.com/andreyladmj/gotils/internal/system"
	"strings"
	"sync"
)

type CommandLine struct {
	wg                *sync.WaitGroup
	maxGoroutines     uint8
	currentGoroutines uint8
	rootDir           *system.File
	filters           *internal.Filters
}

func NewCommandLine() *CommandLine {
	return &CommandLine{}
}

func (cl *CommandLine) Render(dir *system.File) {
	cl.printDir(dir, "-")
}

func (cl *CommandLine) printDir(file *system.File, tab string) {
	if file.IsDir {
		fmt.Println(tab+file.Name, fmt.Sprintf("(%d files):", len(file.Files)))

		for _, file := range file.Files {
			cl.printDir(file, tab+"-")
		}
	} else {
		fmt.Println(tab+file.Name, fmt.Sprintf("%.3fMb", Format(float64(file.Size), "mb")))
	}
}

func Format(n float64, t string) float64 {
	switch strings.ToLower(t) {
	case "kb":
		n = n / 1024
	case "mb":
		n = n / (1024 * 1024)
	case "gb":
		n = n / (1024 * 1024 * 1024)
	case "tb":
		n = n / (1024 * 1024 * 1024 * 1024)
	}

	return n
}
