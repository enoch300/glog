/*
* @Author: wangqilong
* @Description:
* @File: glog
* @Date: 2021/9/3 11:05 上午
 */

package glog

import (
	"bufio"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type MineFormatter struct{}

func (s *MineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Local().Format("2006-01-02 15:04:05"), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func write(baseLogPath string, level string, maxAge time.Duration, rotationTime time.Duration) *rotatelogs.RotateLogs {
	logier, err := rotatelogs.New(
		baseLogPath+"_"+level+"_%Y-%m-%d.log",
		rotatelogs.WithLinkName(baseLogPath+"_"+level), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),                  // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime),      // 日志切割时间间隔
	)

	if err != nil {
		log.Fatalf("config local file system logger error. %s", err.Error())
	}
	return logier
}

// NewLogger logPath 日志目录, logFileName 日志文件名, maxAge 文件最大保存时间, rotationTime 日志切割时间间隔
func NewLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *logrus.Logger {
	fullLogPath := path.Join(logPath, logFileName)
	src, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	output := bufio.NewWriter(src)

	l := logrus.New()
	l.SetOutput(output)

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: write(fullLogPath, "debug", maxAge, rotationTime), // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  write(fullLogPath, "info", maxAge, rotationTime),
		logrus.WarnLevel:  write(fullLogPath, "warn", maxAge, rotationTime),
		logrus.ErrorLevel: write(fullLogPath, "error", maxAge, rotationTime),
		logrus.FatalLevel: write(fullLogPath, "fatal", maxAge, rotationTime),
		logrus.PanicLevel: write(fullLogPath, "panic", maxAge, rotationTime),
	}, &MineFormatter{})

	l.AddHook(lfHook)
	return l
}
