package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type numberedPage struct {
	number   int
	metadata string
	items    []int
}

type integerPage struct {
	items []int
	next  *string
}

func TestCollectPageItems(t *testing.T) {
	t.Parallel()

	t.Run("collects complete pages using endpoint state", func(t *testing.T) {
		t.Parallel()

		calls := 0
		items, err := CollectPageItems(context.Background(), PageOptions[int, int, numberedPage]{
			InitialState: 1,
			Fetch: func(_ context.Context, page int) (*numberedPage, error) {
				calls++
				if calls > 3 {
					return nil, errors.New("pagination exceeded the expected three fetches")
				}
				return &numberedPage{
					number:   page,
					metadata: "metadata from page",
					items:    []int{page},
				}, nil
			},
			Items: func(page *numberedPage) []int { return page.items },
			Next: func(page int, _ *numberedPage) (int, bool) {
				return page + 1, page < 3
			},
		})

		require.NoError(t, err)
		require.Equal(t, []int{1, 2, 3}, items)
	})

	tests := map[string]struct {
		ctx         context.Context
		opts        PageOptions[int, int, numberedPage]
		wantErr     string
		wantErrorIs error
	}{
		"rejects a nil context": {
			opts: PageOptions[int, int, numberedPage]{
				Fetch: func(context.Context, int) (*numberedPage, error) {
					return nil, errors.New("fetch should not be called")
				},
				Items: func(*numberedPage) []int { return nil },
				Next: func(int, *numberedPage) (int, bool) {
					return 0, false
				},
			},
			wantErr: "context cannot be nil",
		},
		"rejects a missing fetch function": {
			ctx: context.Background(),
			opts: PageOptions[int, int, numberedPage]{
				Items: func(*numberedPage) []int { return nil },
				Next: func(int, *numberedPage) (int, bool) {
					return 0, false
				},
			},
			wantErr: "page fetch function cannot be nil",
		},
		"rejects a missing next function": {
			ctx: context.Background(),
			opts: PageOptions[int, int, numberedPage]{
				Fetch: func(context.Context, int) (*numberedPage, error) {
					return &numberedPage{}, nil
				},
				Items: func(*numberedPage) []int { return nil },
			},
			wantErr: "next page function cannot be nil",
		},
		"rejects a missing items function": {
			ctx: context.Background(),
			opts: PageOptions[int, int, numberedPage]{
				Fetch: func(context.Context, int) (*numberedPage, error) {
					return &numberedPage{}, nil
				},
				Next: func(int, *numberedPage) (int, bool) {
					return 0, false
				},
			},
			wantErr: "page items function cannot be nil",
		},
		"returns a fetch error": {
			ctx: context.Background(),
			opts: PageOptions[int, int, numberedPage]{
				Fetch: func(context.Context, int) (*numberedPage, error) {
					return nil, errors.New("fetch failed")
				},
				Items: func(*numberedPage) []int { return nil },
				Next: func(int, *numberedPage) (int, bool) {
					return 0, false
				},
			},
			wantErr: "fetch failed",
		},
		"rejects a nil page": {
			ctx: context.Background(),
			opts: PageOptions[int, int, numberedPage]{
				Fetch: func(context.Context, int) (*numberedPage, error) {
					return nil, nil
				},
				Items: func(*numberedPage) []int { return nil },
				Next: func(int, *numberedPage) (int, bool) {
					return 0, false
				},
			},
			wantErr:     ErrNilPage.Error(),
			wantErrorIs: ErrNilPage,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := CollectPageItems(tc.ctx, tc.opts)
			require.EqualError(t, err, tc.wantErr)
			if tc.wantErrorIs != nil {
				require.ErrorIs(t, err, tc.wantErrorIs)
			}
		})
	}
}

func TestCollect(t *testing.T) {
	t.Parallel()

	cursor := func(value string) *string { return &value }
	zeroLimit := int64(0)
	limit := int64(3)
	negativeLimit := int64(-1)
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	betweenPagesCtx, cancelBetweenPages := context.WithCancel(context.Background())

	tests := map[string]struct {
		context     func() context.Context
		opts        Options[int, integerPage]
		want        []int
		wantErr     string
		wantErrorIs error
	}{
		"rejects a nil context": {
			context: func() context.Context { return nil },
			opts: Options[int, integerPage]{
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return nil, errors.New("fetch should not be called")
				},
			},
			wantErr: "context cannot be nil",
		},
		"rejects a missing fetch function": {
			opts:    Options[int, integerPage]{PageSize: 1},
			wantErr: "page fetch function cannot be nil",
		},
		"rejects a non-positive page size": {
			opts: Options[int, integerPage]{
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return nil, errors.New("fetch should not be called")
				},
			},
			wantErr: "page size must be positive",
		},
		"rejects a negative limit": {
			opts: Options[int, integerPage]{
				PageSize: 1,
				Limit:    &negativeLimit,
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return nil, errors.New("fetch should not be called")
				},
			},
			wantErr: "limit cannot be negative",
		},
		"zero limit does not fetch a page": {
			opts: Options[int, integerPage]{
				PageSize: 1,
				Limit:    &zeroLimit,
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return nil, errors.New("fetch should not be called")
				},
			},
			want: []int(nil),
		},
		"returns a fetch error": {
			opts: Options[int, integerPage]{
				PageSize: 1,
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return nil, errors.New("fetch failed")
				},
			},
			wantErr: "fetch failed",
		},
		"rejects a nil page": {
			opts: Options[int, integerPage]{
				PageSize: 1,
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return nil, nil
				},
			},
			wantErr: ErrNilPage.Error(),
		},
		"stops before fetching when cancelled": {
			context: func() context.Context { return cancelledCtx },
			opts: Options[int, integerPage]{
				PageSize: 1,
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return nil, errors.New("fetch should not be called")
				},
			},
			wantErr:     "context terminated before page retrieval: context canceled",
			wantErrorIs: context.Canceled,
		},
		"stops between pages when cancelled": {
			context: func() context.Context { return betweenPagesCtx },
			opts: Options[int, integerPage]{
				PageSize: 1,
				Fetch: func(_ context.Context, nextCursor string, _ int) (*integerPage, error) {
					if nextCursor != "" {
						return nil, errors.New("second page should not be fetched")
					}
					cancelBetweenPages()
					return &integerPage{next: cursor("second")}, nil
				},
			},
			wantErr:     "context terminated before page retrieval: context canceled",
			wantErrorIs: context.Canceled,
		},
		"collects each page": {
			opts: Options[int, integerPage]{
				PageSize: 2,
				Fetch: func(_ context.Context, nextCursor string, _ int) (*integerPage, error) {
					switch nextCursor {
					case "":
						return &integerPage{items: []int{1, 2}, next: cursor("second")}, nil
					case "second":
						return &integerPage{items: []int{3}}, nil
					default:
						return nil, errors.New("unexpected cursor")
					}
				},
			},
			want: []int{1, 2, 3},
		},
		"follows a cursor after an empty page": {
			opts: Options[int, integerPage]{
				PageSize: 1,
				Fetch: func(_ context.Context, nextCursor string, _ int) (*integerPage, error) {
					if nextCursor == "" {
						return &integerPage{next: cursor("second")}, nil
					}
					return &integerPage{items: []int{3}}, nil
				},
			},
			want: []int{3},
		},
		"rejects an empty cursor": {
			opts: Options[int, integerPage]{
				PageSize: 1,
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return &integerPage{next: cursor("")}, nil
				},
			},
			wantErr: "pagination returned an empty cursor",
		},
		"limits requests and collected items": {
			opts: Options[int, integerPage]{
				Limit:    &limit,
				PageSize: 2,
				Fetch: func(_ context.Context, nextCursor string, pageSize int) (*integerPage, error) {
					if nextCursor == "" {
						if pageSize != 2 {
							return nil, errors.New("first page size must be 2")
						}
						return &integerPage{items: []int{1, 2}, next: cursor("second")}, nil
					}
					if pageSize != 1 {
						return nil, errors.New("second page size must be 1")
					}
					return &integerPage{items: []int{3, 4}}, nil
				},
			},
			want: []int{1, 2, 3},
		},
		"rejects a repeated cursor": {
			opts: Options[int, integerPage]{
				PageSize: 1,
				Fetch: func(_ context.Context, _ string, _ int) (*integerPage, error) {
					return &integerPage{next: cursor("repeat")}, nil
				},
			},
			wantErr: `pagination returned a non-advancing cursor "repeat"`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.opts.Items = func(page *integerPage) []int { return page.items }
			tc.opts.NextCursor = func(page *integerPage) *string { return page.next }

			ctx := context.Background()
			if tc.context != nil {
				ctx = tc.context()
			}

			fetch := tc.opts.Fetch
			calls := 0
			if fetch != nil {
				tc.opts.Fetch = func(ctx context.Context, cursor string, pageSize int) (*integerPage, error) {
					calls++
					if calls > 10 {
						return nil, errors.New("pagination exceeded the test fetch limit")
					}
					return fetch(ctx, cursor, pageSize)
				}
			}

			got, err := Collect(ctx, tc.opts)
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				if tc.wantErrorIs != nil {
					require.ErrorIs(t, err, tc.wantErrorIs)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	t.Run("rejects a multi-cursor cycle", func(t *testing.T) {
		t.Parallel()

		cursors := []string{"", "a", "b", "a"}
		calls := 0
		_, err := Collect(context.Background(), Options[int, integerPage]{
			PageSize: 1,
			Fetch: func(_ context.Context, nextCursor string, _ int) (*integerPage, error) {
				if calls >= len(cursors)-1 {
					return nil, errors.New("cursor cycle was not rejected")
				}
				if nextCursor != cursors[calls] {
					return nil, errors.New("unexpected cursor")
				}
				calls++
				return &integerPage{next: &cursors[calls]}, nil
			},
			Items:      func(page *integerPage) []int { return page.items },
			NextCursor: func(page *integerPage) *string { return page.next },
		})
		require.EqualError(t, err, `pagination returned a non-advancing cursor "a"`)
		require.Equal(t, 3, calls)
	})
}

func TestFind(t *testing.T) {
	t.Parallel()

	cursor := func(value string) *string { return &value }
	got, found, err := Find(context.Background(), Options[int, integerPage]{
		PageSize: 1,
		Fetch: func(_ context.Context, nextCursor string, _ int) (*integerPage, error) {
			switch nextCursor {
			case "":
				return &integerPage{items: []int{1}, next: cursor("second")}, nil
			case "second":
				return &integerPage{items: []int{2}, next: cursor("third")}, nil
			default:
				return nil, errors.New("fetch should stop after the matching item")
			}
		},
		Items:      func(page *integerPage) []int { return page.items },
		NextCursor: func(page *integerPage) *string { return page.next },
	}, func(item int) bool {
		return item == 2
	})

	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, 2, got)
}
