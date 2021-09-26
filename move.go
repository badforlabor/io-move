/**
 * Auth :   liubo
 * Date :   2021/9/26 15:41
 * Comment:
 */

package main

import (
	"github.com/badforlabor/gocrazy/crazyio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func moveFile(srcFullpath, dstFullpath string) error {
	var e = os.Rename(srcFullpath, dstFullpath)
	if e == nil {
		return nil
	}

	var f, e2 = os.Stat(srcFullpath)
	if e2 != nil {
		return e2
	}

	if f.IsDir() {
		return nil
	}
	_, e = copyFile(srcFullpath, dstFullpath)
	if e == nil {
		os.Remove(srcFullpath)
	}
	return e
}
func DelFolderIfEmpty(folder string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	if len(files) != 0 {
		return nil
	}

	err = os.Remove(folder)
	return err
}
func DelEmptyFolder(src string) {

	var folderlist = make([]string, 0, 4096)

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			folderlist = append(folderlist, path)
		}

		return nil
	})

	var cnt = len(folderlist)
	for i, _ := range folderlist {
		DelFolderIfEmpty(folderlist[cnt - i - 1])
	}
}
// 移动文件夹
func MovePath(srcFullpath, dstFullpath string, callback func(string, error)) error {
	var e = os.Rename(srcFullpath, dstFullpath)
	if e != nil {
		// 手动移动
		var f, e = os.Stat(srcFullpath)
		if e == nil {
			if f.IsDir() {
				filepath.Walk(srcFullpath, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() {
						var dst2 = filepath.Join(dstFullpath, path[len(srcFullpath):])
						var e2 = moveFile(path, dst2)
						if e2 == nil {

						} else {
							e = e2
						}
						callback(dst2, e)
					} else if len(path) > len(srcFullpath) {
						crazyio.Mkdir(filepath.Join(dstFullpath, path[len(srcFullpath):]))
					}
					return nil
				})

				// 删除空文件夹
				DelEmptyFolder(srcFullpath)
			}
		} else {
			_, e = copyFile(srcFullpath, dstFullpath)
			if e == nil {
				os.Remove(srcFullpath)
			}
			callback(dstFullpath, e)
		}



	} else {
		callback(dstFullpath, e)
	}
	return e
}
func copyFile(src, dst string) (w int64, err error) {
	w = 0
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dstSlices := strings.Split(dst, string(os.PathSeparator))
	dstSlicesLen := len(dstSlices)
	destDir := ""
	for i := 0; i < dstSlicesLen-1; i++ {
		destDir = destDir + dstSlices[i] + string(os.PathSeparator)
	}
	b, err := pathExists(destDir)
	if b == false {
		err = os.MkdirAll(destDir, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			return
		}
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}

	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
