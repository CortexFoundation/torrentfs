// golang folder copy
// replicate because of vendor issues
// can replace by 'import "github.com/otiai10/copy"'
package torrentfs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	tmpPermissionForDirectory = os.FileMode(0755)
	fcInst                    = folderCopier{}
)

type folderCopier struct{}

func GetFolderCopier() *folderCopier {
	return &fcInst
}

// fclose ANYHOW closes file,
// with asiging error raised during Close,
// BUT respecting the error already reported.
func (fc *folderCopier) fclose(f *os.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}

// chmod ANYHOW changes file mode,
// with asiging error raised during Chmod,
// BUT respecting the error already reported.
func (fc *folderCopier) chmod(dir string, mode os.FileMode, reported *error) {
	if err := os.Chmod(dir, mode); *reported == nil {
		*reported = err
	}
}

// fcopy is for just a file
func (fc *folderCopier) fcopy(src, dest string, info os.FileInfo) (err error) {
	if err = os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return
	}

	f, err := os.Create(dest)
	if err != nil {
		return
	}
	defer fc.fclose(f, &err)

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return
	}

	s, err := os.Open(src)
	if err != nil {
		return
	}
	defer fc.fclose(s, &err)

	var buf []byte = nil
	var w io.Writer = f
	// var r io.Reader = s
	_, err = io.CopyBuffer(w, s, buf)
	return
}

// dcopy is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
func (fc *folderCopier) dcopy(srcdir, destdir string, info os.FileInfo) (err error) {

	_, err = os.Stat(destdir)
	if err == nil {
		err = os.ErrExist
		return
	} else if err != nil && !os.IsNotExist(err) {
		return // Unwelcome error type...!
	}

	originalMode := info.Mode()

	// Make dest dir with 0755 so that everything writable.
	if err = os.MkdirAll(destdir, tmpPermissionForDirectory); err != nil {
		return
	}
	// Recover dir mode with original one.
	defer fc.chmod(destdir, originalMode, &err)

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())

		if err = fc.switchboard(cs, cd, content); err != nil {
			// If any error, exit immediately
			return
		}
	}

	return
}

// switchboard switches proper copy functions regarding file type, etc...
// If there would be anything else here, add a case to this switchboard.
func (fc *folderCopier) switchboard(src, dest string, info os.FileInfo) (err error) {
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return
	case info.IsDir():
		err = fc.dcopy(src, dest, info)
	case info.Mode()&os.ModeNamedPipe != 0:
		return
	default:
		err = fc.fcopy(src, dest, info)
	}

	return
}

// Copy copies src to dest, doesn't matter if src is a directory or a file.
func (fc *folderCopier) Copy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return fc.switchboard(src, dest, info)
}
