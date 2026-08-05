package tx

import "context"

type txKey struct{}

func withTx(ctx context.Context, tx DBTX) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// DBTXFromContext возвращает DBTX, открытую текущим Manager.WithTx, если вызов идёт
// внутри неё; иначе — fallback (обычно *pgxpool.Pool, переданный репозиторию при
// конструировании).
func DBTXFromContext(ctx context.Context, fallback DBTX) DBTX {
	if tx, ok := ctx.Value(txKey{}).(DBTX); ok {
		return tx
	}

	return fallback
}
