package log

import (
	"os"
	"path/filepath"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Option interface {
	apply(*logrus.Logger) error
}

type OptionFunc func(*logrus.Logger) error

func (f OptionFunc) apply(l *logrus.Logger) error {
	return f(l)
}

//StandardOption
func StandardOption(env string) Option {
	return OptionFunc(func(l *logrus.Logger) error {
		switch env {
		case "production":
			// 生产环境中仅在 log 文件中输出
			l.SetFormatter(&NullFormatter{})
			devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			if err != nil {
				return err
			}
			l.SetOutput(devNull)
		default:
			l.SetOutput(os.Stdout)
			l.SetReportCaller(true)
			l.SetFormatter(&nested.Formatter{
				TimestampFormat: DefaultTimestampFormat,
				HideKeys:        true,
			})
		}
		if err := WriteToFileOption("./logs").apply(l); err != nil {
			return err
		}
		return nil
	})
}

//WriteToFileOption
func WriteToFileOption(dir string) Option {
	return OptionFunc(func(l *logrus.Logger) error {
		writer, err := rotatelogs.New(
			filepath.Join(dir, "log.%Y%m%d.log"),
			rotatelogs.WithLinkName(filepath.Join(dir, "log.log")),
			rotatelogs.WithMaxAge(7*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			return err
		}
		errWriter, err := rotatelogs.New(
			filepath.Join(dir, "err.%Y%m%d.log"),
			rotatelogs.WithLinkName(filepath.Join(dir, "err.log")),
			rotatelogs.WithMaxAge(7*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			return err
		}

		lfHook := lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: errWriter,
			logrus.FatalLevel: errWriter,
			logrus.PanicLevel: errWriter,
		}, &logrus.JSONFormatter{})
		l.AddHook(lfHook)
		return nil
	})
}
