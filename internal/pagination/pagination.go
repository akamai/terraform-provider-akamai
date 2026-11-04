// Package pagination provides helpers for traversing paginated APIs.
package pagination

import (
	"context"
	"errors"
	"fmt"
)

// ErrNilPage indicates that a page request returned no response and no error.
var ErrNilPage = errors.New("received nil page without error")

// PageOptions configures pagination whose request state, page response, and items are endpoint-specific.
type PageOptions[State, Item, Page any] struct {
	// InitialState is the request state used to fetch the first page.
	InitialState State

	// Fetch retrieves and validates a complete page response using the current request state.
	Fetch func(context.Context, State) (*Page, error)

	// Items extracts the items from a fetched page.
	Items func(*Page) []Item

	// Next returns the next request state and whether another page should be fetched.
	Next func(State, *Page) (State, bool)
}

// CollectPageItems retrieves and combines page items until Next indicates that pagination is complete.
func CollectPageItems[State, Item, Page any](ctx context.Context, opts PageOptions[State, Item, Page]) ([]Item, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if opts.Fetch == nil {
		return nil, fmt.Errorf("page fetch function cannot be nil")
	}
	if opts.Items == nil {
		return nil, fmt.Errorf("page items function cannot be nil")
	}
	if opts.Next == nil {
		return nil, fmt.Errorf("next page function cannot be nil")
	}

	state := opts.InitialState
	var items []Item
	for {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("context terminated before page retrieval: %w", err)
		}

		page, err := opts.Fetch(ctx, state)
		if err != nil {
			return nil, err
		}
		if page == nil {
			return nil, ErrNilPage
		}
		items = append(items, opts.Items(page)...)

		nextState, more := opts.Next(state, page)
		if !more {
			return items, nil
		}
		state = nextState
	}
}

// Options configures collection of keyset-paginated API results.
type Options[Item, Page any] struct {
	// Fetch returns the complete page for the supplied cursor and page size.
	Fetch func(context.Context, string, int) (*Page, error)

	// Items extracts the items from a fetched page.
	Items func(*Page) []Item

	// NextCursor extracts the next cursor from a fetched page. It returns nil only for the final page.
	NextCursor func(*Page) *string

	// PageSize is the maximum number of items requested per page. It must be positive.
	PageSize int

	// Limit optionally caps the total number of collected items. The final request uses the remaining limit.
	Limit *int64
}

// Collect retrieves and combines all pages. A nil next cursor completes collection. Some endpoints can
// return an empty page with a valid next cursor, so it is followed normally. A repeated cursor always returns an
// error to prevent incomplete or duplicate results. When Limit is set, requests use the remaining page size and
// excess response items are discarded.
func Collect[Item, Page any](ctx context.Context, opts Options[Item, Page]) ([]Item, error) {
	var collected []Item
	err := iterate(ctx, opts, func(item Item) bool {
		collected = append(collected, item)
		return false
	})
	if err != nil {
		return nil, err
	}
	return collected, nil
}

// Find returns the first item matching "matches" and whether an item was found. A nil next cursor completes
// the search. Some endpoints can return an empty page with a valid next cursor, so it is followed normally. A repeated
// cursor always returns an error to prevent incomplete or duplicate results.
func Find[Item, Page any](ctx context.Context, opts Options[Item, Page], matches func(Item) bool) (Item, bool, error) {
	var found Item
	if matches == nil {
		return found, false, fmt.Errorf("item match function cannot be nil")
	}

	foundMatch := false
	err := iterate(ctx, opts, func(item Item) bool {
		if !matches(item) {
			return false
		}
		found = item
		foundMatch = true
		return true
	})
	if err != nil {
		return found, false, err
	}
	return found, foundMatch, nil
}

func iterate[Item, Page any](ctx context.Context, opts Options[Item, Page], visit func(Item) bool) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if opts.Fetch == nil {
		return fmt.Errorf("page fetch function cannot be nil")
	}
	if opts.Items == nil {
		return fmt.Errorf("page items function cannot be nil")
	}
	if opts.NextCursor == nil {
		return fmt.Errorf("next cursor function cannot be nil")
	}
	if opts.PageSize <= 0 {
		return fmt.Errorf("page size must be positive")
	}
	var remaining int64
	if opts.Limit != nil {
		remaining = *opts.Limit
		if remaining < 0 {
			return fmt.Errorf("limit cannot be negative")
		}
		if remaining == 0 {
			return nil
		}
	}

	seen := make(map[string]struct{})

	cursor := ""
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context terminated before page retrieval: %w", err)
		}

		pageSize := opts.PageSize
		if opts.Limit != nil {
			pageSize = int(min(remaining, int64(pageSize)))
		}
		page, err := opts.Fetch(ctx, cursor, pageSize)
		if err != nil {
			return err
		}
		if page == nil {
			return ErrNilPage
		}

		items := opts.Items(page)
		if opts.Limit != nil {
			items = items[:min(int64(len(items)), remaining)]
		}
		for _, item := range items {
			if visit(item) {
				return nil
			}
		}
		if opts.Limit != nil {
			remaining -= int64(len(items))
			if remaining == 0 {
				return nil
			}
		}

		next := opts.NextCursor(page)
		if next == nil {
			return nil
		}
		if *next == "" {
			return fmt.Errorf("pagination returned an empty cursor")
		}
		if _, exists := seen[*next]; exists {
			return fmt.Errorf("pagination returned a non-advancing cursor %q", *next)
		}
		seen[*next] = struct{}{}
		cursor = *next
	}
}
