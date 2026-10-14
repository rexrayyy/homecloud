package client

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/sftp"
)

// Search limits keep a recursive walk from stalling the UI on huge trees.
const (
	maxSearchResults = 500
	maxSearchDepth   = 8
	maxSearchVisited = 50000
)

// searchSkip lists virtual or very expensive mount points that are never walked.
var searchSkip = map[string]bool{
	"/proc": true, "/sys": true, "/dev": true, "/run": true,
	"/snap": true, "/lost+found": true,
}

// SearchEntry is one search hit, absolute remote path included.
type SearchEntry struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

func (a *App) sftpClient() (*sftp.Client, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn == nil {
		return nil, fmt.Errorf("сначала подключитесь по SSH")
	}
	if a.sftp != nil {
		return a.sftp, nil
	}
	c, err := sftp.NewClient(a.conn)
	if err != nil {
		return nil, fmt.Errorf("SFTP-подсистема недоступна на сервере: %w", err)
	}
	a.sftp = c
	return c, nil
}

// safeRemotePath normalises an absolute POSIX path and blocks traversal.
// Spaces are preserved because they are legal in Linux file names.
func safeRemotePath(p string) (string, error) {
	if strings.ContainsAny(p, "\x00\r\n") {
		return "", fmt.Errorf("недопустимые символы в пути")
	}
	if p == "" {
		return "/", nil
	}
	for _, part := range strings.Split(p, "/") {
		if part == ".." {
			return "", fmt.Errorf("обход каталога запрещён")
		}
	}
	clean := path.Clean("/" + strings.TrimPrefix(p, "/"))
	if clean == "." || clean == "" {
		return "/", nil
	}
	return clean, nil
}

// HomeDir returns the remote user's home directory for a quick jump.
func (a *App) HomeDir() (string, error) {
	cl, err := a.sftpClient()
	if err != nil {
		return "", err
	}
	wd, err := cl.Getwd()
	if err != nil {
		return "", err
	}
	return safeRemotePath(wd)
}

// ListFiles returns the real contents of an absolute directory over SFTP.
func (a *App) ListFiles(p string) ([]FileEntry, error) {
	dir, err := safeRemotePath(p)
	if err != nil {
		return nil, err
	}
	cl, err := a.sftpClient()
	if err != nil {
		return nil, err
	}
	entries, err := cl.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать %s: %w", dir, err)
	}
	out := make([]FileEntry, 0, len(entries))
	for _, fi := range entries {
		name := fi.Name()
		if name == "." || name == ".." {
			continue
		}
		e := FileEntry{Name: name, Path: path.Join(dir, name), IsDir: fi.IsDir(), Size: fi.Size()}
		if t := fi.ModTime(); !t.IsZero() {
			e.Modified = t.Local().Format("2006-01-02 15:04")
		}
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// SearchFiles walks the remote tree under dir and matches names case-insensitively.
func (a *App) SearchFiles(dir, query string) ([]SearchEntry, error) {
	root, err := safeRemotePath(dir)
	if err != nil {
		return nil, err
	}
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, fmt.Errorf("введите поисковый запрос")
	}
	if len(q) > 200 || strings.ContainsAny(q, "\x00\r\n") {
		return nil, fmt.Errorf("недопустимый поисковый запрос")
	}
	cl, err := a.sftpClient()
	if err != nil {
		return nil, err
	}
	needle := strings.ToLower(q)
	out := make([]SearchEntry, 0, 64)
	visited := 0
	// walkDepth reads each directory level-by-level, applying depth and count
	// limits so huge trees cannot stall the UI or exhaust memory.
	var walkDepth func(dir string, depth int)
	walkDepth = func(dir string, depth int) {
		if len(out) >= maxSearchResults || visited >= maxSearchVisited {
			return
		}
		entries, err := cl.ReadDir(dir)
		if err != nil {
			return
		}
		for _, fi := range entries {
			if len(out) >= maxSearchResults || visited >= maxSearchVisited {
				return
			}
			name := fi.Name()
			if name == "." || name == ".." {
				continue
			}
			p := path.Join(dir, name)
			visited++
			if strings.Contains(strings.ToLower(name), needle) {
				e := SearchEntry{Name: name, Path: p, IsDir: fi.IsDir(), Size: fi.Size()}
				if t := fi.ModTime(); !t.IsZero() {
					e.Modified = t.Local().Format("2006-01-02 15:04")
				}
				out = append(out, e)
			}
			if fi.IsDir() {
				if depth+1 >= maxSearchDepth || searchSkip[p] {
					continue
				}
				walkDepth(p, depth+1)
			}
		}
	}
	walkDepth(root, 0)
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// DownloadFile copies a remote file into a user-chosen local folder.
func (a *App) DownloadFile(p string) (string, error) {
	src, err := safeRemotePath(p)
	if err != nil {
		return "", err
	}
	cl, err := a.sftpClient()
	if err != nil {
		return "", err
	}
	info, err := cl.Stat(src)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать %s: %w", src, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s is a directory", src)
	}
	remote, err := cl.Open(src)
	if err != nil {
		return "", err
	}
	defer remote.Close()
	dir, err := pickFolder("Choose download folder")
	if err != nil || dir == "" {
		return "", err
	}
	tmp, err := os.CreateTemp(dir, ".homecloud-download-*")
	if err != nil {
		return "", err
	}
	tmpName := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()
	if _, err = io.Copy(tmp, remote); err != nil {
		return "", err
	}
	if err = tmp.Close(); err != nil {
		return "", err
	}
	final := filepath.Join(dir, path.Base(src))
	if err = os.Rename(tmpName, final); err != nil {
		return "", err
	}
	return final, nil
}

// UploadFile copies files picked in a native dialog into the current directory.
func (a *App) UploadFile(directory string) ([]string, error) {
	return a.upload(directory, false)
}

// UploadFolder copies a local folder, including its subfolders, into the
// current remote directory, keeping the folder name and its structure.
func (a *App) UploadFolder(directory string) ([]string, error) {
	return a.upload(directory, true)
}

func (a *App) upload(directory string, folder bool) ([]string, error) {
	dir, err := safeRemotePath(directory)
	if err != nil {
		return nil, err
	}
	cl, err := a.sftpClient()
	if err != nil {
		return nil, err
	}
	// The filesystem root is almost never writable by a normal account, and the
	// file view opens there. Land those uploads in the login home instead of
	// failing with a permission error after the user has already picked files.
	if dir == "/" {
		if home, herr := cl.Getwd(); herr == nil {
			if home, herr = safeRemotePath(home); herr == nil && home != "/" && home != "" {
				dir = home
			}
		}
	}
	info, err := cl.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}
	var picked []string
	if folder {
		chosen, derr := pickFolder("Choose a folder to upload")
		if derr != nil || chosen == "" {
			return nil, derr
		}
		picked = []string{chosen}
	} else {
		picked, err = pickFiles("Choose files to upload")
		if err != nil || len(picked) == 0 {
			return nil, err
		}
	}
	uploaded := make([]string, 0, len(picked))
	for _, local := range picked {
		info, statErr := os.Stat(local)
		if statErr != nil {
			return uploaded, statErr
		}
		name := filepath.Base(local)
		if !validUploadName(name) {
			return uploaded, fmt.Errorf("недопустимое имя %q", local)
		}
		remotePath := path.Join(dir, name)
		if info.IsDir() {
			if err = uploadDir(cl, local, remotePath, &uploaded); err != nil {
				return uploaded, err
			}
			continue
		}
		if err = uploadOneFile(cl, local, remotePath); err != nil {
			return uploaded, err
		}
		uploaded = append(uploaded, remotePath)
	}
	return uploaded, nil
}

// validUploadName rejects names that would escape the destination directory.
func validUploadName(name string) bool {
	return name != "" && name != "." && name != ".." && !strings.ContainsAny(name, "/\\\x00\r\n")
}

// uploadDir recreates a local folder and everything inside it over SFTP.
func uploadDir(cl *sftp.Client, localDir, remoteDir string, uploaded *[]string) error {
	if err := cl.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("не удалось создать %s: %w", remoteDir, err)
	}
	entries, err := os.ReadDir(localDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if !validUploadName(name) {
			return fmt.Errorf("недопустимое имя %q", name)
		}
		local := filepath.Join(localDir, name)
		remote := path.Join(remoteDir, name)
		if entry.IsDir() {
			if err = uploadDir(cl, local, remote, uploaded); err != nil {
				return err
			}
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		// A symlink or other non-regular entry is skipped rather than followed,
		// so an upload can never escape the folder the user picked.
		if !info.Mode().IsRegular() {
			continue
		}
		if err = uploadOneFile(cl, local, remote); err != nil {
			return err
		}
		*uploaded = append(*uploaded, remote)
	}
	return nil
}

// uploadOneFile copies one local file, replacing a remote file of the same name.
func uploadOneFile(cl *sftp.Client, local, remotePath string) error {
	f, err := os.Open(local)
	if err != nil {
		return err
	}
	defer f.Close()
	remote, err := cl.Create(remotePath)
	if err != nil {
		return fmt.Errorf("не удалось создать %s: %w", remotePath, err)
	}
	defer remote.Close()
	if _, err = io.Copy(remote, f); err != nil {
		return err
	}
	return nil
}
