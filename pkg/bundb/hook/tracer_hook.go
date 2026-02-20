package hook

import (
	"context"

	"github.com/uptrace/bun"

	apm "go.elastic.co/apm/v2"
)

var _ bun.QueryHook = (*TracerHook)(nil)

type TracerHook struct{}

func NewTracerHook() *TracerHook {
	return &TracerHook{}
}

func (h *TracerHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	span, ctx := apm.StartSpan(ctx, "SQL "+event.Operation(), "db.query")
	ctx = apm.ContextWithSpan(ctx, span)

	return ctx
}

func (h *TracerHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	span := apm.SpanFromContext(ctx)
	if span != nil {
		if len(event.Query) > 100 {
			truncatedQuery := event.Query[:100] + "..."
			span.Context.SetLabel("query", truncatedQuery)
		} else {
			span.Context.SetLabel("query", event.Query)
		}

		span.End()
	}
}
