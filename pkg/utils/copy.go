/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 12:25:00
 */

package utils

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFile 拷贝文件
func CopyFile(dst, src string) error {
	var srcFile *os.File
	var dstFile *os.File
	var err error

	if srcFile, err = os.Open(src); err != nil {
		return err
	}

	if dstFile, err = os.Create(dst); err != nil {
		return err
	}

	_, err = io.Copy(dstFile, srcFile)

	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	return err
}

// CopyDirectory 拷贝文件夹
func CopyDirectory(src, dst string) error {
	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(src); err != nil {
		return err
	}

	if fileInfo.IsDir() {
		// src 是文件夹，那么定义 dst 也是文件夹
		var dirEntry []os.DirEntry
		if dirEntry, err = os.ReadDir(src); err == nil {
			// 递归每一个文件
			for _, item := range dirEntry {
				if err = CopyDirectory(filepath.Join(src, item.Name()), filepath.Join(dst, item.Name())); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	} else {
		// src 是文件，那么创建 dst 的文件夹
		dir := filepath.Dir(dst)

		if _, err = os.Stat(dir); err != nil {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		}

		return CopyFile(dst, src)
	}

	return nil
}
