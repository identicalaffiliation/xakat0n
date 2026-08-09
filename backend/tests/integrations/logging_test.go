package integrations

import (
	"log/slog"

	"github.com/identicalaffiliation/xakat0n/backend/internal/shared/logctx"
)

type testLogger struct {
	*slog.Logger
	*logctx.LogCtx
}

func newTestLogger() *testLogger {
	return &testLogger{
		Logger: slog.Default(),
		LogCtx: logctx.NewLogCtx(),
	}
}
