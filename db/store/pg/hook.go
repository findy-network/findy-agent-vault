package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lib/pq"
	"github.com/qustavo/sqlhooks/v2"
)

type traceCtxKey string

type traceHook struct{}

func initTraceHook(dataSourceName string) (*sql.DB, error) {
	const driverName = "pgWithHooks"
	sql.Register(driverName, sqlhooks.Wrap(&pq.Driver{}, &traceHook{}))
	return sql.Open(driverName, dataSourceName)
}

func (h *traceHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	utils.LogTrace().Infof("> %s %q", query, args)
	return context.WithValue(ctx, traceCtxKey("QueryTrace"), time.Now()), nil
}

func (h *traceHook) After(ctx context.Context, _ string, _ ...interface{}) (context.Context, error) {
	if begin, ok := ctx.Value(traceCtxKey("QueryTrace")).(time.Time); ok {
		utils.LogTrace().Infof("< took: %s\n", time.Since(begin))
	}
	return ctx, nil
}
