package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sirupsen/logrus"
)

// logRusTFLogHandler a tflog handler which integrates logrus so that logrus output gets handled natively
type logRusTFLogHandler struct {}

// EnableLogrusToTflog turns on passing of ruslog to tflog
func EnableLogrusToTflog() {
	logrus.AddHook(logRusTFLogHandler{})
}

// Fire off a logrus event
func (lh logRusTFLogHandler) Fire(e *logrus.Entry) error {
	go func(event *logrus.Entry) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

		if err := TFLogFire(ctx, event); err != nil {
			tflog.Error(context.Background(), "Log Write error", map[string]interface{}{"error": err.Error()})
		}

		cancel()
	}(e)

	return nil
}

// Levels that this logrus hook will handle
func (lh logRusTFLogHandler) Levels() []logrus.Level {
	return logrus.AllLevels
}

func TFLogFire(ctx context.Context, e *logrus.Entry) error {
	mes := e.Message
	addFields := map[string]interface{}{}

	switch e.Level {
	case logrus.DebugLevel:
		tflog.Debug(ctx, mes, addFields)
	case logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel:
		tflog.Error(ctx, mes, addFields)
	case logrus.InfoLevel:
		tflog.Info(ctx, mes, addFields)
	case logrus.WarnLevel:
		tflog.Warn(ctx, mes, addFields)
	}

	return nil
}
