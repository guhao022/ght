package mod

import (
	"ght/utils"
	fsnotify "gopkg.in/fsnotify.v1"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"runtime"
)

var watchExts = []string{".go", ".php"}

var started = make(chan bool)


var (
	cmd          *exec.Cmd
	eventTime    = make(map[string]int64)
	scheduleTime time.Time
)

// 根据recursive值确定是否递归查找paths每个目录下的子目录。
func recursivePath(recursive bool, paths []string) []string {
	if !recursive {
		return paths
	}

	ret := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			utils.ColorLog("[ERRO] [1]遍历监视目录错误: %s", err)
		}

		//(BUG):不能监视隐藏目录下的文件
		if fi.IsDir() && strings.Index(path, "/.") < 0 {
			ret = append(ret, path)
		}
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			utils.ColorLog("[ERRO] [2]遍历监视目录错误: %s", err)
		}
	}

	return ret
}

func checkIfWatchExt(name string) bool {
	for _, s := range watchExts {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

func Watch(paths []string) {

	//初始化监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.ColorLog("[ERRO] 初始化监视器失败:", err)
		os.Exit(2)
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if !checkIfWatchExt(event.Name) {
					continue
				}
				mt := getFileModTime(event.Name)
				if t := eventTime[event.Name]; mt == t {
					continue
				}

				eventTime[event.Name] = mt
				utils.ColorLog("[SUCC] [ %s ]文件被修改 \n", event.Name)

				Qbuild()

			case err := <-watcher.Errors:
				utils.ColorLog("[ERRO] 监控失败 %s \n", err)
			}
		}
	}()

	utils.ColorLog("[INFO] 初始化文件监视器...\n")
	for _, path := range paths {
		utils.ColorLog("[SUCC] 监视文件夹: [%s] \n", path)
		err = watcher.Add(path)
		if err != nil {

			utils.ColorLog("[ERRO] 监视文件夹失败: [%s] \n", err)
			os.Exit(2)
		}
	}

}

func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {

		utils.ColorLog("[ERRO] 文件打开失败--[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		utils.ColorLog("[ERRO] 获取不到文件信息--[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func getAppName(outputName, wd string) string {
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(outputName, ".exe") {
		outputName += ".exe"
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = /*wd + string(filepath.Separator) + */outputName
	}

	utils.ColorLog("[INFO] 输出文件为--[ %s ]\n", outputName)

	return outputName
}

// 开始编译代码
func Qbuild() {
	utils.ColorLog("[INFO] 编译代码...")

	goCmd := exec.Command("go", "build")
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	if err := goCmd.Run(); err != nil {
		utils.ColorLog("[ERRO] 编译失败: %s", err)
		return
	}

	utils.ColorLog("[SUCC] 编译成功...")

	restart()
}

// 重启被编译的程序
func restart() {
	defer func() {
		if err := recover(); err != nil {
			utils.ColorLog("[ERRO] restart.defer: %s", err)
		}
	}()

	var	appCmd *exec.Cmd
	wd, err := os.Getwd()
	if err != nil {
		utils.ColorLog("[ERROR] 获取当前工作目录时，发生以下错误:", err)
		return
	}
	appname := getAppName("", wd)
	// kill process
	if appCmd != nil && appCmd.Process != nil {
		utils.ColorLog("[INFO] 中止旧进程...")
		if err := appCmd.Process.Kill(); err != nil {
			utils.ColorLog("[ERRO] kill err: %s", err)
		}
		utils.ColorLog("[SUCC] 旧进程被终止!")
	}

	utils.ColorLog("[INFO] 启动新进程...")

	appCmd = exec.Command(appname)
	appCmd.Stderr = os.Stderr
	appCmd.Stdout = os.Stdout
	if err := appCmd.Start(); err != nil {
		utils.ColorLog("[ERRO] 启动进程时出错: %s", err)
	}
}


// 编译代码
func Build() {
	utils.ColorLog("[INFO] 开始编译... \n")

	goCmd := exec.Command("go", "build")
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	if err := goCmd.Run(); err != nil {
		utils.ColorLog("[INFO] 编译失败: %s \n", err)
		return
	}

	utils.ColorLog("[INFO] 编译成功... \n")
	Restart()
}

/*func Kill() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Kill.recover -> ", e)
		}
	}()
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("Kill -> ", err)
		}
	}
}*/

func Restart() {
	utils.ColorLog("[INFO] 终止程序... \n")
	go Start()
}

func Start() {
	var	appCmd *exec.Cmd
	wd, err := os.Getwd()
	if err != nil {
		utils.ColorLog("[ERROR] 获取当前工作目录时，发生以下错误:", err)
		return
	}

	appname := getAppName("", wd)
	utils.ColorLog("[INFO] 重启 %s ...\n", appname)
	if strings.Index(appname, "./") == -1 {
		appname = "./" + appname
	}

	if appCmd != nil && appCmd.Process != nil {
		utils.ColorLog("[INFO] 进程终止... \n")
		if err := appCmd.Process.Kill(); err != nil {
			utils.ColorLog("[ERROR] 终止进程失败 %s ...\n", err)
		}
		utils.ColorLog("[SUCC] 旧进程被终止! \n")
	}

	utils.ColorLog("[INFO] 启动新进程... \n")
	appCmd = exec.Command(appname)
	appCmd.Stderr = os.Stderr
	appCmd.Stdout = os.Stdout
	if err := appCmd.Start(); err != nil {
		utils.ColorLog("[ERRO] 启动进程时出错: %s \n", err)
	}


	//go cmd.Run()
	utils.ColorLog("[SUCC] 新进程已经启动...\n")
	started <- true
}
