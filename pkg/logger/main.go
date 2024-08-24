package logger

import (
	"net"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Args map[string]any

type Logger interface {
	Error(err error, keyValMap Args)
	Info(msg string, keyValMap Args)
	Warn(msg string, keyValMap Args)
}

type logger struct {
	log *zap.SugaredLogger
}

func New(f *os.File, conn net.Conn, isDebug bool) (Logger, error) {
	var cores []zapcore.Core

	var cfg zap.Config
	var level zapcore.LevelEnabler

	if isDebug {
		cfg = zap.NewDevelopmentConfig()
		level = zapcore.DebugLevel
	} else {
		cfg = zap.NewProductionConfig()
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //serialize level by caps string
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   //serialize time by ISO8601

	encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)

	var syncer zapcore.WriteSyncer
	if f != nil && conn != nil {
		fileSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename: f.Name(),
			MaxSize:  20,
			MaxAge:   5,
		})
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(conn), fileSyncer)
	}

	if f != nil {
		cores = append(cores, zapcore.NewCore(encoder, syncer, level))
	}

	if conn != nil {
		cores = append(cores, zapcore.NewCore(encoder, syncer, level))
	}

	core := zapcore.NewTee(cores...)

	log := zap.New(core)

	defer log.Sync()
	return &logger{log: log.Sugar()}, nil
}

func (logger *logger) Info(msg string, keyValueMap Args) {
	args := logger.parseArguments(keyValueMap)
	logger.log.Infof(msg, args...)
}

func (logger *logger) Error(err error, keyValueMap Args) {
	args := logger.parseArguments(keyValueMap)
	logger.log.Errorf(err.Error(), args...)
}

func (logger *logger) Warn(msg string, keyValueMap Args) {
	args := logger.parseArguments(keyValueMap)
	logger.log.Warnf(msg, args...)
}

func (logger *logger) parseArguments(args map[string]any) []any {
	var argsSlice []any
	for key, value := range args {
		argsSlice = append(argsSlice, key, value)
	}
	return argsSlice
}
