package ioutils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func PureName(name string) string {
	return strings.TrimSuffix(name, path.Ext(name))
}

func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	} else {
		return info.IsDir()
	}
}

// 下载文件
func DownloadWithDirAndName(ctx context.Context, srcUrl string, targetDir string, targetName string, funcs ...func(length, downloaded int64)) (err error) {
	var fsize int64 = 0
	var buf = make([]byte, 32*1024)
	var written int64
	err = nil

	req, err := http.NewRequestWithContext(ctx, "GET", srcUrl, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var filename string
	if targetName != "" {
		filename = targetName
	} else {
		uri, errSub := url.ParseRequestURI(srcUrl)
		if errSub != nil {
			err = errSub
			return
		}
		filename = path.Base(uri.Path)
	}

	// 读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return
	}
	targetPath := filepath.Join(targetDir, filename)
	if Exists(targetPath) {
		err = fmt.Errorf("%v already exists", targetPath)
		return
	}

	tmpFilePath := targetPath + ".download"

	// 如果目录不存在，则先创建目录
	if !IsDir(targetDir) {
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return
		}
	}

	// 创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if resp.Body == nil {
		return errors.New("resp.Body is null")
	}
	defer resp.Body.Close()

	// 下面是 io.copyBuffer() 的简化版本
	written = 0
	for {
		// 读取bytes
		cnt, readErr := resp.Body.Read(buf)
		if cnt > 0 {
			// 写入bytes
			cntWrited, writeErr := file.Write(buf[0:cnt])
			// 数据长度大于0
			if cntWrited > 0 {
				written += int64(cntWrited)
			}

			// 写入出错
			if writeErr != nil {
				err = writeErr
				break
			}

			// 读取是数据长度不等于写入的数据长度
			if cnt != cntWrited {
				err = io.ErrShortWrite
				break
			}
		}

		if readErr != nil {
			if readErr != io.EOF {
				err = readErr
			}
			break
		}

		// 没有错误则调用 callback
		for _, cb := range funcs {
			cb(fsize, written)
		}
	}

	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, targetPath)
	}

	return
}

// 将文件下载存储到指定路径
func DownloadWithPath(ctx context.Context, srcUrl string, targetPath string, funcs ...func(length, downloaded int64)) (err error) {
	parent, name := filepath.Split(targetPath)
	return DownloadWithDirAndName(ctx, srcUrl, parent, name, funcs...)
}

// 将文件下载存储到指定目录
func DownloadWithDir(ctx context.Context, srcUrl string, targetDir string, funcs ...func(length, downloaded int64)) (err error) {
	return DownloadWithDirAndName(ctx, srcUrl, targetDir, "", funcs...)
}
