package apis

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type T_Config map[string]string

func (t T_Config) toString() string {
	b, _ := json.Marshal(t)
	return string(b)
}

func (t T_Config) getApp() string {
	return t["APP"]
}

func (t T_Config) getServer() string {
	return t["SERVER"]
}

const (
	CONF_FILE = "tars.release"
	CONF_FLAG = "conf"
	CONF_DESC = "指定taf部署文件"

	DEPLOY_MSG     = "发布备注"
	DEPLOY_DEFAULT = "auto deploy"
	DEPLOY_FLAG    = "msg"

	PRE_FLAG     = "precmd"
	PRE_DESC     = "预处理命令(npm run release)，实验性质"
	EMPTY_STRING = ""

	TARS_LIST_FLAG = "list"
	TARS_LIST_DESC = "获取发布列表"

	TARS_DATA_FLAG = "data"
	TARS_DATA_DESC = "获取发布数据"
)
const (
	PATH_GetServerPatchList = "/pages/server/api/server_patch_list"
	PATH_ServerPatch        = "/pages/server/api/upload_and_publish"
	PATH_DownLoad           = "/pages/server/api/download_package"
)

const (
	DATA    = "data"
	SUCCESS = "success"
)

type servantLog struct {
	__logFile *os.File
	__logger  *log.Logger
}

var T_Log *servantLog
var T_Success_Log *servantLog

func (s *servantLog) close() {
	s.__logFile.Close()
}

func (s *servantLog) print(msg string, v ...any) {
	s.__logger.Printf(msg, v...)
}

func (s *servantLog) panic(format string, v ...any) {
	s.__logger.Panicf(format, v...)
}

type I_ServantLog_Args interface {
	getApp() string
	getServer() string
}

func NewServantLog(args I_ServantLog_Args, t string) *servantLog {
	app := args.getApp()
	name := args.getServer()
	if app == "" || name == "" {
		panic("module or name is empty")
	}
	logDir := AppSet["LOG_DIR"]
	if logDir == "" {
		panic("LOG_DIR environment variable is not set")
	}
	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s.%s.%s.log", app, name, t))
	fmt.Printf("open logFilePath: %s\n", logFilePath)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("无法打开日志文件: %+v \n", err)
	}
	logger := log.New(logFile, "", log.LstdFlags)
	return &servantLog{
		__logFile: logFile,
		__logger:  logger,
	}
}

func RandomHash() string {
	// 生成一个32字节的随机字节切片
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 将随机字节切片转换为Base64编码的字符串
	randomString := base64.StdEncoding.EncodeToString(randomBytes)
	return randomString
}

var CurrentHashs = make(map[string]string)

func RunBuild(cwd, build_cmd string) (string, error) {
	fmt.Println("debug.runbuild >> build_cmd", build_cmd)
	fmt.Println("debug.runbuild >> cwd", cwd)
	// if strings.Index(build_cmd, "npm")!= -1 {
	// 	// /Users/leemulus/.nvm/versions/node/v16.20.1/bin/npm
	// 	build_cmd = strings.ReplaceAll(build_cmd, "npm", "/Users/leemulus/.nvm/versions/node/v16.20.1/bin/npm")
	// }
	cmd := exec.Command("/bin/sh", "-c", build_cmd)
	cmd.Dir = cwd
	cmd.Env = append(cmd.Env, PackEnviron()...)
	logDir :=  AppSet["LOG_DIR"]
	if logDir == "" {
		panic("LOG_DIR environment variable is not set")
	}
	hash := RandomHash()
	CurrentHashs[cwd] = hash
	logPath := filepath.Join(logDir, fmt.Sprintf("tars-release-%s.log", hash))
	fmt.Println("conf_path", logPath)
	// 检查文件是否存在
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		err = nil
		// 如果文件不存在，创建一个新文件
		file, err := os.Create(logPath)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}
	// 打开文件以重定向输出流
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	defer os.Remove(logPath)
	cmd.Stdout = file
	cmd.Stderr = file
	err = cmd.Run()
	return hash, err
}

func RunRelease(cwd, conf_path string) {
	config, done := GetConf(cwd, conf_path)
	if done {
		fmt.Println(" process exit with getConf().done = true ")
		return
	}
	T_Log = NewServantLog(config, DATA)
	T_Success_Log = NewServantLog(config, SUCCESS)
	defer T_Log.close()
	defer T_Success_Log.close()
	patchRequest(config,cwd)
}

func GetConf(cwd, conf_path string) (T_Config, bool) {
	file, err := os.Open(filepath.Join(cwd, conf_path))

	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := make(T_Config)
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
	if config["TAF_SERVER_TICKET"] == "" {
		panic("TAF_SERVER_TICKET is not set")
	}
	TAF_SERVER_PATH, ok := config["TAF_SERVER_PATH"]
	if !ok {
		if config["TAF_PATCH_PATH"] == "" {
			panic("TAF_PATCH_PATH is not set , and TAF_SERVER_PATH is not set")
		}
	} else {
		config["TAF_PATCH_PATH"] = fmt.Sprintf("%s%s", TAF_SERVER_PATH, PATH_ServerPatch)
		config["TAF_SERVER_PATCH_LIST_PATH"] = fmt.Sprintf("%s%s", TAF_SERVER_PATH, PATH_GetServerPatchList)
		config["TAR_PACKAGE_DOWNLOAD_PATH"] = fmt.Sprintf("%s%s", TAF_SERVER_PATH, PATH_DownLoad)
	}
	config["COMMENT"] = fmt.Sprintf("[%s] >>  %s ", config["PROJECT_VERSION"], "AUTO_RELEASE")
	return config, false
}

func patchRequest(config T_Config,cwd string) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	ticket, ok := config["TAF_SERVER_TICKET"]
	if !ok {
		T_Log.panic("TAF_SERVER_TICKET is not set \n")
	}
	TAF_PATCH_PATH, ok := config["TAF_PATCH_PATH"]
	if !ok {
		T_Log.panic("TAF_PATCH_PATH is not set \n")
	}
	target_req_url := fmt.Sprintf("%s?ticket=%s", TAF_PATCH_PATH, ticket)
	fmt.Println("Request URL:", target_req_url)
	writer.WriteField("application", config.getApp())
	writer.WriteField("module_name", config.getServer())
	writer.WriteField("comment", config["COMMENT"])
	package_path := config["PACKAGE_PATH"]

	file, err := os.Open(filepath.Join(cwd, package_path))
	if err != nil {
		T_Log.panic("os.Open.error  %+v \n", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("suse", filepath.Join(cwd, package_path))
	if err != nil {
		T_Log.panic("writer.CreateFormFile.err %+v \n", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		T_Log.panic("io.copy.error %+v \n", err)
	}
	writer.Close()

	// // 发送 HTTP 请求
	req, err := http.NewRequest("POST", target_req_url, body)
	if err != nil {
		T_Log.panic("http.NewRequest.error %+v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		T_Log.panic("client.Doclient.Do.error %+v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	fmt.Printf("Invoke PathRequest Response status:  %s \n", resp.Status)
	var rspBody = make([]byte, 2048)
	_, err = resp.Body.Read(rspBody)
	if err != nil && err != io.EOF {
		T_Log.panic("resp.Body.Read.Error: %s \n", err)
		return
	}
	fmt.Printf("Response data: %s \n", string(rspBody))
}

func debugConfig(config T_Config) {
	for k, v := range config {
		fmt.Printf(" %s : %s \n ", k, v)
	}
}
