package log_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudokyo/log"
)

func TestLog(t *testing.T) {
	format := "Hello %s"
	arg := "world"

	t.Run("Test plain text log", func(t *testing.T) {
		log.Output(log.OutputConsole)

		log.Trace("Hello", arg).Send()
		log.Tracef(format, arg)

		log.Debug("Hello", arg).Send()
		log.Debugf(format, arg)

		log.Info("Hello", arg).Send()
		log.Infof(format, arg)

		log.Warn("Hello", arg).Send()
		log.Warnf(format, arg)

		log.Error("Hello", arg).Send()
		log.Errorf(format, arg)
	})

	t.Run("Test json log", func(t *testing.T) {
		log.Output(log.OutputDefault)

		log.Trace("Hello", arg).Send()
		log.Tracef(format, arg)

		log.Debug("Hello", arg).Send()
		log.Debugf(format, arg)

		log.Info("Hello", arg).Send()
		log.Infof(format, arg)

		log.Warn("Hello", arg).Send()
		log.Warnf(format, arg)

		log.Error("Hello", arg).Send()
		log.Errorf(format, arg)
	})
}

func TestLogContext(t *testing.T) {
	ctx := context.Background()
	ctx = log.WithValue(ctx, log.RequestKey, time.Now().String())

	format := "Hello %s"
	arg := "world"
	err := errors.New("Test e2e")

	log.Trace("Hello", arg, ctx).Send()
	log.Tracef(format, arg, ctx)

	log.Debug("Hello", arg, ctx).Send()
	log.Debugf(format, arg, ctx)

	log.Info("Hello", arg, ctx).Send()
	log.Infof(format, arg, ctx)

	log.Warn("Hello", arg, ctx).Send()
	log.Warnf(format, arg, ctx)

	log.Error("Hello", arg, ctx, true, err).Send()
	log.Errorf(format, arg, ctx, err)
}

func TestLogStackTrace(t *testing.T) {
	ctx := context.Background()
	ctx = log.WithValue(ctx, log.RequestKey, time.Now().String())

	arg := "world"
	err := errors.New("Test e2e")

	log.Error("Hello", arg, ctx, true, err).Send()
}
