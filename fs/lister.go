package fs

import (
	"io"
	"os"

	"github.com/kopia/kopia/cas"
)

const (
	directoryReadAhead = 1024
)

// Lister lists contents of filesystem directories.
type Lister interface {
	List(path string) (Directory, error)
}

type filesystemLister struct {
}

type localStreamingDirectory struct {
	dir     *os.File
	pending []os.FileInfo
}

func (d *filesystemLister) List(path string) (Directory, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var dir Directory

	for {
		fileInfos, err := f.Readdir(16)
		for _, fi := range fileInfos {
			dir = append(dir, entryFromFileSystemInfo(fi))
		}
		if err == nil {
			continue
		}
		if err == io.EOF {
			break
		}
		return nil, err
	}

	return dir, nil
}

type filesystemEntry struct {
	os.FileInfo

	objectID cas.ObjectID
}

func (fse *filesystemEntry) Size() int64 {
	if fse.Mode().IsRegular() {
		return fse.FileInfo.Size()
	}

	return 0
}

func (fse *filesystemEntry) ObjectID() cas.ObjectID {
	return fse.objectID
}

func entryFromFileSystemInfo(fi os.FileInfo) Entry {
	return &filesystemEntry{
		FileInfo: fi,
	}
}
