// antlog
package antlog

import (
	_t "antlinker.com/tools"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var checkInterval float64 = 3  // 日志文件大小检查时间间隔（分钟）
var signalLogSize int64 = 1024 // 单个日志文件大小，单位M
var level logLevel = errLog    // 日志级别
var toConsole bool = false     // 是否输出到控制台
var logFileCount int = 3       // 日志文件数量--该设置未实现

const ( //日志级别，一共五个日志级别，高级别的日志会写入到低级别日志中
	debugLog logLevel = iota
	infoLog
	warningLog
	errLog
	fatalLog
	errLogLvel
)

var lastCheckTime time.Time                   // 记录上次日志文件大小检查时间
const logName string = "antlinker.log"        // 日志文件名称
const bizLogName string = "antlinker_biz.log" // 日志文件名称

type logLevel int         // 定义日志级别类型
var levelName = []string{ // 日志级别的名称
	debugLog:   "DEBUG",
	infoLog:    "INFO",
	warningLog: "WARN",
	errLog:     "ERROR",
	fatalLog:   "FATAL",
}

var _ = fmt.Println       // 避免fmt包未使用错误
var AntLogger *antLog     // 全局Logger
var sysLogger *log.Logger // 系统logger
var bizLogger *log.Logger // 业务logger
var sysLogFile *os.File   // logfile句柄
var bizLogFile *os.File   // bizLogfile句柄

type antLog struct {
	level         logLevel // 日志级别
	filePath      string   // 日志文件路径
	fileCount     int      // 保留日志文件数量----这个设置暂时无效
	fileSize      int64    // 单个日志文件最大大小
	alsoToConsole bool     // 是否同时输出到控制台
}

// 初始化antlog日志模块配置
func init() {
	log.Printf("[启动] 蚂蚁日志模块正在初始化....\r\n")
	readConfig()
	AntLogger = new(antLog)
	AntLogger.level = level
	AntLogger.filePath = getCurrentPath() + "/log/"
	AntLogger.fileCount = logFileCount
	AntLogger.fileSize = signalLogSize * 1024
	AntLogger.alsoToConsole = toConsole
	// 初始化日志文件
	initSysLog()
	initBizLog()
	log.Printf("[启动] 日志级别...........%s\r\n", levelName[AntLogger.level])
	log.Printf("[启动] 日志路径...........%s\r\n", AntLogger.filePath)
	log.Printf("[启动] 最大日志...........%d\r\n", AntLogger.fileSize)
	log.Printf("[启动] 轮训间隔...........%f分钟\r\n", checkInterval)
	log.Printf("[启动] 控制台输出..........%t\r\n", AntLogger.alsoToConsole)
	log.Printf("[启动] 蚂蚁日志模块初始化完成\r\n")
}

// 读取配置文件信息
func readConfig() {
	confPath := _t.GetCurrentPath() + "/config/antlog.conf"
	log.Printf("[启动] 读取日志配置文件....%s\r\n", confPath)
	p, err := _t.ReadConfig(confPath)
	if err != nil {
		log.Printf("[错误] 读取配置文件错误..........%s\r\n", err.Error())
		os.Exit(1)
	}
	if ci, err := strconv.ParseFloat(p["LOG_CHECK_INTERVAL"], 64); err == nil {
		checkInterval = ci
	} else {
		log.Printf("[错误] 日志配置文件错误[LOG_CHECK_INTERVAL]=%s:%s\r\n", p["LOG_CHECK_INTERVAL"], err.Error())
		os.Exit(1)
	}

	if ss, err := strconv.ParseInt(p["LOG_MAX_FILE_SIZE"], 0, 64); err == nil {
		signalLogSize = ss
	} else {
		log.Printf("[错误] 日志配置文件错误[LOG_MAX_FILE_SIZE]=%s:%s\r\n", p["LOG_MAX_FILE_SIZE"], err.Error())
		os.Exit(1)
	}

	logLevelMap := make(map[string]logLevel)
	logLevelMap["DEBUG"] = debugLog
	logLevelMap["INFO"] = infoLog
	logLevelMap["WARN"] = warningLog
	logLevelMap["ERROR"] = errLog
	logLevelMap["FATAL"] = fatalLog

	if strings.TrimSpace(p["LOG_LEVEL"]) == "" {
		log.Printf("[错误] 日志配置文件错误[LOG_LEVEL]=%s:%s\r\n", p["LOG_LEVEL"], err.Error())
		os.Exit(1)
	} else {
		level = logLevelMap[p["LOG_LEVEL"]]
	}

	if strings.ToLower(p["LOG_OUT_PUT_TO_CONSOLE"]) == "true" {
		toConsole = true
	} else {
		toConsole = false
	}

	if lc, err := strconv.Atoi(p["LOG_MAX_FILE_COUNT"]); err == nil {
		logFileCount = lc
	} else {
		log.Printf("[错误] 日志配置文件错误[LOG_MAX_FILE_COUNT]=%s:%s\r\n", p["LOG_MAX_FILE_COUNT"], err.Error())
		os.Exit(1)
	}
}

// 释放日志打开的文件句柄
func ReleaseAntLog() {
	sysLogFile.Close()
	bizLogFile.Close()
}

// 初始化系统日志
func initSysLog() {
	sysLog, err := _t.CreateFile(AntLogger.filePath, logName)
	if err != nil {
		panic(err)
	}
	sysLogFile, err = os.OpenFile(sysLog, os.O_RDWR|os.O_APPEND, 0666)
	sysLogger = log.New(sysLogFile, "", log.Ldate|log.Ltime)
}

// 初始化业务日志
func initBizLog() {
	bizLog, err := _t.CreateFile(AntLogger.filePath, bizLogName)
	if err != nil {
		panic(err)
	}
	bizLogFile, err = os.OpenFile(bizLog, os.O_RDWR|os.O_APPEND, 0666)
	bizLogger = log.New(bizLogFile, "", log.Ldate|log.Ltime)
}

// 获取当前执行程序的绝对路径，不包含文件名
func getCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	return path.Dir(file)
}

// 检查日志文件大小，如果超出规定大小进行重命名处理
func checkLogFileSize() {
	if time.Since(lastCheckTime).Minutes() > checkInterval {
		// 检查系统日志大小是否超出上限
		sysSize, sysErr := _t.GetFileSize(AntLogger.filePath + logName)
		if sysErr != nil {
			panic(sysErr)
		}
		// 如果超出上限，关闭文件备份后重新创建文件
		if sysSize > AntLogger.fileSize {
			tmp := fmt.Sprintf("_%d.", _t.GetTimeStamp())
			sysLogFile.Close()
			_t.Rename(AntLogger.filePath+logName, AntLogger.filePath+strings.Replace(logName, ".", tmp, -1))
			initSysLog()
		}
		// 检查业务日志大小是否超出上限
		bizSize, bizErr := _t.GetFileSize(AntLogger.filePath + bizLogName)
		if bizErr != nil {
			panic(bizErr)
		}
		// 如果超出上限，关闭文件备份后重新创建文件
		if bizSize > AntLogger.fileSize {
			tmp := fmt.Sprintf("_%d.", _t.GetTimeStamp())
			bizLogFile.Close()
			_t.Rename(AntLogger.filePath+bizLogName, AntLogger.filePath+strings.Replace(bizLogName, ".", tmp, -1))
			initBizLog()
		}
		// 更新日志最后检查时间
		lastCheckTime = time.Now()
	}
}

// 把系统日志内容写入文件
func writeSysLog(l logLevel, content string) {
	if l >= AntLogger.level { // 超过规定日志级别的都写入文件
		sysLogger.Printf(content)
	}
	if AntLogger.alsoToConsole {
		log.Println(content)
	}
	checkLogFileSize()
}

// 把业务日志内容写入文件
func writeBizLog(content string) {
	bizLogger.Printf(content)
	if AntLogger.alsoToConsole {
		log.Println(content)
	}
	checkLogFileSize()
}

// 写业务日志的方法
func (l *antLog) BizLog(format string, v ...interface{}) {
	writeBizLog(fmt.Sprintf(format, v...))
}

// 写系统日志的方法
func (l *antLog) Debug(format string, v ...interface{}) {
	writeSysLog(debugLog, fmt.Sprintf("[%s] %s", levelName[debugLog], fmt.Sprintf(format, v...)))
}
func (l *antLog) Info(format string, v ...interface{}) {
	writeSysLog(infoLog, fmt.Sprintf("[%s] %s", levelName[infoLog], fmt.Sprintf(format, v...)))
}
func (l *antLog) Warn(format string, v ...interface{}) {
	writeSysLog(warningLog, fmt.Sprintf("[%s] %s", levelName[warningLog], fmt.Sprintf(format, v...)))
}
func (l *antLog) Error(format string, v ...interface{}) {
	writeSysLog(errLog, fmt.Sprintf("[%s] %s", levelName[errLog], fmt.Sprintf(format, v...)))
}
func (l *antLog) Fatal(format string, v ...interface{}) {
	writeSysLog(fatalLog, fmt.Sprintf("[%s] %s", levelName[fatalLog], fmt.Sprintf(format, v...)))
}
