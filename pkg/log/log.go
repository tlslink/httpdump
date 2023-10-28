package log

import (
    "fmt"
    "log"
    "os"
    "path"
    "strings"
    "time"
)

const (
    LevelTrace = iota
    LevelDebug
    LevelInfo
    LevelWarn
    LevelError
    LevelFatal
)

var (
    lw       *logWriter
    logger   *log.Logger
    logLevel int
    levels   map[int]string

    dateFormat = "2006-01-02"
    logName    = "httpdump.log"
)

// 实现 os.Writer 接口
type logWriter struct {
    UseStdout bool
    FileName  string
    File      *os.File
    NowDate   string
}

// 实现日志文件的切割
func (lw *logWriter) Write(p []byte) (n int, err error) {
    if !lw.UseStdout {
        date := time.Now().Format(dateFormat)
        if lw.NowDate != date {
            _ = lw.File.Close()
            _ = os.Rename(lw.FileName, lw.FileName+"."+lw.NowDate)
            lw.NowDate = date
            lw.newFile()
        }
    }
    return lw.File.Write(p)
}

// 创建新文件
func (lw *logWriter) newFile() {
    if lw.UseStdout {
        lw.File = os.Stdout
        return
    }

    f, err := os.OpenFile(lw.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    lw.File = f
}

func Init(logPath, level string) {
    // 初始化 logger
    lw = &logWriter{
        UseStdout: logPath == "",
        FileName:  path.Join(logPath, logName),
        NowDate:   time.Now().Format(dateFormat),
    }

    lw.newFile()
    logLevel = logLevel2Int(level)
    logger = log.New(lw, "", log.LstdFlags|log.Lshortfile)
}

func LoggerWriter() *logWriter {
    return lw
}

// Logger 获取 log.Logger
func Logger() *log.Logger {
    return logger
}

func Level() int {
    return logLevel
}

func logLevel2Int(l string) int {
    levels = map[int]string{
        LevelTrace: "Trace",
        LevelDebug: "Debug",
        LevelInfo:  "Info",
        LevelWarn:  "Warn",
        LevelError: "Error",
        LevelFatal: "Fatal",
    }
    lvl := LevelInfo
    for k, v := range levels {
        if strings.ToLower(l) == strings.ToLower(v) {
            lvl = k
        }
    }
    return lvl
}

func output(l int, s ...interface{}) {
    logger.SetPrefix(fmt.Sprintf("[%s] ", levels[l]))
    _ = logger.Output(3, fmt.Sprintln(s...))
}

func Trace(v ...interface{}) {
    l := LevelTrace
    if logLevel > l {
        return
    }
    output(l, v...)
}

func Debug(v ...interface{}) {
    l := LevelDebug
    if logLevel > l {
        return
    }
    output(l, v...)
}

func Info(v ...interface{}) {
    l := LevelInfo
    if logLevel > l {
        return
    }
    output(l, v...)
}

func Warn(v ...interface{}) {
    l := LevelWarn
    if logLevel > l {
        return
    }
    output(l, v...)
}

func Error(v ...interface{}) {
    l := LevelError
    if logLevel > l {
        return
    }
    output(l, v...)
}

func Fatal(v ...interface{}) {
    l := LevelFatal
    if logLevel > l {
        return
    }
    output(l, v...)
    os.Exit(1)
}
