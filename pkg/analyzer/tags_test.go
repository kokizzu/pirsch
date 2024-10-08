package analyzer

import (
	"github.com/pirsch-analytics/pirsch/v6/pkg"
	"github.com/pirsch-analytics/pirsch/v6/pkg/db"
	"github.com/pirsch-analytics/pirsch/v6/pkg/model"
	"github.com/pirsch-analytics/pirsch/v6/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTags_Tags(t *testing.T) {
	db.CleanupDB(t, dbClient)
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/foo", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/bar", IsBounce: false, PageViews: 3},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), EntryPath: "/foo", ExitPath: "/foo", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", IsBounce: true, PageViews: 1},
		},
	})
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, Time: util.Today(), Path: "/foo", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, Time: util.Today(), Path: "/bar", TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 2, Time: util.Today(), Path: "/foo", TagKeys: []string{"author", "type"}, TagValues: []string{"John", "blog_post"}},
		{VisitorID: 3, Time: util.Today(), Path: "/", TagKeys: []string{"author", "type"}, TagValues: []string{"Alice", "blog_post"}},
	}))
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, Time: util.Today(), Path: "/", Name: "event", MetaKeys: []string{"key", "author", "amount"}, MetaValues: []string{"value", "John", "99.99"}},
		{VisitorID: 3, Time: util.Today(), Path: "/", Name: "event", MetaKeys: []string{"author", "type"}, MetaValues: []string{"Alice", "blog_post"}},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	stats, err := analyzer.Tags.Keys(nil)
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 3, stats[0].Visitors)
	assert.Equal(t, 2, stats[1].Visitors)
	assert.Equal(t, 5, stats[0].Views)
	assert.Equal(t, 2, stats[1].Views)
	assert.InDelta(t, 1, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.6666, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 1, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.4, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Path: []string{"/"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 2, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 2, stats[0].Views)
	assert.Equal(t, 1, stats[1].Views)
	assert.InDelta(t, 0.6666, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.4, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.2, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From:     util.Today(),
		To:       util.Today(),
		ExitPath: []string{"/foo"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 1, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 1, stats[0].Views)
	assert.Equal(t, 1, stats[1].Views)
	assert.InDelta(t, 0.3333, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.2, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.2, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Tag:  []string{"author"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 3, stats[0].Visitors)
	assert.Equal(t, 2, stats[1].Visitors)
	assert.Equal(t, 5, stats[0].Views)
	assert.Equal(t, 2, stats[1].Views)
	assert.InDelta(t, 1, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.6666, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 1, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.4, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Tag:  []string{"!author"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 0)
	stats, err = analyzer.Tags.Keys(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Tags: map[string]string{"author": "John"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 2, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 3, stats[0].Views)
	assert.Equal(t, 1, stats[1].Views)
	assert.InDelta(t, 0.6666, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.6, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.2, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Tags: map[string]string{"author": "!John"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 2, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 2, stats[0].Views)
	assert.Equal(t, 1, stats[1].Views)
	assert.InDelta(t, 0.6666, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.4, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.2, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Tags: map[string]string{"author": "Alice"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 2, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 2, stats[0].Views)
	assert.Equal(t, 1, stats[1].Views)
	assert.InDelta(t, 0.6666, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.4, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.2, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Tags: map[string]string{"author": "~"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 3, stats[0].Visitors)
	assert.Equal(t, 2, stats[1].Visitors)
	assert.Equal(t, 5, stats[0].Views)
	assert.Equal(t, 2, stats[1].Views)
	assert.InDelta(t, 1, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.6666, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 1, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.4, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From:      util.Today(),
		To:        util.Today(),
		EventName: []string{"event"},
		EventMeta: map[string]string{"key": "value"},
		Tags:      map[string]string{"author": "John"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 1)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, 1, stats[0].Visitors)
	assert.Equal(t, 2, stats[0].Views)
	assert.InDelta(t, 0.3333, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.4, stats[0].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From:      util.Today(),
		To:        util.Today(),
		EventName: []string{"event"},
		EventMeta: map[string]string{"key": "value"},
		Tags:      map[string]string{"author": "!John"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 0)
	stats, err = analyzer.Tags.Keys(&Filter{
		From:      util.Today(),
		To:        util.Today(),
		EventName: []string{"event"},
		Tags:      map[string]string{"author": "!John"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, "type", stats[1].Key)
	assert.Equal(t, 1, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 1, stats[0].Views)
	assert.Equal(t, 1, stats[1].Views)
	assert.InDelta(t, 0.3333, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.2, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.2, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Keys(&Filter{
		From:      util.Today(),
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
	assert.Len(t, stats, 1)
	assert.Equal(t, "author", stats[0].Key)
	assert.Equal(t, 1, stats[0].Visitors)
	assert.Equal(t, 3, stats[0].Views)
	assert.InDelta(t, 0.3333, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.6, stats[0].RelativeViews, 0.001)
}

func TestTags_Breakdown(t *testing.T) {
	db.CleanupDB(t, dbClient)
	saveSessions(t, [][]model.Session{
		{
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/foo", IsBounce: false, PageViews: 2},
			{Sign: 1, VisitorID: 1, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/bar", IsBounce: false, PageViews: 3},
			{Sign: 1, VisitorID: 2, Time: util.Today(), Start: time.Now(), EntryPath: "/foo", ExitPath: "/foo", IsBounce: true, PageViews: 1},
			{Sign: 1, VisitorID: 3, Time: util.Today(), Start: time.Now(), EntryPath: "/", ExitPath: "/", IsBounce: true, PageViews: 1},
		},
	})
	assert.NoError(t, dbClient.SavePageViews([]model.PageView{
		{VisitorID: 1, Time: util.Today(), Path: "/", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, Time: util.Today(), Path: "/foo", TagKeys: []string{"author"}, TagValues: []string{"John"}},
		{VisitorID: 1, Time: util.Today(), Path: "/bar", TagKeys: []string{"author"}, TagValues: []string{"Alice"}},
		{VisitorID: 2, Time: util.Today(), Path: "/foo", TagKeys: []string{"author", "type"}, TagValues: []string{"John", "blog_post"}},
		{VisitorID: 3, Time: util.Today(), Path: "/", TagKeys: []string{"author", "type"}, TagValues: []string{"Alice", "blog_post"}},
	}))
	assert.NoError(t, dbClient.SaveEvents([]model.Event{
		{VisitorID: 1, Time: util.Today(), Path: "/", Name: "event", MetaKeys: []string{"key", "author"}, MetaValues: []string{"value", "John"}},
		{VisitorID: 3, Time: util.Today(), Path: "/", Name: "event", MetaKeys: []string{"author", "type"}, MetaValues: []string{"Alice", "blog_post"}},
	}))
	time.Sleep(time.Millisecond * 100)
	analyzer := NewAnalyzer(dbClient)
	stats, err := analyzer.Tags.Breakdown(nil)
	assert.NoError(t, err)
	assert.Empty(t, stats)
	stats, err = analyzer.Tags.Breakdown(&Filter{
		From: util.Today(),
		To:   util.Today(),
		Tag:  []string{"author"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "Alice", stats[0].Value)
	assert.Equal(t, "John", stats[1].Value)
	assert.Equal(t, 2, stats[0].Visitors)
	assert.Equal(t, 2, stats[1].Visitors)
	assert.Equal(t, 2, stats[0].Views)
	assert.Equal(t, 3, stats[1].Views)
	assert.InDelta(t, 0.6666, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.6666, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.4, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.6, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Breakdown(&Filter{
		From:      util.Today(),
		To:        util.Today(),
		EventName: []string{"event"},
		Tag:       []string{"author"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "Alice", stats[0].Value)
	assert.Equal(t, "John", stats[1].Value)
	assert.Equal(t, 2, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 2, stats[0].Views)
	assert.Equal(t, 2, stats[1].Views)
	assert.InDelta(t, 0.6666, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.4, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.4, stats[1].RelativeViews, 0.001)
	stats, err = analyzer.Tags.Breakdown(&Filter{
		From:      util.Today(),
		To:        util.Today(),
		EventName: []string{"event"},
		EventMeta: map[string]string{"key": "value"},
		Tag:       []string{"author"},
	})
	assert.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, "Alice", stats[0].Value)
	assert.Equal(t, "John", stats[1].Value)
	assert.Equal(t, 1, stats[0].Visitors)
	assert.Equal(t, 1, stats[1].Visitors)
	assert.Equal(t, 1, stats[0].Views)
	assert.Equal(t, 2, stats[1].Views)
	assert.InDelta(t, 0.3333, stats[0].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.3333, stats[1].RelativeVisitors, 0.001)
	assert.InDelta(t, 0.2, stats[0].RelativeViews, 0.001)
	assert.InDelta(t, 0.4, stats[1].RelativeViews, 0.001)
}
