package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

/*
 example:
  mnt: /home/jho/colo00
  pa: /home/jho/colo00/cores/bug_33166/a/test.log
  realPa: testrunner.log
*/
func getRealPath(mnt, pa, realPa string) (string, error) {
	if f, err := os.Open(realPa); err == nil {
		f.Close()
		return realPa, nil
	}

	// if the 1st char is not /, it's a local link
	if realPa[0] != '/' {
		// get the last index of / in pa
		i := strings.LastIndex(pa, "/")
		return pa[:i+1] + realPa, nil
	}

	fn := func(c rune) bool {
		return c == rune('/')
	}

	paths := strings.FieldsFunc(realPa, fn)

	//fmt.Printf("paths %v\n", paths)
	realP := ""
	for i, p := range paths {
		if i == 0 {
			realP = mnt
		}
		realP = realP + "/" + p

		tmp, err := os.Readlink(realP)
		//fmt.Printf("i %d p %s\n", i, p)
		//fmt.Printf("realP is %s\n", realP)
		//fmt.Printf("tmp is %s\n", tmp)

		if err != nil {
			// it's not a link
			continue
		}
		realP = mnt + "/" + tmp
	}

	return realP, nil
}

func SearchPath(root string, cfg *Conf) error {
	const (
		asupdata         = "/asupdata"
		asupdataInternal = "/asupdata-internal"
	)

	if cfg.MountPoint != "" {
		// isMnt = true
		root = cfg.MountPoint + root
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("eachPath: path %s", path))
		}

		// don't care about directory
		if info.IsDir() {
			return nil
		}

		fName := filepath.Base(path)

		// looking for pre-defined files to work on
		for idx, logf := range cfg.LogFiles {
			if (logf.Type == "zcore" && strings.Contains(fName, "zcore")) ||
				strings.HasPrefix(fName, logf.Prefix) {

				// use ReadLink to determine if it's a real file or symblic link
				realPath, err := os.Readlink(path)
				if err != nil {
					// file is not a link
					f, err := os.Open(path)
					if err != nil {
						return errors.Wrap(err, fmt.Sprintf("eachPath: open path %s", path))
					}
					f.Close()
					cfg.LogFiles[idx].DaFiles = append(cfg.LogFiles[idx].DaFiles, NewDaFile(path))

					continue
				}

				// symbolic link
				realPath, err = getRealPath(cfg.MountPoint, path, realPath)
				f, err := os.Open(realPath)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("eachPath: open \n  path %s\n  real %s", path, realPath))
				}
				f.Close()
				cfg.LogFiles[idx].DaFiles = append(cfg.LogFiles[idx].DaFiles, NewDaFileWithLink(path, realPath))
			}
		}
		return nil
	})
}

/*
func SearchPath(root string, cfg *Conf) error {
	// bug_30536
	const (
		asupdata         = "/asupdata"
		asupdataInternal = "/asupdata-internal"
	)

	var (
		asup, asupInt, realPath string
		isMnt                   bool
		err                     error
	)

	if cfg.MountPoint != "" {
		isMnt = true

		root = cfg.MountPoint + root

		// get real asup
		asup = cfg.MountPoint + asupdata
		realPath, err = os.Readlink(asup)
		if err != nil {
			errors.Wrap(err, fmt.Sprintf("eachPath: open path %s", asup))
		}
		asup = cfg.MountPoint + realPath

		// get real asupdataInternal
		asupInt = cfg.MountPoint + asupdataInternal
		realPath, err = os.Readlink(asupInt)
		if err != nil {
			errors.Wrap(err, fmt.Sprintf("eachPath: open path %s", asupInt))
		}
		asupInt = cfg.MountPoint + realPath
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("eachPath: path %s", path))
		}

		// don't care about directory
		if info.IsDir() {
			return nil
		}

		// we are dealing w/ files after this point
		var fPath, fName string

		idxLastSlash := strings.LastIndex(path, "/")
		if idxLastSlash == -1 {
			// files in the root directory
			fName = path
		} else {
			fPath = path[:idxLastSlash+1]
			fName = path[idxLastSlash+1:]
		}
		_ = fPath

		var (
			f        *os.File
			realPath string
		)

		for idx, logf := range cfg.LogFiles {
			if strings.HasPrefix(fName, logf.Prefix) {
				realPath, err = os.Readlink(path)

				if err != nil {
					// file is not a link
					f, err = os.Open(path)
					if err != nil {
						return errors.Wrap(err, fmt.Sprintf("eachPath: open path %s", path))
					}
					f.Close()

					cfg.LogFiles[idx].DaFiles = append(cfg.LogFiles[idx].DaFiles, NewDaFile(path))

					continue
				}

				// need to read the LinkPath
				if isMnt {
					//tmp, err := filepath.EvalSymlinks(path)
					//if err != nil {
					//	return errors.Wrap(err, fmt.Sprintf("EvalSymlinks: open path %s", path))
					//}
					//fmt.Printf("path: %s\n => evalsymlinks %s\n", path, tmp)

					if strings.HasPrefix(realPath, asupdataInternal+"/") {
						realPath = asupInt + realPath[len(asupdataInternal):]
						fmt.Println(realPath)
					} else if strings.HasPrefix(realPath, asupdata+"/") {
						realPath = asup + realPath[len(asupdata):]
						fmt.Println(realPath)
					} else {
						realPath = cfg.MountPoint + realPath
						fmt.Println(realPath)
					}
				}

				f, err = os.Open(realPath)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("eachPath: open path %s", realPath))
				}
				f.Close()

				cfg.LogFiles[idx].DaFiles = append(cfg.LogFiles[idx].DaFiles, NewDaFileWithLink(path, realPath))
				// defer f.Close()

			}
		}
		// fmt.Printf("path %s name %s\n", fPath, fName)
		return nil
	})

}
*/
