package mod

import (
	fsnotify "gopkg.in/fsnotify.v1"
	"os"
	"ght/utils"
	"strings"
	"time"
	"path/filepath"
)

// 根据recursive值确定是否递归查找paths每个目录下的子目录。
func recursivePath(recursive bool, paths []string) []string {
	if !recursive {
		return paths
	}

	ret := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			utils.ColorLog("[ERROR] [1]遍历监视目录错误: %s", err)
		}

		//(BUG):不能监视隐藏目录下的文件
		if fi.IsDir() && strings.Index(path, "/.") < 0 {
			ret = append(ret, path)
		}
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			utils.ColorLog("[ERROR] [2]遍历监视目录错误: %s", err)
		}
	}

	return ret
}

// 确定文件path是否属于被忽略的格式。
func isIgnore(path string) bool {
	var	appCmd *exec.Cmd // appName的命令行包装引用，方便结束其进程。
	if appCmd != nil && appCmd.Path == path { // 忽略程序本身的监视
		return true
	}

	for _, ext := range b.exts {
		if len(ext) == 0 {
			continue
		}
		if ext == "*" {
			return false
		}
		if strings.HasSuffix(path, ext) {
			return false
		}
	}

	return true
}

func Watch(paths []string){

	//初始化监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.ColorLog("[ERROR] 初始化监视器失败:", err)
		os.Exit(2)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					utils.ColorLog("[INFO] 修改文件: %s \n", event.Name)
				}
				utils.ColorLog("[INFO] [ %s ]文件被修改 \n", event.Name)
			case err := <-watcher.Errors:
				utils.ColorLog("[ERROR] %s \n", err)
			}
		}
	}()

	utils.ColorLog("[INFO] 初始化文件监视器...\n")
	for _, path := range paths {
		utils.ColorLog("[TRAC] 监视文件夹: %s \n", path)
		err = watcher.Add(path)
		if err != nil {
			utils.ColorLog("[ERRO] 监视文件夹失败 %s \n", err)
			os.Exit(2)
		}
	}

	<-done

	//
}

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		utils.ColorLog("[ERRO] Fail to open file[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		utils.ColorLog("[ERRO] Fail to get file information[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}



