package logger

import (
  "go.uber.org/zap"
  "encoding/json"
)

var log *zap.Logger

func SetupLogging(loglevel string) {
  rawJSON := []byte(`{
  "level": "` + loglevel + `",
  "development": true,
  "encoding": "json",
  "outputPaths": ["stdout"],
  "errorOutputPaths": ["stderr"],
  "initialFields": {},
  "encoderConfig": {
  "messageKey": "message",
  "levelKey": "level",
  "timeKey": "ts",
  "timeEncoder": "iso8601",
  "levelEncoder": "lowercase"
  }
  }`)

  var cfg zap.Config
  if err := json.Unmarshal(rawJSON, &cfg); err != nil {
    panic(err)
  }
  log = zap.Must(cfg.Build())

}

func Info(message string, fields ...zap.Field) {
  log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
  log.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
  log.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
  log.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
  log.Fatal(message, fields...)
}
