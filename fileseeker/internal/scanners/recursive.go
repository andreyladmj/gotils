package scanners

import (
	"github.com/andreyladmj/gotils/internal"
	"github.com/andreyladmj/gotils/internal/system"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"sync"
)

type DirsRecursiveScanner struct {
	wg                *sync.WaitGroup
	maxGoroutines     uint8
	currentGoroutines uint8
	rootDir           *system.File
	filters           *internal.Filters
}

func NewDirsRecursiveScanner(maxGoroutines uint8) *DirsRecursiveScanner {
	return &DirsRecursiveScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines}
}

func (ds *DirsRecursiveScanner) SetFilter(f *internal.Filters) {
	ds.filters = f
}

func (ds *DirsRecursiveScanner) GetDir() *system.File {
	return ds.rootDir
}

func (ds *DirsRecursiveScanner) Scan(path string) {
	ds.rootDir = &system.File{Path: path, Name: path, IsDir: true}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}

func (ds *DirsRecursiveScanner) DropEmptyDirs() {
	ds.rootDir.DropEmptyDirs()
}

func (ds *DirsRecursiveScanner) ReadDir(dir *system.File) {
	files, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		log.Error(err)
		return
	}
	for _, f := range files {
		fullpath := path.Join(dir.Path, f.Name())

		d := system.NewFile(f.Name(), fullpath, int(f.Size()), f.IsDir())

		if ds.filters != nil && !ds.filters.Filter(d) {
			continue
		}

		if f.IsDir() {

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

			dir.Append(d)

		} else {
			dir.Append(d)
		}
	}
}
