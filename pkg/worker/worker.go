package worker

import (
	_ "embed"
	"errors"
	"io/fs"
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

	_, err := os.Stat(fullPath)
	switch {
	case err == nil:
		return fullPath

	case errors.Is(err, fs.ErrNotExist):
		break

	default:
		log.Fatal("could not check worker file path: ", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		log.Fatal("could not create worker bin file: ", err)
	}

	if _, err := f.Write(wokerBin); err != nil {
		log.Fatal("could not write worker bin file: ", err)
	}
	defer f.Close()

	if err := os.Chmod(fullPath, 0755); err != nil {
		log.Fatal("could not chmod: ", err)
	}

	w.fullPath = fullPath

	return fullPath
}
