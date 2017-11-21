package libmockery

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Same(a, b Content) error {
	if len(a) != len(b) {
		return errors.New("number of files has changed")
	}

	for file, hashA := range a {
		hashB, exists := b[file]
		if !exists {
			return errors.New("missing hash for file: " + file)
		}
		if hashA != hashB {
			return fmt.Errorf("file %s hash has changed %q => %q", file, hashA, hashB)
		}
	}

	return nil
}

type Content map[string]string

type Hasher interface {
	Hash() (Content, error)
}

type hasher struct {
	root string
}

func NewHasher(root string) Hasher {
	return &hasher{root: root}
}

func (h *hasher) Hash() (Content, error) {
	content := make(map[string]string, 10)

	err := filepath.Walk(h.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, f); err != nil {
			return err
		}

		content[path] = fmt.Sprintf("%x", hash.Sum(nil))

		return nil
	})

	return content, err
}
