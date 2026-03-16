package client

import (
	"errors"
	"fmt"
	"testing"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name  string
		pages map[int32]struct {
			items []string
			total int32
			err   error
		}
		wantItems  []string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "single page",
			pages: map[int32]struct {
				items []string
				total int32
				err   error
			}{
				0: {items: []string{"a", "b", "c"}, total: 3},
			},
			wantItems: []string{"a", "b", "c"},
		},
		{
			name: "multiple pages",
			pages: map[int32]struct {
				items []string
				total int32
				err   error
			}{
				0:   {items: makeStrings(int(defaultLimit)), total: 120},
				50:  {items: makeStrings(int(defaultLimit)), total: 120},
				100: {items: makeStrings(20), total: 120},
			},
			wantItems: makeStrings(int(defaultLimit)*2 + 20),
		},
		{
			name: "empty result",
			pages: map[int32]struct {
				items []string
				total int32
				err   error
			}{
				0: {items: []string{}, total: 0},
			},
			wantItems: nil,
		},
		{
			name: "error on first page",
			pages: map[int32]struct {
				items []string
				total int32
				err   error
			}{
				0: {err: errors.New("api error")},
			},
			wantErr:    true,
			wantErrMsg: "api error",
		},
		{
			name: "filtered count (GetTestRuns pattern)",
			pages: map[int32]struct {
				items []string
				total int32
				err   error
			}{
				0:  {items: makeStrings(int(defaultLimit)), total: 80},
				50: {items: makeStrings(30), total: 80},
			},
			wantItems: makeStrings(int(defaultLimit) + 30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := paginate(func(offset int32) ([]string, int32, error) {
				page, ok := tt.pages[offset]
				if !ok {
					t.Fatalf("unexpected offset %d", offset)
				}
				return page.items, page.total, page.err
			})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantErrMsg {
					t.Errorf("error = %q, want %q", err.Error(), tt.wantErrMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.wantItems) {
				t.Errorf("got %d items, want %d", len(got), len(tt.wantItems))
			}
		})
	}
}

func TestPaginationLimit_Default(t *testing.T) {
	got := paginationLimit()
	if got != defaultLimit {
		t.Errorf("paginationLimit() = %d, want %d", got, defaultLimit)
	}
}

func TestPaginationLimit_Custom(t *testing.T) {
	t.Setenv("QASE_TESTOPS_PAGINATION_LIMIT", "100")
	got := paginationLimit()
	if got != 100 {
		t.Errorf("paginationLimit() = %d, want 100", got)
	}
}

func TestPaginationLimit_InvalidValue(t *testing.T) {
	t.Setenv("QASE_TESTOPS_PAGINATION_LIMIT", "0")
	got := paginationLimit()
	if got != defaultLimit {
		t.Errorf("paginationLimit() with 0 = %d, want %d (default)", got, defaultLimit)
	}
}

func TestPaginate_CustomLimit(t *testing.T) {
	t.Setenv("QASE_TESTOPS_PAGINATION_LIMIT", "10")

	callCount := 0
	got, err := paginate(func(offset int32) ([]string, int32, error) {
		callCount++
		expectedOffset := int32((callCount - 1) * 10)
		if offset != expectedOffset {
			t.Errorf("call %d: offset = %d, want %d", callCount, offset, expectedOffset)
		}
		if callCount == 1 {
			return makeStrings(10), 25, nil
		}
		if callCount == 2 {
			return makeStrings(10), 25, nil
		}
		return makeStrings(5), 25, nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 25 {
		t.Errorf("got %d items, want 25", len(got))
	}
	if callCount != 3 {
		t.Errorf("expected 3 pages with limit=10, got %d calls", callCount)
	}
}

// makeStrings creates n sequential string items for testing
func makeStrings(n int) []string {
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = fmt.Sprintf("item-%d", i)
	}
	return s
}
