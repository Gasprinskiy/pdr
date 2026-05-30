package worker

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
)

//go:embed worker
var wokerBin []byte

type Woker struct {
	dir      string
	binName  string
	fullPath string
}

func NewWorker(dir, binName string) *Woker {
	return &Woker{
		dir:     dir,
		binName: binName,
	}
}

func (w *Woker) Init() string {
	fullPath := filepath.Join(w.dir, w.binName)

	if _, err := os.Stat(fullPath); err == nil {
		return fullPath
	}

	f, err := os.Create(fullPath)
	if err != nil {
		log.Panic("could not create worker bin file: ", err)
	}

	if _, err := f.Write(wokerBin); err != nil {
		log.Panic("could not write worker bin file: ", err)
	}
	w.fullPath = fullPath

	return fullPath
}
