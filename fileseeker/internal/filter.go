package internal

import (
	"github.com/andreyladmj/gotils/internal/system"
	"strings"
)

type Filters struct {
	filters []func(f *system.File) bool
}

func (f *Filters) Append(fn func(f *system.File) bool) {
	f.filters = append(f.filters, fn)
}

func (f *Filters) Filter(file *system.File) bool {
	for _, fn := range f.filters {
		if !fn(file) {
			return false
		}
	}

	return true
}

func (f *Filters) NameFilter(name string) {
	f.Append(func(file *system.File) bool {
		return strings.Contains(file.Name, name)
	})
}

func (f *Filters) FileSizeFilter(n int, t string) {
	switch strings.ToLower(t) {
	case "kb":
		n = 1024 * n
	case "mb":
		n = 1024 * 1024 * n
	case "gb":
		n = 1024 * 1024 * 1024 * n
	case "tb":
		n = 1024 * 1024 * 1024 * 1024 * n
	}

	f.Append(func(file *system.File) bool {
		if !file.IsDir {
			return file.Size > int(n)
		}
		return true
	})
}
