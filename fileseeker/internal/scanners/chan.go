package scanners

import (
	"fmt"
	"github.com/andreyladmj/gotils/internal"
	"github.com/andreyladmj/gotils/internal/system"
	"io/ioutil"
	"path"
	"sync"
)

type DirsChanScanner struct {
	wg                *sync.WaitGroup
	maxGoroutines     uint8
	currentGoroutines uint8
	rootDir           *system.File
	filters           *internal.Filters
}

func NewDirsChanScanner(maxGoroutines uint8) *DirsChanScanner {
	return &DirsChanScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines}
}

func (ds *DirsChanScanner) SetFilter(f *internal.Filters) {
	ds.filters = f
}

func (ds *DirsChanScanner) Scan(path string) {
	ds.rootDir = &system.File{Path: path}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}

func (ds *DirsChanScanner) ReadDir(dir *system.File) {
	ch := make(chan *system.File)

	go func(dirname string, ch chan *system.File) {
		ds.wg.Add(1)
		defer ds.wg.Done()

		files, err := ioutil.ReadDir(dir.Path)
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}

		for _, f := range files {
			fullpath := path.Join(dir.Path, f.Name())

			if f.IsDir() {
				d := system.NewFile(f.Name(), fullpath, int(f.Size()), f.IsDir())

				if ds.filters != nil && !ds.filters.Filter(d) {
					continue
				}

				if ds.currentGoroutines < ds.maxGoroutines {
					ds.currentGoroutines++
					ds.wg.Add(1)
					go func(d1 *system.File) {
						ds.ReadDir(d1)
						ds.currentGoroutines--
						ds.wg.Done()
					}(d)
				} else {
					ds.ReadDir(d)
				}
				ch <- d
			} else {
				ch <- system.NewFile(f.Name(), fullpath, int(f.Size()), f.IsDir())
			}
		}
		close(ch)

	}(dir.Path, ch)

	for file := range ch {
		dir.Append(file)
	}
}
