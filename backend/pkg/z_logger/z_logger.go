package z_logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	log *zerolog.Logger
}

func NewLogger(dir, fileName string, isDev bool) *ZeroLogger {
	var writer io.Writer

	writer = &lumberjack.Logger{
		Filename:   filepath.Join(dir, "app.log"),
		MaxSize:    5,
		MaxBackups: 1,
		MaxAge:     2,
		Compress:   true,
	}

	if isDev {
		writer = io.MultiWriter(
			zerolog.ConsoleWriter{Out: os.Stdout},
			writer,
		)
	}

	log := zerolog.New(writer).
		With().
		Timestamp().
		Logger()

	return &ZeroLogger{
		log: &log,
	}
}

func (l *ZeroLogger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *ZeroLogger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

func (l *ZeroLogger) Error(msg string, err error) {
	l.log.Error().Err(err).Msg(msg)
}

func (l *ZeroLogger) Fatal(msg string, err error) {
	l.log.Fatal().Err(err).Msg(msg)
}
