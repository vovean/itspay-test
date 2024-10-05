package testkit

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func NewOTELZapTestLogger(t zaptest.TestingT) *otelzap.Logger {
	return otelzap.New(zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller())))
}
