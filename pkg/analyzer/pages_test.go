package analyzer

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/timezone"
	"github.com/pirsch-analytics/pirsch/v6/pkg"
	"github.com/pirsch-analytics/pirsch/v6/pkg/db"
	"github.com/pirsch-analytics/pirsch/v6/pkg/model"
	"github.com/pirsch-analytics/pirsch/v6/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAnalyzer_Hostname(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home", Hostname: "example.com"},
		{VisitorID: 1, Time: util.PastDay(4).Add(time.Minute * 3), SessionID: 4, DurationSeconds: 180, Path: "/foo", Title: "Foo", Hostname: "example.com"},
		{VisitorID: 1, Time: util.PastDay(4).Add(time.Hour), SessionID: 41, Path: "/", Title: "Home", Hostname: "example.com"},
		{VisitorID: 2, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home", Hostname: "foo.com"},
		{VisitorID: 2, Time: util.PastDay(4).Add(time.Minute * 2), SessionID: 4, DurationSeconds: 120, Path: "/bar", Title: "Bar", Hostname: "foo.com"},
		{VisitorID: 3, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home", Hostname: "example.com"},
		{VisitorID: 4, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home", Hostname: "example.com"},
		{VisitorID: 5, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home", Hostname: "foo.com"},
		{VisitorID: 5, Time: util.PastDay(2).Add(time.Minute * 5), SessionID: 21, Path: "/bar", Title: "Bar", Hostname: "foo.com"},
		{VisitorID: 6, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home", Hostname: "foo.com"},
		{VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 10), SessionID: 2, DurationSeconds: 600, Path: "/bar", Title: "Bar", Hostname: "foo.com"},
		{VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 11), SessionID: 21, Path: "/bar", Title: "Bar", Hostname: "foo.com"},
		{VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 21), SessionID: 21, DurationSeconds: 600, Path: "/foo", Title: "Foo", Hostname: "foo.com"},
		{VisitorID: 7, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home", Hostname: "example.com"},
		{VisitorID: 8, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home", Hostname: "example.com"},
		{VisitorID: 9, Time: util.Today(), SessionID: 2, Path: "/", Title: "Home", Hostname: "bar.com"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.PastDay(4).Add(time.Minute * 3), Start: time.Now(), SessionID: 4, DurationSeconds: 180, ExitPath: "/foo", EntryTitle: "Foo", IsBounce: false, PageViews: 2, Hostname: "example.com"},
			{Sign: 1, VisitorID: 1, Time: util.PastDay(4).Add(time.Hour), Start: time.Now(), SessionID: 41, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "example.com"},
			{Sign: 1, VisitorID: 2, Time: util.PastDay(4).Add(time.Minute * 2), Start: time.Now(), SessionID: 4, DurationSeconds: 120, ExitPath: "/bar", EntryTitle: "Bar", IsBounce: false, PageViews: 2, Hostname: "foo.com"},
			{Sign: 1, VisitorID: 3, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "example.com"},
			{Sign: 1, VisitorID: 4, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "example.com"},
		},
		{
			{Sign: -1, VisitorID: 4, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "example.com"},
			{Sign: 1, VisitorID: 4, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/foo", EntryTitle: "Foo", IsBounce: false, PageViews: 2, Hostname: "example.com"},
			{Sign: 1, VisitorID: 5, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "foo.com"},
			{Sign: 1, VisitorID: 5, Time: util.PastDay(2).Add(time.Minute * 5), Start: time.Now(), SessionID: 21, ExitPath: "/bar", EntryTitle: "Bar", IsBounce: true, PageViews: 1, Hostname: "foo.com"},
			{Sign: 1, VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 10), Start: time.Now(), SessionID: 2, DurationSeconds: 600, ExitPath: "/bar", EntryTitle: "Bar", IsBounce: false, PageViews: 2, Hostname: "foo.com"},
			{Sign: 1, VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 21), Start: time.Now(), SessionID: 21, DurationSeconds: 600, ExitPath: "/foo", EntryTitle: "Foo", IsBounce: false, PageViews: 2, Hostname: "foo.com"},
			{Sign: 1, VisitorID: 7, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "example.com"},
			{Sign: 1, VisitorID: 8, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "example.com"},
			{Sign: 1, VisitorID: 9, Time: util.Today(), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1, Hostname: "bar.com"},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.Hostname(nil)
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "example.com", visitors[0].Hostname)
	assert.Equal(t, "foo.com", visitors[1].Hostname)
	assert.Equal(t, "bar.com", visitors[2].Hostname)
	assert.Equal(t, 5, visitors[0].Visitors)
	assert.Equal(t, 3, visitors[1].Visitors)
	assert.Equal(t, 1, visitors[2].Visitors)
	assert.Equal(t, 8, visitors[0].Views)
	assert.Equal(t, 8, visitors[1].Views)
	assert.Equal(t, 1, visitors[2].Views)
	assert.Equal(t, 6, visitors[0].Sessions)
	assert.Equal(t, 5, visitors[1].Sessions)
	assert.Equal(t, 1, visitors[2].Sessions)
	assert.Equal(t, 4, visitors[0].Bounces)
	assert.Equal(t, 2, visitors[1].Bounces)
	assert.Equal(t, 1, visitors[2].Bounces)
	assert.InDelta(t, 0.5555, visitors[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, visitors[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.1111, visitors[2].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.4705, visitors[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.4705, visitors[1].RelativeViews, 0.001)
	assert.InDelta(t, 0.0588, visitors[2].RelativeViews, 0.001)
	assert.InDelta(t, 0.6666, visitors[0].BounceRate, 0.001)
	assert.InDelta(t, 0.4, visitors[1].BounceRate, 0.001)
	assert.InDelta(t, 1, visitors[2].BounceRate, 0.001)

	visitors, err = analyzer.Pages.Hostname(&Filter{
		Hostname: []string{"bar.com"},
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "bar.com", visitors[0].Hostname)
	assert.Equal(t, 1, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[0].Views)
	assert.Equal(t, 1, visitors[0].Sessions)
	assert.Equal(t, 1, visitors[0].Bounces)
	assert.InDelta(t, 0.1111, visitors[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.0588, visitors[0].RelativeViews, 0.001)
	assert.InDelta(t, 1, visitors[0].BounceRate, 0.001)

	visitors, err = analyzer.Pages.Hostname(&Filter{
		Hostname: []string{"example.com"},
		Path:     []string{"/"},
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "example.com", visitors[0].Hostname)
	assert.Equal(t, 5, visitors[0].Visitors)
	assert.Equal(t, 6, visitors[0].Views)
	assert.Equal(t, 6, visitors[0].Sessions)
	assert.Equal(t, 4, visitors[0].Bounces)
	assert.InDelta(t, 0.5555, visitors[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3529, visitors[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.6666, visitors[0].BounceRate, 0.001)

	_, err = analyzer.Pages.Hostname(getMaxFilter(""))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Hostname(getMaxFilter("event"))
	assert.NoError(t, err)
}

// TODO remove after migration
func TestAnalyzer_HostnameFallback(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/entry", Hostname: "example.com"},
		{VisitorID: 1, Time: util.Today(), Path: "/exit", Hostname: "example.com"},
		{VisitorID: 2, Time: util.Today(), Path: "/", Hostname: "example.com"},
		{VisitorID: 2, Time: util.Today(), Path: "/exit", Hostname: "example.com"},
		{VisitorID: 3, Time: util.Today(), Path: "/entry"},
		{VisitorID: 3, Time: util.Today(), Path: "/"},
	}))
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, Time: util.Today(), Name: "event", Path: "/", Hostname: "example.com"},
		{VisitorID: 3, Time: util.Today(), Name: "event", Path: "/"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/entry", ExitPath: "/exit", PageViews: 2, Hostname: "example.com"},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/exit", PageViews: 2, Hostname: "example.com"},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), EntryPath: "/entry", ExitPath: "/", PageViews: 2},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.Hostname(&Filter{
		HostnameFallback: "example.com",
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "example.com", visitors[0].Hostname)
	assert.Equal(t, 3, visitors[0].Visitors)
	assert.Equal(t, 6, visitors[0].Views)
	assert.Equal(t, 3, visitors[0].Sessions)
	assert.Equal(t, 0, visitors[0].Bounces)
	assert.InDelta(t, 1, visitors[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 1, visitors[0].RelativeViews, 0.001)
	assert.InDelta(t, 0, visitors[0].BounceRate, 0.001)
	pages, err := analyzer.Pages.ByPath(&Filter{
		HostnameFallback: "example.com",
	})
	assert.NoError(t, err)
	assert.Len(t, pages, 3)
	assert.Equal(t, "example.com", pages[0].Hostname)
	assert.Equal(t, "example.com", pages[1].Hostname)
	assert.Equal(t, "example.com", pages[2].Hostname)
	assert.Equal(t, "/", pages[0].Path)
	assert.Equal(t, "/entry", pages[1].Path)
	assert.Equal(t, "/exit", pages[2].Path)
	assert.Equal(t, 2, pages[0].Visitors)
	assert.Equal(t, 2, pages[1].Visitors)
	assert.Equal(t, 2, pages[2].Visitors)
	eventPages, err := analyzer.Pages.ByEventPath(&Filter{
		EventName:        []string{"event"},
		HostnameFallback: "example.com",
	})
	assert.NoError(t, err)
	assert.Len(t, eventPages, 1)
	assert.Equal(t, "example.com", eventPages[0].Hostname)
	assert.Equal(t, "/", eventPages[0].Path)
	assert.Equal(t, 2, eventPages[0].Visitors)
	entryPages, err := analyzer.Pages.Entry(&Filter{
		HostnameFallback: "example.com",
	})
	assert.NoError(t, err)
	assert.Len(t, entryPages, 2)
	assert.Equal(t, "example.com", entryPages[0].Hostname)
	assert.Equal(t, "example.com", entryPages[1].Hostname)
	assert.Equal(t, "/entry", entryPages[0].Path)
	assert.Equal(t, "/", entryPages[1].Path)
	assert.Equal(t, 2, entryPages[0].Entries)
	assert.Equal(t, 1, entryPages[1].Entries)
	exitPages, err := analyzer.Pages.Exit(&Filter{
		HostnameFallback: "example.com",
	})
	assert.NoError(t, err)
	assert.Len(t, exitPages, 2)
	assert.Equal(t, "example.com", exitPages[0].Hostname)
	assert.Equal(t, "example.com", exitPages[1].Hostname)
	assert.Equal(t, "/exit", exitPages[0].Path)
	assert.Equal(t, "/", exitPages[1].Path)
	assert.Equal(t, 2, exitPages[0].Exits)
	assert.Equal(t, 1, exitPages[1].Exits)
}

func TestAnalyzer_ByPathAndAvgTimeOnPage(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home"},
		{VisitorID: 1, Time: util.PastDay(4).Add(time.Minute * 3), SessionID: 4, DurationSeconds: 180, Path: "/foo", Title: "Foo"},
		{VisitorID: 1, Time: util.PastDay(4).Add(time.Hour), SessionID: 41, Path: "/", Title: "Home"},
		{VisitorID: 2, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home"},
		{VisitorID: 2, Time: util.PastDay(4).Add(time.Minute * 2), SessionID: 4, DurationSeconds: 120, Path: "/bar", Title: "Bar"},
		{VisitorID: 3, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home"},
		{VisitorID: 4, Time: util.PastDay(4), SessionID: 4, Path: "/", Title: "Home"},
		{VisitorID: 5, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home"},
		{VisitorID: 5, Time: util.PastDay(2).Add(time.Minute * 5), SessionID: 21, Path: "/bar", Title: "Bar"},
		{VisitorID: 6, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home"},
		{VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 10), SessionID: 2, DurationSeconds: 600, Path: "/bar", Title: "Bar"},
		{VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 11), SessionID: 21, Path: "/bar", Title: "Bar"},
		{VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 21), SessionID: 21, DurationSeconds: 600, Path: "/foo", Title: "Foo"},
		{VisitorID: 7, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home"},
		{VisitorID: 8, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home"},
		{VisitorID: 9, Time: util.Today(), SessionID: 2, Path: "/", Title: "Home"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.PastDay(4).Add(time.Minute * 3), Start: time.Now(), SessionID: 4, DurationSeconds: 180, ExitPath: "/foo", EntryTitle: "Foo", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 1, Time: util.PastDay(4).Add(time.Hour), Start: time.Now(), SessionID: 41, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 2, Time: util.PastDay(4).Add(time.Minute * 2), Start: time.Now(), SessionID: 4, DurationSeconds: 120, ExitPath: "/bar", EntryTitle: "Bar", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 3, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 4, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
		},
		{
			{Sign: -1, VisitorID: 4, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 4, Time: util.PastDay(4), Start: time.Now(), SessionID: 4, ExitPath: "/foo", EntryTitle: "Foo", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 5, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 5, Time: util.PastDay(2).Add(time.Minute * 5), Start: time.Now(), SessionID: 21, ExitPath: "/bar", EntryTitle: "Bar", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 10), Start: time.Now(), SessionID: 2, DurationSeconds: 600, ExitPath: "/bar", EntryTitle: "Bar", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 6, Time: util.PastDay(2).Add(time.Minute * 21), Start: time.Now(), SessionID: 21, DurationSeconds: 600, ExitPath: "/foo", EntryTitle: "Foo", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 7, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 8, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 9, Time: util.Today(), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home", IsBounce: true, PageViews: 1},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByPath(&Filter{IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/bar", visitors[1].Path)
	assert.Equal(t, "/foo", visitors[2].Path)
	assert.Empty(t, visitors[0].Title)
	assert.Empty(t, visitors[1].Title)
	assert.Empty(t, visitors[2].Title)
	assert.Equal(t, 9, visitors[0].Visitors)
	assert.Equal(t, 3, visitors[1].Visitors)
	assert.Equal(t, 2, visitors[2].Visitors)
	assert.InDelta(t, 1, visitors[0].RelativeVisitors, 0.01)
	assert.InDelta(t, 0.3333, visitors[1].RelativeVisitors, 0.01)
	assert.InDelta(t, 0.2222, visitors[2].RelativeVisitors, 0.01)
	assert.Equal(t, 10, visitors[0].Sessions)
	assert.Equal(t, 4, visitors[1].Sessions)
	assert.Equal(t, 2, visitors[2].Sessions)
	assert.Equal(t, 10, visitors[0].Views)
	assert.Equal(t, 4, visitors[1].Views)
	assert.Equal(t, 2, visitors[2].Views)
	assert.InDelta(t, 0.5882, visitors[0].RelativeViews, 0.01)
	assert.InDelta(t, 0.2352, visitors[1].RelativeViews, 0.01)
	assert.InDelta(t, 0.125, visitors[2].RelativeViews, 0.01)
	assert.Equal(t, 6, visitors[0].Bounces)
	assert.Equal(t, 1, visitors[1].Bounces)
	assert.Equal(t, 0, visitors[2].Bounces)
	assert.InDelta(t, 0.6, visitors[0].BounceRate, 0.01)
	assert.InDelta(t, 0.25, visitors[1].BounceRate, 0.01)
	assert.InDelta(t, 0, visitors[2].BounceRate, 0.01)
	assert.Equal(t, 300, visitors[0].AverageTimeSpentSeconds)
	assert.Equal(t, 600, visitors[1].AverageTimeSpentSeconds)
	assert.Equal(t, 0, visitors[2].AverageTimeSpentSeconds)
	visitors, err = analyzer.Pages.ByPath(&Filter{Sort: []Sort{
		{Field: FieldPath, Direction: pkg.DirectionDESC},
	}})
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "/foo", visitors[0].Path)
	assert.Equal(t, "/bar", visitors[1].Path)
	assert.Equal(t, "/", visitors[2].Path)
	top, err := analyzer.Time.AvgTimeOnPage(nil)
	assert.NoError(t, err)
	assert.Len(t, top, 2)
	assert.Equal(t, util.PastDay(4), top[0].Day.Time)
	assert.Equal(t, util.PastDay(2), top[1].Day.Time)
	assert.Equal(t, 150, top[0].AverageTimeSpentSeconds)
	assert.Equal(t, 600, top[1].AverageTimeSpentSeconds)
	ttop, err := analyzer.Visitors.totalTimeOnPage(&Filter{
		Ctx: context.Background(),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1500, ttop)
	visitors, err = analyzer.Pages.ByPath(&Filter{From: util.PastDay(3), To: util.PastDay(1), IncludeTitle: true, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/bar", visitors[1].Path)
	assert.Equal(t, "/foo", visitors[2].Path)
	assert.Equal(t, "Home", visitors[0].Title)
	assert.Equal(t, "Bar", visitors[1].Title)
	assert.Equal(t, "Foo", visitors[2].Title)
	assert.Equal(t, 600, visitors[0].AverageTimeSpentSeconds)
	assert.Equal(t, 600, visitors[1].AverageTimeSpentSeconds)
	assert.Equal(t, 0, visitors[2].AverageTimeSpentSeconds)
	top, err = analyzer.Time.AvgTimeOnPage(&Filter{From: util.PastDay(3), To: util.PastDay(1), IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, top, 3)
	assert.Equal(t, util.PastDay(3), top[0].Day.Time)
	assert.Equal(t, util.PastDay(2), top[1].Day.Time)
	assert.Equal(t, util.PastDay(1), top[2].Day.Time)
	assert.Equal(t, 0, top[0].AverageTimeSpentSeconds)
	assert.Equal(t, 600, top[1].AverageTimeSpentSeconds)
	assert.Equal(t, 0, top[2].AverageTimeSpentSeconds)
	ttop, err = analyzer.Visitors.totalTimeOnPage(&Filter{
		Ctx:  context.Background(),
		From: util.PastDay(3),
		To:   util.PastDay(1),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1200, ttop)
	_, err = analyzer.Pages.ByPath(getMaxFilter(""))
	assert.NoError(t, err)
	_, err = analyzer.Pages.ByPath(getMaxFilter("event"))
	assert.NoError(t, err)
	_, err = analyzer.Pages.ByPath(&Filter{Offset: 1, Limit: 10, Sort: []Sort{
		{
			Field:     FieldPath,
			Direction: pkg.DirectionASC,
		},
	}, Search: []Search{
		{
			Field: FieldPath,
			Input: "/",
		},
	}})
	assert.NoError(t, err)
	_, err = analyzer.Visitors.totalTimeOnPage(getMaxFilter(""))
	assert.NoError(t, err)
	visitors, err = analyzer.Pages.ByPath(&Filter{Limit: 1})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	ttop, err = analyzer.Visitors.totalTimeOnPage(&Filter{
		Ctx:                  context.Background(),
		MaxTimeOnPageSeconds: 200,
	})
	assert.NoError(t, err)
	assert.Equal(t, 180+120+200+200, ttop)
	visitors, err = analyzer.Pages.ByPath(&Filter{Search: []Search{{Field: FieldPath, Input: "%foo%"}}, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		ExitPath:          []string{"/foo"},
		Search:            []Search{{Field: FieldPath, Input: "%foo%"}},
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)

	// imported statistics
	past5Days := util.PastDay(5).Format(time.DateOnly)
	_, err = dbClient.Exec(fmt.Sprintf(`INSERT INTO "imported_page" (date, path, visitors, views, sessions, bounces) VALUES
		('%s', '/bar', 2, 4, 3, 1), ('%s', '/', 1, 2, 1, 1)`, past5Days, past5Days))
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		From:          util.PastDay(5),
		To:            util.Today(),
		ImportedUntil: util.PastDay(4),
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/bar", visitors[1].Path)
	assert.Equal(t, "/foo", visitors[2].Path)
	assert.Equal(t, 10, visitors[0].Visitors)
	assert.Equal(t, 5, visitors[1].Visitors)
	assert.Equal(t, 2, visitors[2].Visitors)
	assert.InDelta(t, 0.8333, visitors[0].RelativeVisitors, 0.01)
	assert.InDelta(t, 0.4166, visitors[1].RelativeVisitors, 0.01)
	assert.InDelta(t, 0.1666, visitors[2].RelativeVisitors, 0.01)
	assert.Equal(t, 11, visitors[0].Sessions)
	assert.Equal(t, 7, visitors[1].Sessions)
	assert.Equal(t, 2, visitors[2].Sessions)
	assert.Equal(t, 12, visitors[0].Views)
	assert.Equal(t, 8, visitors[1].Views)
	assert.Equal(t, 2, visitors[2].Views)
	assert.InDelta(t, 0.5217, visitors[0].RelativeViews, 0.01)
	assert.InDelta(t, 0.3478, visitors[1].RelativeViews, 0.01)
	assert.InDelta(t, 0.0869, visitors[2].RelativeViews, 0.01)
	assert.Equal(t, 7, visitors[0].Bounces)
	assert.Equal(t, 2, visitors[1].Bounces)
	assert.Equal(t, 0, visitors[2].Bounces)
	assert.InDelta(t, 0.6363, visitors[0].BounceRate, 0.01)
	assert.InDelta(t, 0.2857, visitors[1].BounceRate, 0.01)
	assert.InDelta(t, 0, visitors[2].BounceRate, 0.01)
	assert.Equal(t, 0, visitors[0].AverageTimeSpentSeconds)
	assert.Equal(t, 0, visitors[1].AverageTimeSpentSeconds)
	assert.Equal(t, 0, visitors[2].AverageTimeSpentSeconds)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		From:          util.PastDay(5),
		To:            util.Today(),
		ImportedUntil: util.PastDay(4),
		Path:          []string{"/"},
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, 10, visitors[0].Visitors)
	assert.InDelta(t, 0.8333, visitors[0].RelativeVisitors, 0.01)
	assert.Equal(t, 11, visitors[0].Sessions)
	assert.Equal(t, 12, visitors[0].Views)
	assert.InDelta(t, 0.5217, visitors[0].RelativeViews, 0.01)
	assert.Equal(t, 7, visitors[0].Bounces)
	assert.InDelta(t, 0.6363, visitors[0].BounceRate, 0.01)
	assert.Equal(t, 0, visitors[0].AverageTimeSpentSeconds)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		From:              util.PastDay(5),
		To:                util.Today(),
		ImportedUntil:     util.PastDay(4),
		Path:              []string{"/"},
		IncludeTitle:      true,
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "Home", visitors[0].Title)
	assert.Equal(t, 10, visitors[0].Visitors)
	assert.InDelta(t, 0.8333, visitors[0].RelativeVisitors, 0.01)
	assert.Equal(t, 11, visitors[0].Sessions)
	assert.Equal(t, 12, visitors[0].Views)
	assert.InDelta(t, 0.5217, visitors[0].RelativeViews, 0.01)
	assert.Equal(t, 7, visitors[0].Bounces)
	assert.InDelta(t, 0.6363, visitors[0].BounceRate, 0.01)
	assert.Equal(t, 300, visitors[0].AverageTimeSpentSeconds)
}

func TestAnalyzer_PageTitle(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.PastDay(2).Add(time.Hour * 23).Add(time.Minute * 57), SessionID: 1, Path: "/", Title: "Home 1"},
		{VisitorID: 1, Time: util.PastDay(1).Add(time.Second * 14), SessionID: 1, Path: "/", Title: "Home 2", DurationSeconds: 314},
		{VisitorID: 2, Time: util.Today(), SessionID: 3, Path: "/foo", Title: "Foo"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.PastDay(2), Start: time.Now(), SessionID: 1, ExitPath: "/foo", EntryTitle: "Foo"},
		},
		{
			{Sign: -1, VisitorID: 1, Time: util.PastDay(2), Start: time.Now(), SessionID: 1, ExitPath: "/foo", EntryTitle: "Foo"},
			{Sign: 1, VisitorID: 1, Time: util.PastDay(2), Start: time.Now(), SessionID: 1, ExitPath: "/", EntryTitle: "Home 1"},
			{Sign: 1, VisitorID: 1, Time: util.PastDay(1), Start: time.Now(), SessionID: 2, ExitPath: "/", EntryTitle: "Home 2", DurationSeconds: 42},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), SessionID: 3, ExitPath: "/foo", EntryTitle: "Foo"},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByPath(&Filter{IncludeTitle: true, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "Home 1", visitors[0].Title)
	assert.Equal(t, "Home 2", visitors[1].Title)
	assert.Equal(t, "Foo", visitors[2].Title)
	assert.Equal(t, 314, visitors[0].AverageTimeSpentSeconds)
	assert.Equal(t, 314, visitors[1].AverageTimeSpentSeconds)
	assert.Equal(t, 0, visitors[2].AverageTimeSpentSeconds)
}

func TestAnalyzer_PageTitleEvent(t *testing.T) {
	db.CleanupDB(t, dbClient)
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.PastDay(2), Start: time.Now(), EntryPath: "/", ExitPath: "/", EntryTitle: "Home 1", ExitTitle: "Home 1"},
		},
		{
			{Sign: -1, VisitorID: 1, SessionID: 1, Time: util.PastDay(2), Start: time.Now(), EntryPath: "/", ExitPath: "/", EntryTitle: "Home 1", ExitTitle: "Home 1"},
			{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.PastDay(2), Start: time.Now(), EntryPath: "/", ExitPath: "/foo", EntryTitle: "Home 1", ExitTitle: "Foo"},
			{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.PastDay(1), Start: time.Now(), EntryPath: "/", ExitPath: "/", EntryTitle: "Home 2", ExitTitle: "Home 2"},
			{Sign: 1, VisitorID: 2, SessionID: 3, Time: util.PastDay(1), Start: time.Now(), EntryPath: "/foo", ExitPath: "/foo", EntryTitle: "Foo", ExitTitle: "Foo"},
		},
	})
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.PastDay(2), SessionID: 1, Path: "/", Title: "Home 1"},
		{VisitorID: 1, Time: util.PastDay(2), SessionID: 1, Path: "/foo", Title: "Foo", DurationSeconds: 42},
		{VisitorID: 1, Time: util.PastDay(1), SessionID: 1, Path: "/", Title: "Home 2"},
		{VisitorID: 2, Time: util.PastDay(1), SessionID: 2, Path: "/foo", Title: "Foo"},
	}))
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{Name: "event", VisitorID: 1, Time: util.PastDay(2), SessionID: 1, Path: "/", Title: "Home 1"},
		{Name: "event", VisitorID: 1, Time: util.PastDay(1), SessionID: 1, Path: "/", Title: "Home 2", DurationSeconds: 42},
		{Name: "event", VisitorID: 2, Time: util.Today(), SessionID: 3, Path: "/foo", Title: "Foo"},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByPath(&Filter{EventName: []string{"event"}, IncludeTitle: true, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "Home 1", visitors[0].Title)
	assert.Equal(t, "Home 2", visitors[1].Title)
	assert.Equal(t, "Foo", visitors[2].Title)
	assert.Equal(t, 42, visitors[0].AverageTimeSpentSeconds)
	assert.Equal(t, 42, visitors[1].AverageTimeSpentSeconds)
	assert.Equal(t, 0, visitors[2].AverageTimeSpentSeconds)
}

func TestAnalyzer_ByPathEvent(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/"},
		{VisitorID: 1, Time: util.Today(), Path: "/foo"},
		{VisitorID: 1, Time: util.Today(), Path: "/bar"},
		{VisitorID: 2, Time: util.Today(), Path: "/foo"},
		{VisitorID: 3, Time: util.Today(), Path: "/"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/foo", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/bar", IsBounce: false, PageViews: 3},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), EntryPath: "/foo", ExitPath: "/foo", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", IsBounce: true, PageViews: 1},
		},
	})
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, Time: util.Today(), Name: "event"},
		{VisitorID: 3, Time: util.Today(), Name: "event"},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByPath(&Filter{EventName: []string{"event"}})
	assert.NoError(t, err)
	assert.Len(t, visitors, 3)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/bar", visitors[1].Path)
	assert.Equal(t, "/foo", visitors[2].Path)
	assert.Equal(t, 2, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	assert.Equal(t, 1, visitors[2].Visitors)
	visitors, err = analyzer.Pages.ByPath(&Filter{EventName: []string{"!event"}})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "/foo", visitors[0].Path)
	assert.Equal(t, 1, visitors[0].Visitors)

	entries, err := analyzer.Pages.Entry(&Filter{EventName: []string{"event"}})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, 2, entries[0].Entries)
	assert.Equal(t, 2, entries[0].Visitors)
	entries, err = analyzer.Pages.Entry(&Filter{EventName: []string{"!event"}})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/foo", entries[0].Path)
	assert.Equal(t, 1, entries[0].Entries)
	assert.Equal(t, 2, entries[0].Visitors)

	exits, err := analyzer.Pages.Exit(&Filter{EventName: []string{"event"}})
	assert.NoError(t, err)
	assert.Len(t, exits, 2)
	assert.Equal(t, "/", exits[0].Path)
	assert.Equal(t, "/bar", exits[1].Path)
	assert.Equal(t, 1, exits[0].Exits)
	assert.Equal(t, 1, exits[1].Exits)
	assert.Equal(t, 2, exits[0].Visitors)
	assert.Equal(t, 1, exits[1].Visitors)
	exits, err = analyzer.Pages.Exit(&Filter{EventName: []string{"!event"}})
	assert.NoError(t, err)
	assert.Len(t, exits, 1)
	assert.Equal(t, "/foo", exits[0].Path)
	assert.Equal(t, 1, exits[0].Exits)
	assert.Equal(t, 2, exits[0].Visitors)
}

func TestAnalyzer_ByPathEventPath(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/", Title: "Home"},
		{VisitorID: 1, Time: util.Today(), Path: "/foo", Title: "Foo"},
		{VisitorID: 1, Time: util.Today(), Path: "/bar", Title: "Bar"},
		{VisitorID: 2, Time: util.Today(), Path: "/foo", Title: "Foo"},
		{VisitorID: 3, Time: util.Today(), Path: "/", Title: "Home"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", EntryTitle: "Home", ExitTitle: "Home", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/foo", EntryTitle: "Home", ExitTitle: "Foo", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/bar", EntryTitle: "Home", ExitTitle: "Bar", IsBounce: false, PageViews: 3},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), EntryPath: "/foo", ExitPath: "/foo", EntryTitle: "Foo", ExitTitle: "Foo", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", EntryTitle: "Home", ExitTitle: "Home", IsBounce: true, PageViews: 1},
		},
	})
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, Time: util.Today(), Name: "event", Path: "/", Title: "Home", DurationSeconds: 5},
		{VisitorID: 3, Time: util.Today(), Name: "event", Path: "/foo", Title: "Foo", DurationSeconds: 9},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByEventPath(&Filter{
		EventName: []string{"event"},
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/foo", visitors[1].Path)
	assert.Equal(t, 1, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	visitors, err = analyzer.Pages.ByEventPath(&Filter{
		EventName: []string{"!event"},
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Empty(t, visitors[0].Path) // "unknown"
	assert.Equal(t, 1, visitors[0].Visitors)
	visitors, err = analyzer.Pages.ByEventPath(&Filter{
		EventName:    []string{"event"},
		IncludeTitle: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/foo", visitors[1].Path)
	assert.Equal(t, "Home", visitors[0].Title)
	assert.Equal(t, "Foo", visitors[1].Title)
	assert.Equal(t, 1, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	visitors, err = analyzer.Pages.ByEventPath(&Filter{
		EventName:    []string{"!event"},
		IncludeTitle: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Empty(t, visitors[0].Path)  // "unknown"
	assert.Empty(t, visitors[0].Title) // "unknown"
	assert.Equal(t, 1, visitors[0].Visitors)
	tz, err := timezone.Load("Europe/Berlin")
	assert.NoError(t, err)
	visitors, err = analyzer.Pages.ByEventPath(&Filter{
		Timezone:  tz,
		From:      util.PastDay(1),
		To:        util.Today(),
		EventName: []string{"event"},
		Path:      []string{"/foo"},
		Period:    pkg.PeriodDay,
		Limit:     10,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "/foo", visitors[0].Path)
	assert.Equal(t, 1, visitors[0].Visitors)
	visitors, err = analyzer.Pages.ByEventPath(&Filter{
		From:         util.Today(),
		To:           util.Today(),
		EventName:    []string{"event"},
		Path:         []string{"/foo"},
		Period:       pkg.PeriodDay,
		Limit:        10,
		IncludeTitle: true,
		Sample:       10_000_000,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, "/foo", visitors[0].Path)
	assert.Equal(t, "Foo", visitors[0].Title)
	assert.Equal(t, 1, visitors[0].Visitors)
	_, err = analyzer.Pages.ByEventPath(&Filter{
		EventName:         []string{"event"},
		Search:            []Search{{Field: FieldEventPath, Input: "%foo%"}},
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)
	_, err = analyzer.Pages.ByEventPath(&Filter{
		EventName:         []string{"event"},
		EntryPath:         []string{"/"},
		Search:            []Search{{Field: FieldEventPath, Input: "%foo%"}},
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)
}

func TestAnalyzer_ByPathTags(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/", Title: "Home", TagKeys: []string{"author", "foo"}, TagValues: []string{"John", "bar"}},
		{VisitorID: 1, Time: util.Today(), Path: "/foo", Title: "Foo", TagKeys: []string{"author", "title"}, TagValues: []string{"John", "Foo"}},
		{VisitorID: 1, Time: util.Today(), Path: "/bar", Title: "Bar", TagKeys: []string{"foo"}, TagValues: []string{"bar"}},
		{VisitorID: 2, Time: util.Today(), Path: "/foo", Title: "Foo"},
		{VisitorID: 3, Time: util.Today(), Path: "/", Title: "Home", TagKeys: []string{"author", "key"}, TagValues: []string{"Alice", "bar"}},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/bar", EntryTitle: "Home", ExitTitle: "Bar", IsBounce: false, PageViews: 3},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), EntryPath: "/foo", ExitPath: "/foo", EntryTitle: "Foo", ExitTitle: "Foo", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", EntryTitle: "Home", ExitTitle: "Home", IsBounce: true, PageViews: 1},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByPath(&Filter{
		Tags: map[string]string{"foo": "bar"},
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/bar", visitors[1].Path)
	assert.Equal(t, 1, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		Tags: map[string]string{"foo": "!bar"},
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	assert.Equal(t, visitors[0].Path, "/foo")
	assert.Equal(t, visitors[1].Path, "/")
	assert.Equal(t, 2, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		Tags:         map[string]string{"foo": "bar"},
		IncludeTitle: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	assert.Equal(t, "/", visitors[0].Path)
	assert.Equal(t, "/bar", visitors[1].Path)
	assert.Equal(t, "Home", visitors[0].Title)
	assert.Equal(t, "Bar", visitors[1].Title)
	assert.Equal(t, 1, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		Tags:         map[string]string{"foo": "!bar"},
		IncludeTitle: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	assert.Equal(t, visitors[0].Path, "/foo")
	assert.Equal(t, visitors[1].Path, "/")
	assert.Equal(t, visitors[0].Title, "Foo")
	assert.Equal(t, visitors[1].Title, "Home")
	assert.Equal(t, 2, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		EntryPath:    []string{"/"},
		ExitPath:     []string{"/bar"},
		Tags:         map[string]string{"foo": "bar"},
		IncludeTitle: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	assert.Equal(t, visitors[0].Path, "/")
	assert.Equal(t, visitors[1].Path, "/bar")
	assert.Equal(t, visitors[0].Title, "Home")
	assert.Equal(t, visitors[1].Title, "Bar")
	assert.Equal(t, 1, visitors[0].Visitors)
	assert.Equal(t, 1, visitors[1].Visitors)
	visitors, err = analyzer.Pages.ByPath(&Filter{
		EntryPath:    []string{"/"},
		ExitPath:     []string{"/bar"},
		Tags:         map[string]string{"title": "Foo"},
		IncludeTitle: true,
	})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	assert.Equal(t, visitors[0].Path, "/foo")
	assert.Equal(t, visitors[0].Title, "Foo")
	assert.Equal(t, 1, visitors[0].Visitors)
}

func TestAnalyzer_ByPathManyPages(t *testing.T) {
	db.CleanupDB(t, dbClient)
	var sessions []model.Session
	var pageViews []model.PageView

	for i := 0; i < 100_000; i++ {
		sessions = append(sessions, model.Session{
			Sign:            1,
			VisitorID:       uint64(i + 1),
			Time:            util.Today(),
			Start:           time.Now(),
			SessionID:       1,
			DurationSeconds: 60,
			PageViews:       1,
		})
		pageViews = append(pageViews, model.PageView{
			VisitorID: uint64(i + 1),
			Time:      util.Today(),
			SessionID: 1,
			Path:      fmt.Sprintf("/relatively/long/unique/path/%d", i),
		})
	}

	saveSessions(t, [][]model.Session{sessions})
	assert.NoError(t, dbClient.SavePageViews(pageViews))
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByPath(&Filter{IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, visitors, 100_000)
}

func TestAnalyzer_PathPattern(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/"},
		{VisitorID: 2, Time: util.Today(), Path: "/simple/page"},
		{VisitorID: 3, Time: util.Today(), Path: "/siMple/page/"},
		{VisitorID: 4, Time: util.Today(), Path: "/simple/page/with/many/slashes"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), ExitPath: "/exit"},
		},
		{
			{Sign: -1, VisitorID: 1, Time: util.Today(), Start: time.Now(), ExitPath: "/exit"},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), ExitPath: "/"},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), ExitPath: "/simple/page"},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), ExitPath: "/siMple/page/"},
			{Sign: 1, VisitorID: 4, Time: util.Today(), Start: time.Now(), ExitPath: "/simple/page/with/many/slashes"},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	visitors, err := analyzer.Pages.ByPath(nil)
	assert.NoError(t, err)
	assert.Len(t, visitors, 4)
	visitors, err = analyzer.Pages.ByPath(&Filter{PathPattern: []string{"(?i)^/simple/[^/]+$"}})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
	visitors, err = analyzer.Pages.ByPath(&Filter{PathPattern: []string{"(?i)^/simple/[^/]+/.*"}})
	assert.NoError(t, err)
	assert.Len(t, visitors, 2)
	visitors, err = analyzer.Pages.ByPath(&Filter{PathPattern: []string{"(?i)^/simple/[^/]+/slashes$"}})
	assert.NoError(t, err)
	assert.Len(t, visitors, 0)
	visitors, err = analyzer.Pages.ByPath(&Filter{PathPattern: []string{"(?i)^/simple/.+/slashes$"}})
	assert.NoError(t, err)
	assert.Len(t, visitors, 1)
}

func TestAnalyzer_EntryExitPages(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.PastDay(2), SessionID: 1, Path: "/", Title: "Home"},
		{VisitorID: 1, Time: util.PastDay(2).Add(time.Second), SessionID: 2, Path: "/", Title: "Home"},
		{VisitorID: 1, Time: util.PastDay(2).Add(time.Second * 10), SessionID: 2, DurationSeconds: 10, Path: "/foo", Title: "Foo"},
		{VisitorID: 2, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home"},
		{VisitorID: 3, Time: util.PastDay(2), SessionID: 2, Path: "/", Title: "Home"},
		{VisitorID: 4, Time: util.PastDay(1), SessionID: 1, Path: "/", Title: "Home"},
		{VisitorID: 4, Time: util.PastDay(1).Add(time.Second * 20), SessionID: 1, DurationSeconds: 20, Path: "/bar", Title: "Bar"},
		{VisitorID: 5, Time: util.PastDay(1), SessionID: 1, Path: "/", Title: "Home"},
		{VisitorID: 5, Time: util.PastDay(1).Add(time.Second * 40), SessionID: 1, DurationSeconds: 40, Path: "/bar", Title: "Bar"},
		{VisitorID: 6, Time: util.PastDay(1), SessionID: 1, Path: "/bar", Title: "Bar"},
		{VisitorID: 7, Time: util.PastDay(1), SessionID: 1, Path: "/bar", Title: "Bar"},
		{VisitorID: 7, Time: util.PastDay(1).Add(time.Minute), SessionID: 2, Path: "/", Title: "Home"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.PastDay(2).Add(time.Second * 10), Start: time.Now(), SessionID: 1, DurationSeconds: 10, EntryPath: "/bar", ExitPath: "/foo", EntryTitle: "Bar", ExitTitle: "Foo"},
		},
		{
			{Sign: -1, VisitorID: 1, Time: util.PastDay(2).Add(time.Second * 10), Start: time.Now(), SessionID: 1, DurationSeconds: 10, EntryPath: "/bar", ExitPath: "/foo", EntryTitle: "Bar", ExitTitle: "Foo"},
			{Sign: 1, VisitorID: 1, Time: util.PastDay(2).Add(time.Second * 10), Start: time.Now(), SessionID: 1, DurationSeconds: 10, EntryPath: "/", ExitPath: "/", EntryTitle: "Home", ExitTitle: "Home"},
			{Sign: 1, VisitorID: 1, Time: util.PastDay(2).Add(time.Second * 10), Start: time.Now(), SessionID: 2, DurationSeconds: 10, EntryPath: "/", ExitPath: "/foo", EntryTitle: "Home", ExitTitle: "Foo"},
			{Sign: 1, VisitorID: 2, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, EntryPath: "/", ExitPath: "/", EntryTitle: "Home", ExitTitle: "Home"},
			{Sign: 1, VisitorID: 3, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, EntryPath: "/", ExitPath: "/", EntryTitle: "Home", ExitTitle: "Home"},
			{Sign: 1, VisitorID: 4, Time: util.PastDay(1).Add(time.Second * 20), Start: time.Now(), SessionID: 1, DurationSeconds: 20, EntryPath: "/", ExitPath: "/bar", EntryTitle: "Home", ExitTitle: "Bar"},
			{Sign: 1, VisitorID: 5, Time: util.PastDay(1).Add(time.Second * 40), Start: time.Now(), SessionID: 1, DurationSeconds: 40, EntryPath: "/", ExitPath: "/bar", EntryTitle: "Home", ExitTitle: "Bar"},
			{Sign: 1, VisitorID: 6, Time: util.PastDay(1), Start: time.Now(), SessionID: 1, EntryPath: "/bar", ExitPath: "/bar", EntryTitle: "Bar", ExitTitle: "Bar"},
			{Sign: 1, VisitorID: 7, Time: util.PastDay(1).Add(time.Minute), Start: time.Now(), SessionID: 1, EntryPath: "/bar", ExitPath: "/bar", EntryTitle: "Bar", ExitTitle: "Bar"},
			{Sign: 1, VisitorID: 7, Time: util.PastDay(1).Add(time.Minute), Start: time.Now(), SessionID: 2, EntryPath: "/", ExitPath: "/", EntryTitle: "Home", ExitTitle: "Home"},
		},
	})
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, SessionID: 2, Time: util.PastDay(2).Add(time.Second), Path: "/", Name: "event", MetaKeys: []string{"key", "author", "amount"}, MetaValues: []string{"value", "John", "99.99"}},
		{VisitorID: 3, SessionID: 2, Time: util.PastDay(2), Path: "/", Name: "event", MetaKeys: []string{"author", "type"}, MetaValues: []string{"Alice", "blog_post"}},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	entries, err := analyzer.Pages.Entry(&Filter{IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, "/bar", entries[1].Path)
	assert.Empty(t, entries[0].Title)
	assert.Empty(t, entries[1].Title)
	assert.Equal(t, 6, entries[0].Visitors)
	assert.Equal(t, 4, entries[1].Visitors)
	assert.Equal(t, 7, entries[0].Sessions)
	assert.Equal(t, 4, entries[1].Sessions)
	assert.Equal(t, 7, entries[0].Entries)
	assert.Equal(t, 2, entries[1].Entries)
	assert.InDelta(t, 0.7777, entries[0].EntryRate, 0.001)
	assert.InDelta(t, 0.2222, entries[1].EntryRate, 0.001)
	assert.Equal(t, 23, entries[0].AverageTimeSpentSeconds)
	assert.Equal(t, 0, entries[1].AverageTimeSpentSeconds)
	entries, err = analyzer.Pages.Entry(&Filter{PathPattern: []string{"(?i)^/.*$"}})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, "/bar", entries[1].Path)
	assert.Equal(t, 6, entries[0].Visitors)
	assert.Equal(t, 4, entries[1].Visitors)
	assert.Equal(t, 7, entries[0].Sessions)
	assert.Equal(t, 4, entries[1].Sessions)
	assert.Equal(t, 7, entries[0].Entries)
	assert.Equal(t, 2, entries[1].Entries)
	entries, err = analyzer.Pages.Entry(&Filter{From: util.PastDay(1), To: util.Today(), IncludeTitle: true, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, "/bar", entries[1].Path)
	assert.Equal(t, "Home", entries[0].Title)
	assert.Equal(t, "Bar", entries[1].Title)
	assert.Equal(t, 3, entries[0].Visitors)
	assert.Equal(t, 4, entries[1].Visitors)
	assert.Equal(t, 3, entries[0].Sessions)
	assert.Equal(t, 4, entries[1].Sessions)
	assert.Equal(t, 3, entries[0].Entries)
	assert.Equal(t, 2, entries[1].Entries)
	assert.InDelta(t, 0.6, entries[0].EntryRate, 0.001)
	assert.InDelta(t, 0.4, entries[1].EntryRate, 0.001)
	assert.Equal(t, 30, entries[0].AverageTimeSpentSeconds)
	assert.Equal(t, 0, entries[1].AverageTimeSpentSeconds)
	entries, err = analyzer.Pages.Entry(&Filter{From: util.PastDay(1), To: util.Today(), EntryPath: []string{"/"}, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, 3, entries[0].Visitors)
	assert.Equal(t, 3, entries[0].Entries)
	assert.InDelta(t, 0.6, entries[0].EntryRate, 0.001)
	assert.Equal(t, 30, entries[0].AverageTimeSpentSeconds)
	_, err = analyzer.Pages.Entry(&Filter{Path: []string{"/bar"}, IncludeTitle: true})
	assert.NoError(t, err)
	_, err = analyzer.Pages.Entry(getMaxFilter(""))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Entry(getMaxFilter("event"))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Entry(&Filter{Offset: 1, Limit: 10, Sort: []Sort{
		{
			Field:     FieldEntryPath,
			Direction: pkg.DirectionASC,
		},
	}, Search: []Search{
		{
			Field: FieldEntryPath,
			Input: "/",
		},
	}})
	assert.NoError(t, err)
	_, err = analyzer.Pages.Entry(&Filter{Search: []Search{{Field: FieldPath, Input: "%foo%"}}, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	_, err = analyzer.Pages.Entry(&Filter{
		EntryPath:         []string{"/"},
		Search:            []Search{{Field: FieldPath, Input: "%foo%"}},
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)
	entries, err = analyzer.Pages.Entry(&Filter{
		From:      util.PastDay(2),
		To:        util.Today(),
		EventName: []string{"event"},
		EventMeta: map[string]string{
			"author": "John",
		},
		CustomMetricKey:  "amount",
		CustomMetricType: pkg.CustomMetricTypeFloat,
		Sample:           10_000_000,
	})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, 1, entries[0].Visitors)
	assert.Equal(t, 1, entries[0].Entries)
	assert.InDelta(t, 0.1111, entries[0].EntryRate, 0.001)
	exits, err := analyzer.Pages.Exit(nil)
	assert.NoError(t, err)
	assert.Len(t, exits, 3)
	assert.Equal(t, "/", exits[0].Path)
	assert.Equal(t, "/bar", exits[1].Path)
	assert.Equal(t, "/foo", exits[2].Path)
	assert.Empty(t, exits[0].Title)
	assert.Empty(t, exits[1].Title)
	assert.Empty(t, exits[2].Title)
	assert.Equal(t, 6, exits[0].Visitors)
	assert.Equal(t, 4, exits[1].Visitors)
	assert.Equal(t, 1, exits[2].Visitors)
	assert.Equal(t, 7, exits[0].Sessions)
	assert.Equal(t, 4, exits[1].Sessions)
	assert.Equal(t, 1, exits[2].Sessions)
	assert.Equal(t, 4, exits[0].Exits)
	assert.Equal(t, 4, exits[1].Exits)
	assert.Equal(t, 1, exits[2].Exits)
	assert.InDelta(t, 0.4444, exits[0].ExitRate, 0.001)
	assert.InDelta(t, 0.4444, exits[1].ExitRate, 0.001)
	assert.InDelta(t, 0.1111, exits[2].ExitRate, 0.001)
	exits, err = analyzer.Pages.Exit(&Filter{PathPattern: []string{"(?i)^/.*$"}})
	assert.NoError(t, err)
	assert.Len(t, exits, 3)
	assert.Equal(t, "/", exits[0].Path)
	assert.Equal(t, "/bar", exits[1].Path)
	assert.Equal(t, "/foo", exits[2].Path)
	assert.Equal(t, 6, exits[0].Visitors)
	assert.Equal(t, 4, exits[1].Visitors)
	assert.Equal(t, 1, exits[2].Visitors)
	assert.Equal(t, 7, exits[0].Sessions)
	assert.Equal(t, 4, exits[1].Sessions)
	assert.Equal(t, 1, exits[2].Sessions)
	assert.Equal(t, 4, exits[0].Exits)
	assert.Equal(t, 4, exits[1].Exits)
	assert.Equal(t, 1, exits[2].Exits)
	exits, err = analyzer.Pages.Exit(&Filter{From: util.PastDay(1), To: util.Today(), IncludeTitle: true})
	assert.NoError(t, err)
	assert.Len(t, exits, 2)
	assert.Equal(t, "/bar", exits[0].Path)
	assert.Equal(t, "/", exits[1].Path)
	assert.Equal(t, "Bar", exits[0].Title)
	assert.Equal(t, "Home", exits[1].Title)
	assert.Equal(t, 4, exits[0].Visitors)
	assert.Equal(t, 3, exits[1].Visitors)
	assert.Equal(t, 4, exits[0].Exits)
	assert.Equal(t, 1, exits[1].Exits)
	assert.InDelta(t, 0.8, exits[0].ExitRate, 0.001)
	assert.InDelta(t, 0.2, exits[1].ExitRate, 0.01)
	exits, err = analyzer.Pages.Exit(&Filter{From: util.PastDay(1), To: util.Today(), ExitPath: []string{"/"}})
	assert.NoError(t, err)
	assert.Len(t, exits, 1)
	assert.Equal(t, "/", exits[0].Path)
	assert.Equal(t, 3, exits[0].Visitors)
	assert.Equal(t, 1, exits[0].Exits)
	assert.InDelta(t, 0.2, exits[0].ExitRate, 0.01)
	exits, err = analyzer.Pages.Exit(&Filter{
		From:      util.PastDay(2),
		To:        util.Today(),
		EventName: []string{"event"},
		EventMeta: map[string]string{
			"author": "John",
		},
		CustomMetricKey:  "amount",
		CustomMetricType: pkg.CustomMetricTypeFloat,
		Sample:           10_000_000,
	})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, 1, entries[0].Visitors)
	assert.Equal(t, 1, entries[0].Entries)
	assert.InDelta(t, 0.1111, entries[0].EntryRate, 0.001)
	_, err = analyzer.Pages.Exit(&Filter{Path: []string{"/bar"}, IncludeTitle: true})
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(getMaxFilter(""))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(getMaxFilter("event"))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(&Filter{Offset: 1, Limit: 10, Sort: []Sort{
		{
			Field:     FieldExitPath,
			Direction: pkg.DirectionASC,
		},
	}, Search: []Search{
		{
			Field: FieldExitPath,
			Input: "/",
		},
	}})
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(&Filter{Search: []Search{{Field: FieldPath, Input: "%foo%"}}, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(&Filter{
		ExitPath:          []string{"/"},
		Search:            []Search{{Field: FieldPath, Input: "%foo%"}},
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)

	// imported statistics
	past3Days := util.PastDay(3).Format(time.DateOnly)
	_, err = dbClient.Exec(fmt.Sprintf(`INSERT INTO "imported_page" (date, path, visitors, views, sessions) VALUES
		('%s', '/bar', 2, 3, 3), ('%s', '/', 1, 1, 1)`, past3Days, past3Days))
	_, err = dbClient.Exec(fmt.Sprintf(`INSERT INTO "imported_entry_page" (date, entry_path, visitors, sessions) VALUES
		('%s', '/bar', 2, 3), ('%s', '/', 1, 1)`, past3Days, past3Days))
	_, err = dbClient.Exec(fmt.Sprintf(`INSERT INTO "imported_exit_page" (date, exit_path, visitors, sessions) VALUES
		('%s', '/bar', 2, 3), ('%s', '/', 1, 1)`, past3Days, past3Days))
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
	entries, err = analyzer.Pages.Entry(&Filter{
		From:          util.PastDay(3),
		To:            util.Today(),
		ImportedUntil: util.PastDay(2),
	})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, "/bar", entries[1].Path)
	assert.Equal(t, 8, entries[0].Entries)
	assert.Equal(t, 4, entries[1].Entries)
	assert.InDelta(t, 0.6153, entries[0].EntryRate, 0.001)
	assert.InDelta(t, 0.3076, entries[1].EntryRate, 0.001)
	exits, err = analyzer.Pages.Exit(&Filter{
		From:          util.PastDay(3),
		To:            util.Today(),
		ImportedUntil: util.PastDay(2),
	})
	assert.NoError(t, err)
	assert.Len(t, exits, 3)
	assert.Equal(t, "/bar", exits[0].Path)
	assert.Equal(t, "/", exits[1].Path)
	assert.Equal(t, "/foo", exits[2].Path)
	assert.Equal(t, 6, exits[0].Exits)
	assert.Equal(t, 5, exits[1].Exits)
	assert.Equal(t, 1, exits[2].Exits)
	assert.InDelta(t, 0.4615, exits[0].ExitRate, 0.001)
	assert.InDelta(t, 0.3846, exits[1].ExitRate, 0.001)
	assert.InDelta(t, 0.0769, exits[2].ExitRate, 0.001)
	entries, err = analyzer.Pages.Entry(&Filter{
		From:              util.PastDay(3),
		To:                util.Today(),
		ImportedUntil:     util.PastDay(2),
		EntryPath:         []string{"/"},
		IncludeTimeOnPage: true,
		IncludeTitle:      true,
	})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, "Home", entries[0].Title)
	assert.Equal(t, 8, entries[0].Entries)
	assert.InDelta(t, 0.6153, entries[0].EntryRate, 0.001)
	assert.Equal(t, 23, entries[0].AverageTimeSpentSeconds)
	exits, err = analyzer.Pages.Exit(&Filter{
		From:          util.PastDay(3),
		To:            util.Today(),
		ImportedUntil: util.PastDay(2),
		ExitPath:      []string{"/bar"},
		IncludeTitle:  true,
	})
	assert.NoError(t, err)
	assert.Len(t, exits, 1)
	assert.Equal(t, "/bar", exits[0].Path)
	assert.Equal(t, "Bar", exits[0].Title)
	assert.Equal(t, 6, exits[0].Exits)
	assert.InDelta(t, 0.4615, exits[0].ExitRate, 0.001)
}

func TestAnalyzer_EntryExitPagesSortVisitors(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.PastDay(2), SessionID: 1, Path: "/"},
		{VisitorID: 1, Time: util.PastDay(2), SessionID: 2, Path: "/foo"},
		{VisitorID: 2, Time: util.PastDay(2), SessionID: 3, Path: "/"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.PastDay(2), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/", PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.PastDay(2), Start: time.Now(), SessionID: 2, EntryPath: "/foo", ExitPath: "/foo", PageViews: 1},
			{Sign: 1, VisitorID: 2, Time: util.PastDay(2), Start: time.Now(), SessionID: 3, EntryPath: "/", ExitPath: "/", PageViews: 1},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	entries, err := analyzer.Pages.Entry(&Filter{Sort: []Sort{{Field: FieldVisitors, Direction: pkg.DirectionDESC}}})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, "/foo", entries[1].Path)
	assert.Equal(t, 2, entries[0].Visitors)
	assert.Equal(t, 1, entries[1].Visitors)
	entries, err = analyzer.Pages.Entry(&Filter{Sort: []Sort{{Field: FieldVisitors, Direction: pkg.DirectionASC}}})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/foo", entries[0].Path)
	assert.Equal(t, "/", entries[1].Path)
	assert.Equal(t, 1, entries[0].Visitors)
	assert.Equal(t, 2, entries[1].Visitors)
	entries, err = analyzer.Pages.Entry(&Filter{Sort: []Sort{{Field: FieldEntryRate, Direction: pkg.DirectionDESC}}})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, "/foo", entries[1].Path)
	assert.Equal(t, 2, entries[0].Visitors)
	assert.Equal(t, 1, entries[1].Visitors)
	entries, err = analyzer.Pages.Entry(&Filter{Sort: []Sort{{Field: FieldEntryRate, Direction: pkg.DirectionASC}}})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "/foo", entries[0].Path)
	assert.Equal(t, "/", entries[1].Path)
	assert.Equal(t, 1, entries[0].Visitors)
	assert.Equal(t, 2, entries[1].Visitors)
	exits, err := analyzer.Pages.Exit(&Filter{Sort: []Sort{{Field: FieldVisitors, Direction: pkg.DirectionDESC}}})
	assert.NoError(t, err)
	assert.Len(t, exits, 2)
	assert.Equal(t, "/", exits[0].Path)
	assert.Equal(t, "/foo", exits[1].Path)
	assert.Equal(t, 2, exits[0].Visitors)
	assert.Equal(t, 1, exits[1].Visitors)
	exits, err = analyzer.Pages.Exit(&Filter{Sort: []Sort{{Field: FieldVisitors, Direction: pkg.DirectionASC}}})
	assert.NoError(t, err)
	assert.Len(t, exits, 2)
	assert.Equal(t, "/foo", exits[0].Path)
	assert.Equal(t, "/", exits[1].Path)
	assert.Equal(t, 1, exits[0].Visitors)
	assert.Equal(t, 2, exits[1].Visitors)
	exits, err = analyzer.Pages.Exit(&Filter{Sort: []Sort{{Field: FieldExitRate, Direction: pkg.DirectionDESC}}})
	assert.NoError(t, err)
	assert.Len(t, exits, 2)
	assert.Equal(t, "/", exits[0].Path)
	assert.Equal(t, "/foo", exits[1].Path)
	assert.Equal(t, 2, exits[0].Visitors)
	assert.Equal(t, 1, exits[1].Visitors)
	exits, err = analyzer.Pages.Exit(&Filter{Sort: []Sort{{Field: FieldExitRate, Direction: pkg.DirectionASC}}})
	assert.NoError(t, err)
	assert.Len(t, exits, 2)
	assert.Equal(t, "/foo", exits[0].Path)
	assert.Equal(t, "/", exits[1].Path)
	assert.Equal(t, 1, exits[0].Visitors)
	assert.Equal(t, 2, exits[1].Visitors)
}

func TestAnalyzer_EntryExitPagesEvents(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), SessionID: 1, Path: "/"},
		{VisitorID: 1, Time: util.Today().Add(time.Second), SessionID: 1, Path: "/", DurationSeconds: 8},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 15), SessionID: 1, Path: "/foo", DurationSeconds: 31},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 20), SessionID: 1, Path: "/", DurationSeconds: 9},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 25), SessionID: 1, Path: "/bar", DurationSeconds: 21},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/", PageViews: 1},
		},
		{
			{Sign: -1, VisitorID: 1, Time: util.Today(), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/", PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/foo", PageViews: 2},
			{Sign: -1, VisitorID: 1, Time: util.Today().Add(time.Second), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/foo", PageViews: 2},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second * 2), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/bar", PageViews: 3},
		},
	})
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{Name: "event", VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/foo"},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	entries, err := analyzer.Pages.Entry(&Filter{EventName: []string{"event"}, IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, 1, entries[0].Visitors)
	assert.Equal(t, 1, entries[0].Sessions)
	assert.Equal(t, 1, entries[0].Entries)
	assert.InDelta(t, 1, entries[0].EntryRate, 0.001)
	assert.Equal(t, 20, entries[0].AverageTimeSpentSeconds)
	exits, err := analyzer.Pages.Exit(&Filter{EventName: []string{"event"}})
	assert.NoError(t, err)
	assert.Len(t, exits, 1)
	assert.Equal(t, "/bar", exits[0].Path)
	assert.Equal(t, 1, exits[0].Visitors)
	assert.Equal(t, 1, exits[0].Sessions)
	assert.Equal(t, 1, exits[0].Exits)
	assert.InDelta(t, 1, exits[0].ExitRate, 0.001)
	_, err = analyzer.Pages.Entry(getMaxFilter("event"))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(getMaxFilter("event"))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Entry(getMaxFilter("!event"))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(getMaxFilter("!event"))
	assert.NoError(t, err)
}

func TestAnalyzer_EntryExitPagesTags(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), SessionID: 1, Path: "/", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, Time: util.Today().Add(time.Second), SessionID: 1, Path: "/", DurationSeconds: 8, TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 15), SessionID: 1, Path: "/foo", DurationSeconds: 31, TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 20), SessionID: 1, Path: "/", DurationSeconds: 9, TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 25), SessionID: 1, Path: "/bar", DurationSeconds: 21, TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/", PageViews: 1},
		},
		{
			{Sign: -1, VisitorID: 1, Time: util.Today(), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/", PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/foo", PageViews: 2},
			{Sign: -1, VisitorID: 1, Time: util.Today().Add(time.Second), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/foo", PageViews: 2},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second * 2), Start: time.Now(), SessionID: 1, EntryPath: "/", ExitPath: "/bar", PageViews: 3},
		},
	})
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{Name: "event", VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/foo", MetaKeys: []string{"author"}, MetaValues: []string{"Alice"}},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	entries, err := analyzer.Pages.Entry(&Filter{
		EventName:         []string{"event"},
		Tags:              map[string]string{"author": "Alice"}, // references the event
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, 1, entries[0].Visitors)
	assert.Equal(t, 1, entries[0].Sessions)
	assert.Equal(t, 1, entries[0].Entries)
	assert.InDelta(t, 1, entries[0].EntryRate, 0.001)
	assert.Equal(t, 0, entries[0].AverageTimeSpentSeconds)
	entries, err = analyzer.Pages.Entry(&Filter{
		Tags:              map[string]string{"author": "John"}, // page view tag
		IncludeTimeOnPage: true,
	})
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "/", entries[0].Path)
	assert.Equal(t, 1, entries[0].Visitors)
	assert.Equal(t, 1, entries[0].Sessions)
	assert.Equal(t, 1, entries[0].Entries)
	assert.InDelta(t, 1, entries[0].EntryRate, 0.001)
	assert.Equal(t, 20, entries[0].AverageTimeSpentSeconds)
	exits, err := analyzer.Pages.Exit(&Filter{
		EventName: []string{"event"},
		Tags:      map[string]string{"author": "Alice"},
	})
	assert.NoError(t, err)
	assert.Len(t, exits, 1)
	assert.Equal(t, "/bar", exits[0].Path)
	assert.Equal(t, 1, exits[0].Visitors)
	assert.Equal(t, 1, exits[0].Sessions)
	assert.Equal(t, 1, exits[0].Exits)
	assert.InDelta(t, 1, exits[0].ExitRate, 0.001)
	exits, err = analyzer.Pages.Exit(&Filter{
		Tags: map[string]string{"author": "Alice"},
	})
	assert.NoError(t, err)
	assert.Len(t, exits, 1)
	assert.Equal(t, "/bar", exits[0].Path)
	assert.Equal(t, 1, exits[0].Visitors)
	assert.Equal(t, 1, exits[0].Sessions)
	assert.Equal(t, 1, exits[0].Exits)
	assert.InDelta(t, 1, exits[0].ExitRate, 0.001)
}

func TestAnalyzer_EntryExitPagePathFilter(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), DurationSeconds: 0, Path: "/"},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 3), DurationSeconds: 3, Path: "/account/billing/"},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 5), DurationSeconds: 2, Path: "/settings/general/"},
		{VisitorID: 1, Time: util.Today().Add(time.Second * 7), DurationSeconds: 2, Path: "/integrations/wordpress/"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second * 7), Start: time.Now(), DurationSeconds: 0, EntryPath: "/", ExitPath: "/", PageViews: 1, IsBounce: true},
		},
		{
			{Sign: -1, VisitorID: 1, Time: util.Today().Add(time.Second * 7), Start: time.Now(), DurationSeconds: 0, EntryPath: "/", ExitPath: "/", PageViews: 1, IsBounce: true},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second * 10), Start: time.Now(), DurationSeconds: 3, EntryPath: "/", ExitPath: "/integrations/wordpress/", PageViews: 2, IsBounce: false},
			{Sign: -1, VisitorID: 1, Time: util.Today().Add(time.Second * 12), Start: time.Now(), DurationSeconds: 3, EntryPath: "/", ExitPath: "/integrations/wordpress/", PageViews: 2, IsBounce: false},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second * 14), Start: time.Now(), DurationSeconds: 5, EntryPath: "/", ExitPath: "/account/billing/", PageViews: 3, IsBounce: false},
			{Sign: -1, VisitorID: 1, Time: util.Today().Add(time.Second * 14), Start: time.Now(), DurationSeconds: 5, EntryPath: "/", ExitPath: "/account/billing/", PageViews: 3, IsBounce: false},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second * 14), Start: time.Now(), DurationSeconds: 7, EntryPath: "/", ExitPath: "/settings/general/", PageViews: 4, IsBounce: false},
			{Sign: -1, VisitorID: 1, Time: util.Today().Add(time.Second * 14), Start: time.Now(), DurationSeconds: 7, EntryPath: "/", ExitPath: "/settings/general/", PageViews: 4, IsBounce: false},
			{Sign: 1, VisitorID: 1, Time: util.Today().Add(time.Second * 14), Start: time.Now(), DurationSeconds: 7, EntryPath: "/", ExitPath: "/integrations/wordpress/", PageViews: 5, IsBounce: false},
		},
	})
	analyzer := NewAnalyzer(dbClient)
	filter := &Filter{
		Path:  []string{"/account/billing/"},
		Limit: 11,
	}
	entry, err := analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entry, 1)
	assert.Equal(t, "/", entry[0].Path)
	assert.Equal(t, 1, entry[0].Visitors)
	assert.Equal(t, 1, entry[0].Entries)
	exit, err := analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exit, 1)
	assert.Equal(t, "/integrations/wordpress/", exit[0].Path)
	assert.Equal(t, 1, exit[0].Visitors)
	assert.Equal(t, 1, exit[0].Exits)

	filter.Path = []string{"/foo"}
	entry, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entry, 0)
	exit, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exit, 0)
}

func TestAnalyzer_EntryExitPageFilterCombination(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		// / -> /foo -> /bar -> /exit
		{VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/"},
		{VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Second * 10), Path: "/foo"},
		{VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Second * 20), Path: "/bar"},
		{VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Second * 30), Path: "/exit"},

		// / -> /bar -> /
		{VisitorID: 2, SessionID: 2, Time: util.Today(), Path: "/"},
		{VisitorID: 2, SessionID: 2, Time: util.Today().Add(time.Second * 10), Path: "/bar"},
		{VisitorID: 2, SessionID: 2, Time: util.Today().Add(time.Second * 20), Path: "/"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Second * 30), Start: time.Now(), ExitPath: "/", EntryPath: "/exit", PageViews: 4, IsBounce: true},
		},
		{
			{Sign: -1, VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Second * 30), Start: time.Now(), ExitPath: "/", EntryPath: "/exit", PageViews: 4, IsBounce: true},
			{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Second * 30), Start: time.Now(), ExitPath: "/exit", EntryPath: "/", PageViews: 4, IsBounce: false},
			{Sign: 1, VisitorID: 2, SessionID: 2, Time: util.Today().Add(time.Second * 20), Start: time.Now(), ExitPath: "/", EntryPath: "/", PageViews: 3, IsBounce: true},
		},
	})
	analyzer := NewAnalyzer(dbClient)

	// no filter
	pages, err := analyzer.Pages.ByPath(nil)
	assert.NoError(t, err)
	assert.Len(t, pages, 4)
	assert.Equal(t, "/", pages[0].Path)
	assert.Equal(t, "/bar", pages[1].Path)
	assert.Equal(t, "/exit", pages[2].Path)
	assert.Equal(t, "/foo", pages[3].Path)
	assert.Equal(t, 2, pages[0].Visitors)
	assert.Equal(t, 2, pages[1].Visitors)
	assert.Equal(t, 1, pages[2].Visitors)
	assert.Equal(t, 1, pages[3].Visitors)
	entryPages, err := analyzer.Pages.Entry(nil)
	assert.NoError(t, err)
	assert.Len(t, entryPages, 1)
	assert.Equal(t, "/", entryPages[0].Path)
	assert.Equal(t, 2, entryPages[0].Visitors)
	assert.Equal(t, 2, entryPages[0].Sessions)
	assert.Equal(t, 2, entryPages[0].Entries)
	exitPages, err := analyzer.Pages.Exit(nil)
	assert.NoError(t, err)
	assert.Len(t, exitPages, 2)
	assert.Equal(t, "/", exitPages[0].Path)
	assert.Equal(t, "/exit", exitPages[1].Path)
	assert.Equal(t, 2, exitPages[0].Visitors)
	assert.Equal(t, 1, exitPages[1].Visitors)
	assert.Equal(t, 1, exitPages[0].Exits)
	assert.Equal(t, 1, exitPages[1].Exits)

	// filter for a path
	filter := &Filter{Path: []string{"/bar"}}
	pages, err = analyzer.Pages.ByPath(filter)
	assert.NoError(t, err)
	assert.Len(t, pages, 1)
	assert.Equal(t, "/bar", pages[0].Path)
	assert.Equal(t, 2, pages[0].Visitors)
	entryPages, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entryPages, 1)
	assert.Equal(t, "/", entryPages[0].Path)
	assert.Equal(t, 2, entryPages[0].Visitors)
	assert.Equal(t, 2, entryPages[0].Entries)
	exitPages, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exitPages, 2)
	assert.Equal(t, "/", exitPages[0].Path)
	assert.Equal(t, "/exit", exitPages[1].Path)
	assert.Equal(t, 2, exitPages[0].Visitors)
	assert.Equal(t, 1, exitPages[1].Visitors)
	assert.Equal(t, 1, exitPages[0].Exits)
	assert.Equal(t, 1, exitPages[1].Exits)

	// filter entry page
	filter.Path = nil
	filter.EntryPath = []string{"/bar"}
	pages, err = analyzer.Pages.ByPath(filter)
	assert.NoError(t, err)
	assert.Len(t, pages, 0)
	entryPages, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entryPages, 0)
	exitPages, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exitPages, 0)

	filter.EntryPath = []string{"/"}
	pages, err = analyzer.Pages.ByPath(filter)
	assert.NoError(t, err)
	assert.Len(t, pages, 4)
	assert.Equal(t, "/", pages[0].Path)
	assert.Equal(t, "/bar", pages[1].Path)
	assert.Equal(t, "/exit", pages[2].Path)
	assert.Equal(t, "/foo", pages[3].Path)
	assert.Equal(t, 2, pages[0].Visitors)
	assert.Equal(t, 2, pages[1].Visitors)
	assert.Equal(t, 1, pages[2].Visitors)
	assert.Equal(t, 1, pages[3].Visitors)
	entryPages, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entryPages, 1)
	assert.Equal(t, "/", entryPages[0].Path)
	assert.Equal(t, 2, entryPages[0].Visitors)
	assert.Equal(t, 2, entryPages[0].Entries)
	exitPages, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exitPages, 2)
	assert.Equal(t, "/", exitPages[0].Path)
	assert.Equal(t, "/exit", exitPages[1].Path)
	assert.Equal(t, 2, exitPages[0].Visitors)
	assert.Equal(t, 1, exitPages[1].Visitors)
	assert.Equal(t, 1, exitPages[0].Exits)
	assert.Equal(t, 1, exitPages[1].Exits)

	// filter entry + exit page
	filter.ExitPath = []string{"/bar"}
	pages, err = analyzer.Pages.ByPath(filter)
	assert.NoError(t, err)
	assert.Len(t, pages, 0)
	entryPages, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entryPages, 0)
	exitPages, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exitPages, 0)

	filter.ExitPath = []string{"/exit"}
	pages, err = analyzer.Pages.ByPath(filter)
	assert.NoError(t, err)
	assert.Len(t, pages, 4)
	assert.Equal(t, "/", pages[0].Path)
	assert.Equal(t, "/bar", pages[1].Path)
	assert.Equal(t, "/exit", pages[2].Path)
	assert.Equal(t, "/foo", pages[3].Path)
	assert.Equal(t, 1, pages[0].Visitors)
	assert.Equal(t, 1, pages[1].Visitors)
	assert.Equal(t, 1, pages[2].Visitors)
	assert.Equal(t, 1, pages[3].Visitors)
	entryPages, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entryPages, 1)
	assert.Equal(t, "/", entryPages[0].Path)
	assert.Equal(t, 2, entryPages[0].Visitors)
	assert.Equal(t, 1, entryPages[0].Entries)
	exitPages, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exitPages, 1)
	assert.Equal(t, "/exit", exitPages[0].Path)
	assert.Equal(t, 1, exitPages[0].Visitors)
	assert.Equal(t, 1, exitPages[0].Exits)

	// filter entry + exit page + page
	filter.Path = []string{"/bar"}
	pages, err = analyzer.Pages.ByPath(filter)
	assert.NoError(t, err)
	assert.Len(t, pages, 1)
	assert.Equal(t, "/bar", pages[0].Path)
	assert.Equal(t, 1, pages[0].Visitors)
	entryPages, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	assert.Len(t, entryPages, 1)
	assert.Equal(t, "/", entryPages[0].Path)
	assert.Equal(t, 2, entryPages[0].Visitors)
	assert.Equal(t, 1, entryPages[0].Entries)
	exitPages, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
	assert.Len(t, exitPages, 1)
	assert.Equal(t, "/exit", exitPages[0].Path)
	assert.Equal(t, 1, exitPages[0].Visitors)
	assert.Equal(t, 1, exitPages[0].Exits)

	// filter conversion goal
	filter = &Filter{PathPattern: []string{"(?i)^/bar$"}}
	_, err = analyzer.Pages.ByPath(filter)
	assert.NoError(t, err)
	_, err = analyzer.Pages.Entry(filter)
	assert.NoError(t, err)
	_, err = analyzer.Pages.Exit(filter)
	assert.NoError(t, err)
}

func TestAnalyzer_EntryExitManyPages(t *testing.T) {
	db.CleanupDB(t, dbClient)
	var sessions []model.Session
	var pageViews []model.PageView

	for i := 0; i < 100_000; i++ {
		sessions = append(sessions, model.Session{
			Sign:            1,
			VisitorID:       uint64(i + 1),
			Time:            util.Today(),
			Start:           time.Now(),
			SessionID:       1,
			DurationSeconds: 60,
			PageViews:       1,
			EntryPath:       fmt.Sprintf("/relatively/long/unique/entry/%d", i),
			ExitPath:        fmt.Sprintf("/relatively/long/unique/exit/%d", i),
		})
		pageViews = append(pageViews, model.PageView{
			VisitorID: uint64(i + 1),
			Time:      util.Today(),
			SessionID: 1,
			Path:      fmt.Sprintf("/relatively/long/unique/path/%d", i),
		})
	}

	saveSessions(t, [][]model.Session{sessions})
	assert.NoError(t, dbClient.SavePageViews(pageViews))
	analyzer := NewAnalyzer(dbClient)
	entries, err := analyzer.Pages.Entry(&Filter{IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, entries, 100_000)
	exits, err := analyzer.Pages.Exit(&Filter{IncludeTimeOnPage: true})
	assert.NoError(t, err)
	assert.Len(t, exits, 100_000)
}

func TestAnalyzer_Conversions(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/"},
		{VisitorID: 2, Time: util.Today(), Path: "/simple/page"},
		{VisitorID: 2, Time: util.Today().Add(time.Minute), Path: "/simple/page"},
		{VisitorID: 3, Time: util.Today(), Path: "/siMple/page/"},
		{VisitorID: 3, Time: util.Today().Add(time.Minute), Path: "/siMple/page/"},
		{VisitorID: 4, Time: util.Today(), Path: "/simple/page/with/many/slashes"},
	}))
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), ExitPath: "/foo", PageViews: 1},
		},
		{
			{Sign: -1, VisitorID: 1, Time: util.Today(), Start: time.Now(), ExitPath: "/foo", PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), ExitPath: "/", PageViews: 2},
			{Sign: 1, VisitorID: 2, Time: util.Today().Add(time.Minute), Start: time.Now(), ExitPath: "/simple/page", PageViews: 1},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), ExitPath: "/siMple/page/", PageViews: 1},
			{Sign: 1, VisitorID: 3, Time: util.Today().Add(time.Minute), Start: time.Now(), ExitPath: "/siMple/page/", PageViews: 2},
			{Sign: 1, VisitorID: 4, Time: util.Today(), Start: time.Now(), ExitPath: "/simple/page/with/many/slashes", PageViews: 1},
		},
	})
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, Time: util.Today(), Name: "Sale", MetaKeys: []string{"amount", "currency"}, MetaValues: []string{"189", "EUR"}, Path: "/simple/page"},
		{VisitorID: 2, Time: util.Today(), Name: "Sale", MetaKeys: []string{"amount", "currency"}, MetaValues: []string{"312", "EUR"}, Path: "/simple"},
		{VisitorID: 4, Time: util.Today(), Name: "Sale", MetaKeys: []string{"amount", "currency"}, MetaValues: []string{"177", "USD"}, Path: "/simple/page"},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	stats, err := analyzer.Pages.Conversions(nil)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 4, stats.Visitors)
	assert.Equal(t, 7, stats.Views)
	assert.InDelta(t, 1, stats.CR, 0.01)
	assert.InDelta(t, 0, stats.CustomMetricAvg, 0.001)
	assert.InDelta(t, 0, stats.CustomMetricTotal, 0.001)
	stats, err = analyzer.Pages.Conversions(&Filter{
		From:        util.Today(),
		To:          util.Today(),
		PathPattern: []string{"(?i)^/.*$"},
		Period:      pkg.PeriodDay,
	})
	assert.NoError(t, err)
	assert.Equal(t, 4, stats.Visitors)
	assert.Equal(t, 6, stats.Views)
	assert.InDelta(t, 1, stats.CR, 0.01)
	assert.InDelta(t, 0, stats.CustomMetricAvg, 0.001)
	assert.InDelta(t, 0, stats.CustomMetricTotal, 0.001)
	stats, err = analyzer.Pages.Conversions(&Filter{PathPattern: []string{"(?i)^/simple/[^/]+/.*"}})
	assert.NoError(t, err)
	assert.Equal(t, 2, stats.Visitors)
	assert.Equal(t, 3, stats.Views)
	assert.InDelta(t, 0.5, stats.CR, 0.01)
	assert.InDelta(t, 0, stats.CustomMetricAvg, 0.001)
	assert.InDelta(t, 0, stats.CustomMetricTotal, 0.001)
	stats, err = analyzer.Pages.Conversions(&Filter{
		EventName:        []string{"Sale"},
		CustomMetricKey:  "amount",
		CustomMetricType: pkg.CustomMetricTypeInteger,
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, stats.Visitors)
	assert.Equal(t, 4, stats.Views)
	assert.InDelta(t, 0.75, stats.CR, 0.01)
	assert.InDelta(t, 226, stats.CustomMetricAvg, 0.001)
	assert.InDelta(t, 678, stats.CustomMetricTotal, 0.001)
	stats, err = analyzer.Pages.Conversions(&Filter{
		EventName:        []string{"Sale"},
		EventMeta:        map[string]string{"currency": "EUR"},
		CustomMetricKey:  "amount",
		CustomMetricType: pkg.CustomMetricTypeInteger,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, stats.Visitors)
	assert.Equal(t, 3, stats.Views)
	assert.InDelta(t, 0.5, stats.CR, 0.01)
	assert.InDelta(t, 250.5, stats.CustomMetricAvg, 0.001)
	assert.InDelta(t, 501, stats.CustomMetricTotal, 0.001)
	stats, err = analyzer.Pages.Conversions(&Filter{
		EventName:        []string{"Sale"},
		EventMeta:        map[string]string{"currency": "EUR"},
		CustomMetricKey:  "amount",
		CustomMetricType: pkg.CustomMetricTypeInteger,
		PathPattern:      []string{"(?i)^/.*"},
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, stats.Visitors)
	assert.Equal(t, 3, stats.Views)
	assert.InDelta(t, 0.5, stats.CR, 0.01)
	assert.InDelta(t, 250.5, stats.CustomMetricAvg, 0.001)
	assert.InDelta(t, 501, stats.CustomMetricTotal, 0.001)
	_, err = analyzer.Pages.Conversions(getMaxFilter(""))
	assert.NoError(t, err)
	_, err = analyzer.Pages.Conversions(getMaxFilter("event"))
	assert.NoError(t, err)
	filter := getMaxFilter("Sale")
	filter.CustomMetricType = pkg.CustomMetricTypeInteger
	filter.CustomMetricKey = "amount"
	filter.From = util.Today()
	filter.To = util.Today()
	_, err = analyzer.Pages.Conversions(filter)
	assert.NoError(t, err)
}

func TestAnalyzer_avgTimeOnPage(t *testing.T) {
	db.CleanupDB(t, dbClient)
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/"},
		},
		{
			{Sign: -1, VisitorID: 1, SessionID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/"},
			{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/bar"},
			{Sign: 1, VisitorID: 2, SessionID: 2, Time: util.Today(), Start: time.Now(), EntryPath: "/bar", ExitPath: "/"},
		},
	})
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Minute * 2), Path: "/foo", DurationSeconds: 120, TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Minute*2 + time.Second*23), Path: "/bar", DurationSeconds: 23, TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 2, SessionID: 2, Time: util.Today(), Path: "/bar", TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 2, SessionID: 2, Time: util.Today().Add(time.Second * 16), Path: "/foo", DurationSeconds: 16, TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 2, SessionID: 2, Time: util.Today().Add(time.Second*16 + time.Second*8), Path: "/", DurationSeconds: 7, TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
	}))
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Minute*2 + time.Second*2), Path: "/foo", Name: "event", MetaKeys: []string{"key", "author"}, MetaValues: []string{"value", "John"}},
		{VisitorID: 1, SessionID: 1, Time: util.Today().Add(time.Minute*2 + time.Second*5), Path: "/foo", Name: "another event", MetaKeys: []string{"author"}, MetaValues: []string{"John"}},
		{VisitorID: 2, SessionID: 2, Time: util.Today().Add(time.Second*16 + time.Second*9), Path: "/", Name: "event", MetaKeys: []string{"author"}, MetaValues: []string{"Alice"}},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	stats, err := analyzer.Pages.avgTimeOnPage(nil, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 3)
	paths := []string{stats[0].Path, stats[1].Path, stats[2].Path}
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")
	top := []int{stats[0].AverageTimeSpentSeconds, stats[1].AverageTimeSpentSeconds, stats[2].AverageTimeSpentSeconds}
	assert.Contains(t, top, 120)
	assert.Contains(t, top, 16)
	assert.Contains(t, top, 15)
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From:      util.PastDay(1),
		To:        util.Today(),
		EventName: []string{"event"},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 3)
	paths = []string{stats[0].Path, stats[1].Path, stats[2].Path}
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")
	top = []int{stats[0].AverageTimeSpentSeconds, stats[1].AverageTimeSpentSeconds, stats[2].AverageTimeSpentSeconds}
	assert.Contains(t, top, 120)
	assert.Contains(t, top, 16)
	assert.Contains(t, top, 15)
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From:      util.PastDay(1),
		To:        util.Today(),
		EventName: []string{"event"},
		EventMeta: map[string]string{
			"key": "value",
		},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	paths = []string{stats[0].Path, stats[1].Path}
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	top = []int{stats[0].AverageTimeSpentSeconds, stats[1].AverageTimeSpentSeconds}
	assert.Contains(t, top, 120)
	assert.Contains(t, top, 23)
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From:      util.PastDay(1),
		To:        util.Today(),
		EntryPath: []string{"/"},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	paths = []string{stats[0].Path, stats[1].Path}
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	top = []int{stats[0].AverageTimeSpentSeconds, stats[1].AverageTimeSpentSeconds}
	assert.Contains(t, top, 120)
	assert.Contains(t, top, 23)
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From: util.PastDay(1),
		To:   util.Today(),
		Tags: map[string]string{"author": "!John"},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	paths = []string{stats[0].Path, stats[1].Path}
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From: util.PastDay(1),
		To:   util.Today(),
		Tags: map[string]string{"author": "John"},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, stats[0].Path, "/")
	assert.Equal(t, stats[1].Path, "/foo")
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From:      util.PastDay(1),
		To:        util.Today(),
		EventName: []string{"event"},
		Tags:      map[string]string{"author": "Alice"},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	paths = []string{stats[0].Path, stats[1].Path}
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From:     util.PastDay(1),
		To:       util.Today(),
		ExitPath: []string{"/", "/bar"},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 3)
	paths = []string{stats[0].Path, stats[1].Path, stats[2].Path}
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")
	top = []int{stats[0].AverageTimeSpentSeconds, stats[1].AverageTimeSpentSeconds, stats[2].AverageTimeSpentSeconds}
	assert.Contains(t, top, 120)
	assert.Contains(t, top, (23+8)/2)
	assert.Contains(t, top, 16)
	stats, err = analyzer.Pages.avgTimeOnPage(&Filter{
		From:      util.PastDay(1),
		To:        util.Today(),
		EntryPath: []string{"/foo"},
		ExitPath:  []string{"/bar"},
	}, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, stats, 0)
	_, err = analyzer.Pages.avgTimeOnPage(getMaxFilter(""), []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	_, err = analyzer.Pages.avgTimeOnPage(getMaxFilter("event"), []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
}

func TestAnalyzer_totalVisitorsSessions(t *testing.T) {
	db.CleanupDB(t, dbClient)
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/foo", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/bar", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/bar", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, SessionID: 2, Time: util.Today(), Path: "/foo", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 2, SessionID: 1, Time: util.Today(), Path: "/", TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 2, SessionID: 2, Time: util.Today(), Path: "/foo", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 3, SessionID: 1, Time: util.Today(), Path: "/", TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 3, SessionID: 1, Time: util.Today(), Path: "/foo", TagKeys: []string{"author"}, TagValues: []string{"John"}},
	}))
	assert.NoError(t, dbClient.SaveSessions([]model.Session{
		{Sign: 1, VisitorID: 1, SessionID: 1, Time: util.Today(), Start: time.Now()},
		{Sign: 1, VisitorID: 1, SessionID: 2, Time: util.Today(), Start: time.Now()},
		{Sign: 1, VisitorID: 2, SessionID: 1, Time: util.Today(), Start: time.Now()},
		{Sign: 1, VisitorID: 2, SessionID: 2, Time: util.Today(), Start: time.Now()},
		{Sign: 1, VisitorID: 3, SessionID: 1, Time: util.Today(), Start: time.Now()},
		{Sign: 1, VisitorID: 3, SessionID: 1, Time: util.Today(), Start: time.Now()},
	}))
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{Name: "event", VisitorID: 1, SessionID: 1, Time: util.Today(), Path: "/foo", MetaKeys: []string{"author"}, MetaValues: []string{"John"}},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	total, err := analyzer.Pages.totalVisitorsSessions(nil, []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	assert.Len(t, total, 3)
	assert.Equal(t, "/foo", total[0].Path)
	assert.Equal(t, "/", total[1].Path)
	assert.Equal(t, "/bar", total[2].Path)
	assert.Equal(t, 4, total[0].Views)
	assert.Equal(t, 3, total[1].Views)
	assert.Equal(t, 2, total[2].Views)
	assert.Equal(t, 3, total[0].Visitors)
	assert.Equal(t, 3, total[1].Visitors)
	assert.Equal(t, 1, total[2].Visitors)
	assert.Equal(t, 4, total[0].Sessions)
	assert.Equal(t, 3, total[1].Sessions)
	assert.Equal(t, 1, total[2].Sessions)
	total, err = analyzer.Pages.totalVisitorsSessions(nil, []string{"/"})
	assert.NoError(t, err)
	assert.Len(t, total, 1)
	assert.Equal(t, "/", total[0].Path)
	assert.Equal(t, 3, total[0].Views)
	assert.Equal(t, 3, total[0].Visitors)
	assert.Equal(t, 3, total[0].Sessions)
	total, err = analyzer.Pages.totalVisitorsSessions(&Filter{EventName: []string{"event"}}, []string{"/"})
	assert.NoError(t, err)
	assert.Len(t, total, 1)
	assert.Equal(t, "/", total[0].Path)
	assert.Equal(t, 1, total[0].Views)
	assert.Equal(t, 1, total[0].Visitors)
	assert.Equal(t, 1, total[0].Sessions)
	total, err = analyzer.Pages.totalVisitorsSessions(&Filter{Tags: map[string]string{"author": "Alice"}}, []string{"/"})
	assert.NoError(t, err)
	assert.Len(t, total, 1)
	assert.Equal(t, "/", total[0].Path)
	assert.Equal(t, 3, total[0].Views)
	assert.Equal(t, 3, total[0].Visitors)
	assert.Equal(t, 3, total[0].Sessions)
	_, err = analyzer.Pages.totalVisitorsSessions(getMaxFilter(""), []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
	_, err = analyzer.Pages.totalVisitorsSessions(getMaxFilter("event"), []string{"/", "/foo", "/bar"})
	assert.NoError(t, err)
}

func TestGetPathList(t *testing.T) {
	paths := getPathList([]model.PageStats{
		{Path: "/"},
		{Path: "/foo"},
		{Path: "/bar"},
	})
	assert.Len(t, paths, 3)
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")

	paths = getPathList([]model.EntryStats{
		{Path: "/"},
		{Path: "/foo"},
		{Path: "/bar"},
	})
	assert.Len(t, paths, 3)
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")

	paths = getPathList([]model.ExitStats{
		{Path: "/"},
		{Path: "/foo"},
		{Path: "/bar"},
	})
	assert.Len(t, paths, 3)
	assert.Contains(t, paths, "/")
	assert.Contains(t, paths, "/foo")
	assert.Contains(t, paths, "/bar")
}
