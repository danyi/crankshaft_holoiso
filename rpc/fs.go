package rpc

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// substituteHomeDir takes a path that might be prefixed with `~`, and returns
// the path with the `~` replaced by the user's home directory.
func substituteHomeDir(path string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	if strings.HasPrefix(path, "~") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

type FSService struct{}

type ListDirArgs struct {
	Path string `json:"path"`
}

type Content struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

type ListDirReply struct {
	Contents []Content `json:"contents"`
}

func (service *FSService) ListDir(r *http.Request, req *ListDirArgs, res *ListDirReply) error {
	path := substituteHomeDir(req.Path)

	c, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading dir", err)
		return err
	}

	for _, entry := range c {
		res.Contents = append(res.Contents, Content{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		})
	}

	return nil
}

type ReadFileArgs struct {
	Path string `json:"path"`
}

type ReadFileReply struct {
	Data string `json:"data"`
}

func (service *FSService) ReadFile(r *http.Request, req *ReadFileArgs, res *ReadFileReply) error {
	path := substituteHomeDir(req.Path)

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file", err)
		return err
	}

	res.Data = string(data)

	return nil
}
