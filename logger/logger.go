package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger logrus.Logger

const loggerDateTimeFormat = "02.01.2006 15:04:05"

// NewLogger returns logger instance
func NewLogger(level string) (*Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(lvl)
	customFormatter := &logrus.TextFormatter{}
	customFormatter.TimestampFormat = loggerDateTimeFormat
	customFormatter.FullTimestamp = true
	logger.SetFormatter(customFormatter)

	return (*Logger)(logger), nil
}

// Start start log for endpoint start
func (log *Logger) Start(service, action, clientID, managerID string) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, clientID, managerID)).
		Infof("Start %s in %s", action, service)
}

// End end log for endpoint end
func (log *Logger) End(service, action, message, clientID, managerID string, data interface{}) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, clientID, managerID)).
		Infof("End %s in %s, message: %s, data: %v", action, service, message, data)
}

// Info custom info log
func (log *Logger) Info(service, action, message, clientID, managerID string, data interface{}) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, clientID, managerID)).
		Infof("message: %s, data: %v", message, data)
}

// Trace custom trace log
func (log *Logger) Trace(service, action, message string) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, "", "")).
		Trace(message)
}

// Debug custom debug log
func (log *Logger) Debug(service, action, message, clientID, managerID string, data interface{}) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, clientID, managerID)).
		Debugf(" message: %s, data: %v", message, data)
}

// Warn custom warn log
func (log *Logger) Warn(service, action, message, clientID, managerID string, code int, data interface{}) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, clientID, managerID)).
		Warnf("code: %d, message: %s, data: %v", code, message, data)
}

// Error custom error message log
func (log *Logger) Error(service, action, message, clientID, managerID string, code int, error interface{}) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, clientID, managerID)).
		Errorf("code: %d, message: %s, error: %v", code, message, error)
}

// Fatal custom fatal error log
func (log *Logger) Fatal(service, action, message, clientID, managerID string, error interface{}) {
	(*logrus.Logger)(log).
		WithFields(getFields(service, action, clientID, managerID)).
		Fatalf("message: %s, error: %v", message, error)
}

func getFields(service, action, clientID, managerID string) map[string]interface{} {
	fields := make(map[string]interface{})
	fields["service"] = service
	fields["action"] = action
	fields["clientID"] = clientID
	fields["managerID"] = managerID
	return fields
}
