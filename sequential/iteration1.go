// sequential/iteration1.go
package sequential

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Pair struct {
	Hash string
	Path string
}

type FileList []string
type Results map[string]FileList

func hashFile(path string) Pair {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return Pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

func SearchTreeSequential(dir string) (Results, error) {
	hashes := make(Results)
	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			h := hashFile(p)
			hashes[h.Hash] = append(hashes[h.Hash], h.Path)
		}
		return nil
	}
	err := filepath.Walk(dir, visit)
	return hashes, err
}
