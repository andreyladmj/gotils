package system

import (
	"fmt"
	"os"
	"sync"
)

type File struct {
	Name  string
	Path  string
	Size  int
	IsDir bool
	Files []*File
	mu    sync.Mutex
}

func (f *File) Append(dir *File) {
	f.mu.Lock()
	f.Files = append(f.Files, dir)
	f.mu.Unlock()
}

func (f *File) GetAllFilesCount() int {
	c := 0
	for _, file := range f.Files {
		if file.IsDir {
			c += file.GetAllFilesCount()
		} else {
			c++
		}

	}
	return c
}
func (f *File) GetSize() int {
	c := 0
	for _, file := range f.Files {
		if file.IsDir {
			c += file.GetSize()
		} else {
			c += file.Size
		}

	}
	return c
}

func (f *File) DropEmptyDirs() {
	f.mu.Lock()
	for i := 0; i < len(f.Files); i++ {
		if f.Files[i].GetAllFilesCount() == 0 {
			f.Files = append(f.Files[:i], f.Files[i+1:]...)
			i--
		} else {
			f.Files[i].DropEmptyDirs()
		}
	}
	f.mu.Unlock()
}

func (d *File) String() string {
	sf := ""
	sd := ""
	for _, f := range d.Files {
		if f.IsDir {
			sd += fmt.Sprintf("- %s\n", f.String())
		} else {
			sf += fmt.Sprintf("- %s\n", f.String())
		}
	}

	return sd + sf
}

func (d *File) GetFilesCount() int {
	c := 0
	for _, file := range d.Files {
		if !file.IsDir {
			c++
		}
	}
	return c
}

func (d *File) GetDirsCount() int {
	c := 0
	for _, file := range d.Files {
		if file.IsDir {
			c++
		}
	}
	return c
}

func (d *File) CalcAllObjects() int {
	n := 0

	for _, f := range d.Files {
		if f.IsDir {
			n += f.CalcAllObjects() + 1
		} else {
			n++
		}
	}

	return n
}

func NewFileFromFileInfo(f os.FileInfo, fullpath string) *File {
	return &File{Name: f.Name(), Path: fullpath, Size: int(f.Size()), IsDir: f.IsDir()}
}

func NewFile(name string, fullpath string, size int, isDir bool) *File {
	return &File{Name: name, Path: fullpath, Size: size, IsDir: isDir}
}
