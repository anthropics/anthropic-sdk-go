package pagination

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type intPager interface {
	Next() bool
	Current() int
	Err() error
	Index() int
}

func loadIntPager(t *testing.T, style, endpoint string) intPager {
	t.Helper()
	load := func(dst any) *requestconfig.RequestConfig {
		cfg, err := requestconfig.NewRequestConfig(context.Background(), http.MethodGet, endpoint, nil, dst,
			option.WithBaseURL(endpoint), option.WithMaxRetries(0))
		if err != nil {
			t.Fatal(err)
		}
		if err := cfg.Execute(); err != nil {
			t.Fatal(err)
		}
		return cfg
	}
	switch style {
	case "page-after", "page-before":
		var page *Page[int]
		cfg := load(&page)
		page.SetPageConfig(cfg, nil)
		return NewPageAutoPager(page, nil)
	case "token":
		var page *TokenPage[int]
		cfg := load(&page)
		page.SetPageConfig(cfg, nil)
		return NewTokenPageAutoPager(page, nil)
	case "cursor":
		var page *PageCursor[int]
		cfg := load(&page)
		page.SetPageConfig(cfg, nil)
		return NewPageCursorAutoPager(page, nil)
	case "bidirectional":
		var page *BidirectionalPageCursor[int]
		cfg := load(&page)
		page.SetPageConfig(cfg, nil)
		return NewBidirectionalPageCursorAutoPager(page, nil)
	default:
		t.Fatalf("unknown style %s", style)
		return nil
	}
}

func TestAutoPagersContinueThroughEmptyPages(t *testing.T) {
	styles := []struct{ name, query string }{
		{"page-after", "after_id"}, {"page-before", "before_id"},
		{"token", "page_token"}, {"cursor", "page"}, {"bidirectional", "page"},
	}
	cases := []struct {
		name     string
		pages    [][]int
		want     []int
		failPage int
	}{
		{name: "empty first", pages: [][]int{{}, {1, 2}}, want: []int{1, 2}},
		{name: "empty middle", pages: [][]int{{1}, {}, {2}}, want: []int{1, 2}},
		{name: "consecutive empty", pages: [][]int{{1}, {}, {}, {2}}, want: []int{1, 2}},
		{name: "empty final", pages: [][]int{{1}, {}}, want: []int{1}},
		{name: "empty only", pages: [][]int{{}}, want: []int{}},
		{name: "error after empty", pages: [][]int{{}, {1}}, want: []int{}, failPage: 1},
	}
	for _, style := range styles {
		for _, tc := range cases {
			t.Run(style.name+"/"+tc.name, func(t *testing.T) {
				requests := 0
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					page := requests
					requests++
					if page >= len(tc.pages) {
						t.Errorf("unexpected request %d", page)
						http.Error(w, "too many requests", 500)
						return
					}
					if page > 0 {
						wantCursor := fmt.Sprintf("cursor-%d+=/", page)
						if got := req.URL.Query().Get(style.query); got != wantCursor {
							t.Errorf("cursor = %q, want %q", got, wantCursor)
						}
					}
					if req.URL.Query().Get("limit") != "2" {
						t.Error("lost original query parameter")
					}
					w.Header().Set("Content-Type", "application/json")
					if tc.failPage > 0 && page == tc.failPage {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(w, `{"type":"error","error":{"type":"api_error","message":"test failure"}}`)
						return
					}
					next := ""
					if page+1 < len(tc.pages) {
						next = fmt.Sprintf("cursor-%d+=/", page+1)
					}
					payload := map[string]any{
						"data": tc.pages[page], "has_more": next != "",
						"first_id": next, "last_id": next, "next_page": next,
					}
					if err := json.NewEncoder(w).Encode(payload); err != nil {
						t.Error(err)
					}
				}))
				defer srv.Close()
				endpoint := srv.URL + "/items?limit=2"
				if style.name == "page-before" {
					endpoint += "&before_id=start"
				}
				pager := loadIntPager(t, style.name, endpoint)
				got := []int{}
				for pager.Next() {
					got = append(got, pager.Current())
					if pager.Index() != len(got) {
						t.Fatalf("Index() = %d", pager.Index())
					}
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("items = %v, want %v", got, tc.want)
				}
				if (pager.Err() != nil) != (tc.failPage > 0) {
					t.Errorf("Err() = %v", pager.Err())
				}
				if requests != len(tc.pages) {
					t.Errorf("requests = %d, want %d", requests, len(tc.pages))
				}
				if pager.Next() {
					t.Error("iterator restarted after completion")
				}
				if requests != len(tc.pages) {
					t.Error("iterator made another request after completion")
				}
			})
		}
	}
}

func TestEmptyPageHonorsHasMoreFalse(t *testing.T) {
	for _, style := range []string{"page-after", "page-before", "token"} {
		t.Run(style, func(t *testing.T) {
			requests := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				requests++
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"data":[],"has_more":false,"first_id":"unused","last_id":"unused","next_page":"unused"}`)
			}))
			defer srv.Close()
			endpoint := srv.URL + "/items"
			if style == "page-before" {
				endpoint += "?before_id=start"
			}
			pager := loadIntPager(t, style, endpoint)
			if pager.Next() || pager.Err() != nil {
				t.Fatalf("unexpected item/error: %v", pager.Err())
			}
			if requests != 1 {
				t.Fatalf("requests = %d, want 1", requests)
			}
		})
	}
}
