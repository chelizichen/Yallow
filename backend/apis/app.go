package apis

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	Ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx
	// initListenKeyboard()
}

// Greet returns a greeting for the given name
func (a *App) OpenTarsReleaseFile() string {
	result, err := runtime.OpenFileDialog(a.Ctx, runtime.OpenDialogOptions{
		Title: "Open File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Scan tars.*.release",
				Pattern:     "*.release",
			},
		},
	})
	// 设置回参
	if err != nil {
		fmt.Println(err)
	}
	if result == "" {
		return ""
	}
	cwd := filepath.Dir(result)
	_, file_name := filepath.Split(result)
	fmt.Println(cwd, file_name)
	file, _ := GetConf(cwd, file_name)
	file["filePath"] = result
	return file.toString()
}

func (a *App) RunRelease(targetReleaseFilePath string) {
	cwd := filepath.Dir(targetReleaseFilePath)
	_, file_name := filepath.Split(targetReleaseFilePath)
	fmt.Println(cwd, file_name)
	RunRelease(cwd, file_name)
}

func (a *App) RunReleaseBeforeBuild(targetReleaseFilePath, build_cmd string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	cwd := filepath.Dir(targetReleaseFilePath)
	_, file_name := filepath.Split(targetReleaseFilePath)
	fmt.Println(cwd, file_name)
	hash, err := RunBuild(cwd, build_cmd)
	if err != nil {
		panic(err)
	}
	RunRelease(cwd, file_name)
	return hash
}

var confLock sync.RWMutex

func (a *App) LoadConf() (string, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	logDir := os.Getenv("LOG_DIR")
	conf_path := filepath.Join(logDir, "tars-release-conf.json")
	if _, err := os.Stat(conf_path); os.IsNotExist(err) {
		file, err := os.Create(conf_path)
		if err != nil {
			fmt.Println("创建文件失败:", err)
		}
		defer file.Close()
		return "", nil
	}
	s, err := os.OpenFile(conf_path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	b, err := io.ReadAll(s)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(b), nil
}

func (a *App) MergeConf(conf string) (bool, error) {
	confLock.Lock()
	defer confLock.Unlock()
	logDir := os.Getenv("LOG_DIR")
	conf_path := filepath.Join(logDir, "tars-release-conf.json")
	err := os.WriteFile(conf_path, []byte(conf), 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *App) CheckBuildLog(filePath string) []string {
	logDir := os.Getenv("LOG_DIR")
	cwd := filepath.Dir(filePath)
	hash := CurrentHashs[cwd]
	conf_path := filepath.Join(logDir, fmt.Sprintf("tars-release-%s.log", hash))
	rsp := make([]string, 10)
	// 打开日志文件
	file, err := os.Open(conf_path)
	if err != nil {
		fmt.Printf("无法打开日志文件: %v \n", err)
	}
	defer file.Close()

	// 创建一个 scanner 来读取文件内容
	scanner := bufio.NewScanner(file)
	// 逐行读取并打印日志内容
	for scanner.Scan() {
		rsp = append(rsp, scanner.Text())
	}
	// 检查 scanner 是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取日志文件时发生错误: %v \n", err)
	}
	return rsp
}
