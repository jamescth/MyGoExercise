// https://github.com/tmc/scp/blob/master/scp.go
package scp

import (
	"fmt"
	"io"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

func Copy(size int64, mode os.FileMode, fileName string,
	content io.Reader, destinationPath string, session *ssh.Session) error {
	return copy(size, mode, fileName, content, destinationPath, session)
}

func CopyPath(filePath, destinationPath string, session *ssh.Session) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	return copy(s.Size(), s.Mode().Perm(), path.Base(filePath), f, destinationPath, session)
}

func copy(size int64, mode os.FileMode, fileName string,
	content io.Reader, destination string, session *ssh.Session) error {
	defer session.Close()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fmt.Fprintf(w, "C%#o %d %s\n", mode, size, fileName)
		io.Copy(w, content)
		fmt.Fprint(w, "\x00")
	}()
	cmd := fmt.Sprintf("scp -t %s", destination)
	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}
