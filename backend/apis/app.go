package apis

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	conf := a.initConf()
	a.SetAppSet(conf)
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

func (a *App) GetFolderPath() string {
	result, err := runtime.OpenDirectoryDialog(a.Ctx, runtime.OpenDialogOptions{
		Title: "Open Folder",
	})
	// 设置回参
	if err != nil {
		fmt.Println(err)
	}
	if result == "" {
		return ""
	}
	return result
}

func (a *App) RunRelease(targetReleaseFilePath string) {
	cwd := filepath.Dir(targetReleaseFilePath)
	_, file_name := filepath.Split(targetReleaseFilePath)
	fmt.Println(cwd, file_name)
	RunRelease(cwd, file_name)
}

func (a *App) RunReleaseBeforeBuild(targetReleaseFilePath, build_cmd string) (string,error) {
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
		return "", err
	}
	RunRelease(cwd, file_name)
	return hash,nil
}

var confLock sync.RWMutex

func (a *App) LoadConf() (string, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	logDir := AppSet["LOG_DIR"]
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
	logDir := AppSet["LOG_DIR"]
	conf_path := filepath.Join(logDir, "tars-release-conf.json")
	err := os.WriteFile(conf_path, []byte(conf), 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *App) CheckBuildLog(filePath string) []string {
	logDir := AppSet["LOG_DIR"]
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

func PackEnviron() []string {
	rsp := make([]string, 0)
	rsp = append(rsp, os.Environ()...)
	rsp = append(rsp,"PATH=" + os.Getenv("PATH") + ":/usr/local/bin")
   return rsp
}

func (a *App) OpenProject(filePath string) string {
	cwd := filepath.Dir(filePath)
	var cmd *exec.Cmd = exec.Command("marscode", "./")
	cmd.Dir = cwd
	cmd.Env = append(cmd.Env, PackEnviron()...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		es := fmt.Sprintf("Error: %v \n Output: %s \n OsEnv: %v", err.Error(), output, cmd.Env)
		return es
	}
	return ""
}

var AppSet = make(map[string]string)

func (a *App) GetAppSet() map[string]string {
	cwd, _ := os.Getwd()
	confPath := filepath.Join(cwd, "etc", "yallow.conf")
	AppSet["confPath"] = confPath
	AppSet["cwd"] = cwd
	return AppSet
}

func (a *App) SetAppSet(sets map[string]string) {
	cwd, _ := os.Getwd()
	sets["cwd"] = cwd
	confPath := filepath.Join(cwd, "etc", "yallow.conf")
	sets["confPath"] = confPath
	AppSet = sets
	err := writeAppSetToFile(AppSet)
	if err != nil {
		fmt.Println("写入 AppSet 到文件失败:", err)
	}
}
func writeAppSetToFile(appSet map[string]string) error {
	// 定义文件路径
	cwd, _ := os.Getwd()
	filePath := filepath.Join(cwd, "etc", "yallow.conf")

	// 打开文件，如果文件不存在则创建，存在则截断（覆盖）
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将 map 数据转换为 key = value 的形式并写入文件
	for key, value := range appSet {
		_, err := fmt.Fprintf(file, "%s = %s\n", key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConf() map[string]string {
	cwd, _ := os.Getwd()
	file, err := os.Open(filepath.Join(cwd, "etc", "yallow.conf"))
	config := make(map[string]string)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Println("yallow.conf 不存在")
		return config
	}
	defer file.Close()
	defer debugConfig(config)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return config
}
