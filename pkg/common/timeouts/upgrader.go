package timeouts

import "context"

// MigrateToExplicit is used as a state upgrader function from implicit SDK's timeout object to explicit single-block based timeout
func MigrateToExplicit() func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	return func(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
		timeouts, ok := rawState["timeouts"]
		if !ok || timeouts == nil {
			return rawState, nil
		}

		rawState["timeouts"] = []any{timeouts}

		return rawState, nil
	}
}
