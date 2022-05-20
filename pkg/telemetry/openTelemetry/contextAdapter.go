package openTelemetry

import (
	"context"
	"fmt"
)

type ContextAdapter struct {
	Context context.Context
}

func (a *ContextAdapter) Get(key string) string {
	return fmt.Sprint(a.Context.Value(key))
}

func (a *ContextAdapter) Set(key string, value string) {
	a.Context = context.WithValue(a.Context, key, value)
}

func (a *ContextAdapter) Keys() []string {
	return nil
}
