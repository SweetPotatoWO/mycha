package internal

import (
	"fmt"
	"mycha/helper/log"
	"os"
	"path/filepath"
	)

var logger = log.DLogger()

func checkDirPath(dirPath string) (absDirPath string,err error) {
	if dirPath == "" {
		err = fmt.Errorf("空的文件的路径:%s",dirPath)
		return
	}
	//涉及到一个文件的操作类？？
	if filepath.IsAbs(dirPath) {
		absDirPath = dirPath
	} else {
		absDirPath, err = filepath.Abs(dirPath)
		if err != nil {
			return
		}
	}
	var dir  *os.File
	dir,err = os.Open(absDirPath)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if  dir ==  nil {
		err = os.MkdirAll(absDirPath,0700)
		if err != nil && !os.IsExist(err) {
			return
		}
	} else {
		var fileInfo os.FileInfo
		fileInfo, err = dir.Stat()
		if err != nil {
			return
		}
		if !fileInfo.IsDir() {
			err = fmt.Errorf("不存在文件夹:%s",absDirPath)
			return
		}
	}
	return
}

func Record(level byte,content string) {
	if content == "" {
		return
	}
	switch level {
	case 0:
		logger.Infoln(content)
	case 1:
		logger.Infoln(content)
	case 2:
		logger.Infoln(content)
	}
}






















