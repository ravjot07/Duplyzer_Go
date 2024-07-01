// shared/types.go
package shared

type Pair struct {
	Hash string
	Path string
}

type FileList []string
type Results map[string]FileList
