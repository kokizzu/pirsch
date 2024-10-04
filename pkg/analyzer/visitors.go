package analyzer

import (
	"errors"
	"fmt"
	"github.com/pirsch-analytics/pirsch/v6/pkg/db"
	"github.com/pirsch-analytics/pirsch/v6/pkg/model"
	"github.com/pirsch-analytics/pirsch/v6/pkg/util"
	"strings"
	"time"
)

var (
	// ErrNoPeriodOrDay is returned in case no period or day was specified to calculate the growth rate.
	ErrNoPeriodOrDay = errors.New("no period or day specified")
)

// Visitors aggregates statistics regarding visitors.
type Visitors struct {
	analyzer *Analyzer
	store    db.Store
}

// Active returns the active visitors per hostname, path, and (optional) page title and the total number of active visitors for given duration.
// Use time.Minute*5 for example to get the active visitors for the past 5 minutes.
func (visitors *Visitors) Active(filter *Filter, duration time.Duration) ([]model.ActiveVisitorStats, int, error) {
	filter = visitors.analyzer.getFilter(filter)
	filter.From = time.Now().UTC().Add(-duration)
	filter.IncludeTime = true
	fields := []Field{FieldHostname, FieldPath}
	groupBy := []Field{FieldHostname, FieldPath}
	orderBy := []Field{FieldVisitors, FieldHostname, FieldPath}

	if filter.IncludeTitle {
		fields = append(fields, FieldTitle)
		groupBy = append(groupBy, FieldTitle)
		orderBy = append(orderBy, FieldTitle)
	}

	fields = append(fields, FieldVisitors)
	q, args := filter.buildQuery(fields, groupBy, orderBy, nil, "")
	stats, err := visitors.store.SelectActiveVisitorStats(filter.Ctx, filter.IncludeTitle, q, args...)

	if err != nil {
		return nil, 0, err
	}

	q, args = filter.buildQuery([]Field{FieldVisitors}, nil, nil, nil, "")
	count, err := visitors.store.Count(filter.Ctx, q, args...)

	if err != nil {
		return nil, 0, err
	}

	return stats, count, nil
}

// Total returns the total visitor count, session count, bounce rate, views, CR, and average and total custom metric.
func (visitors *Visitors) Total(filter *Filter) (*model.TotalVisitorStats, error) {
	filter = visitors.analyzer.getFilter(filter)
	fields := []Field{
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
		FieldBounceRate,
	}

	if filter.IncludeCR {
		fields = append(fields, FieldCR)
	}

	includeCustomMetric := false

	if len(filter.EventName) > 0 && filter.CustomMetricType != "" && filter.CustomMetricKey != "" {
		fields = append(fields, FieldEventMetaCustomMetricAvg, FieldEventMetaCustomMetricTotal)
		includeCustomMetric = true
	}

	q, args := filter.buildQuery(fields, nil, nil, []Field{
		FieldVisitors,
		FieldViews,
		FieldSessions,
		FieldBounces,
	}, "imported_visitors")
	stats, err := visitors.store.GetTotalVisitorStats(filter.Ctx, q, filter.IncludeCR, includeCustomMetric, args...)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// TotalVisitors returns the total unique visitor count.
func (visitors *Visitors) TotalVisitors(filter *Filter) (int, error) {
	filter = visitors.analyzer.getFilter(filter)
	f := &Filter{
		ClientID:      filter.ClientID,
		Timezone:      filter.Timezone,
		From:          filter.From,
		To:            filter.To,
		ImportedUntil: filter.ImportedUntil,
		Sample:        filter.Sample,
		IncludeTime:   filter.IncludeTime,
	}
	q, args := f.buildQuery([]Field{FieldVisitors}, nil, nil, []Field{FieldVisitors}, "imported_visitors")
	total, err := visitors.store.GetTotalUniqueVisitorStats(filter.Ctx, q, args...)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// TotalPageViews returns the total number of page views.
func (visitors *Visitors) TotalPageViews(filter *Filter) (int, error) {
	filter = visitors.analyzer.getFilter(filter)
	f := &Filter{
		ClientID:      filter.ClientID,
		Timezone:      filter.Timezone,
		From:          filter.From,
		To:            filter.To,
		ImportedUntil: filter.ImportedUntil,
		Sample:        filter.Sample,
	}
	q, args := f.buildQuery([]Field{FieldViews}, nil, nil, []Field{FieldViews}, "imported_visitors")
	total, err := visitors.store.GetTotalPageViewStats(filter.Ctx, q, args...)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// TotalSessions returns the total number of sessions.
func (visitors *Visitors) TotalSessions(filter *Filter) (int, error) {
	filter = visitors.analyzer.getFilter(filter)
	f := &Filter{
		ClientID:      filter.ClientID,
		Timezone:      filter.Timezone,
		From:          filter.From,
		To:            filter.To,
		ImportedUntil: filter.ImportedUntil,
		Sample:        filter.Sample,
	}
	q, args := f.buildQuery([]Field{FieldSessions}, nil, nil, []Field{FieldSessions}, "imported_visitors")
	total, err := visitors.store.GetTotalSessionStats(filter.Ctx, q, args...)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// TotalVisitorsPageViews returns the total visitor count and number of page views including the growth.
func (visitors *Visitors) TotalVisitorsPageViews(filter *Filter) (*model.TotalVisitorsPageViewsStats, error) {
	filter = visitors.analyzer.getFilter(filter)

	if filter.From.IsZero() || filter.To.IsZero() {
		return nil, ErrNoPeriodOrDay
	}

	q, args := filter.buildQuery([]Field{
		FieldVisitors,
		FieldViews,
	}, nil, nil, []Field{
		FieldVisitors,
		FieldViews,
	}, "imported_visitors")
	current, err := visitors.store.GetTotalVisitorsPageViewsStats(filter.Ctx, q, args...)

	if err != nil {
		return nil, err
	}

	visitors.getPreviousPeriod(filter)
	q, args = filter.buildQuery([]Field{
		FieldVisitors,
		FieldViews,
	}, nil, nil, []Field{
		FieldVisitors,
		FieldViews,
	}, "imported_visitors")
	previous, err := visitors.store.GetTotalVisitorsPageViewsStats(filter.Ctx, q, args...)

	if err != nil {
		return nil, err
	}

	return &model.TotalVisitorsPageViewsStats{
		Visitors:       current.Visitors,
		Views:          current.Views,
		VisitorsGrowth: calculateGrowth(current.Visitors, previous.Visitors),
		ViewsGrowth:    calculateGrowth(current.Views, previous.Views),
	}, nil
}

// ByPeriod returns the visitor count, session count, bounce rate, views, CR, and average and total custom metric
// grouped by day, week, month, or year.
func (visitors *Visitors) ByPeriod(filter *Filter) ([]model.VisitorStats, error) {
	filter = visitors.analyzer.getFilter(filter)
	fields := []Field{
		FieldDay,
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
		FieldBounceRate,
	}

	if filter.IncludeCR {
		fields = append(fields, FieldCRPeriod)
	}

	includeCustomMetric := false

	if len(filter.EventName) > 0 && filter.CustomMetricType != "" && filter.CustomMetricKey != "" {
		fields = append(fields, FieldEventMetaCustomMetricAvg, FieldEventMetaCustomMetricTotal)
		includeCustomMetric = true
	}

	q, args := filter.buildQuery(fields, []Field{
		FieldDay,
	}, []Field{
		FieldDay,
		FieldVisitors,
	}, []Field{
		FieldDay,
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
	}, "imported_visitors")
	stats, err := visitors.store.SelectVisitorStats(filter.Ctx, filter.Period, q, filter.IncludeCR, includeCustomMetric, args...)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// ByHour returns the visitor count grouped by time of day.
func (visitors *Visitors) ByHour(filter *Filter) ([]model.VisitorHourStats, error) {
	filter = visitors.analyzer.getFilter(filter)
	fields := []Field{
		FieldHour,
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
		FieldBounceRate,
	}

	if filter.IncludeCR {
		fields = append(fields, FieldCRPeriod)
	}

	includeCustomMetric := false

	if len(filter.EventName) > 0 && filter.CustomMetricType != "" && filter.CustomMetricKey != "" {
		fields = append(fields, FieldEventMetaCustomMetricAvg, FieldEventMetaCustomMetricTotal)
		includeCustomMetric = true
	}

	q, args := filter.buildQuery(fields, []Field{
		FieldHour,
	}, []Field{
		FieldHour,
		FieldVisitors,
	}, []Field{
		FieldHour,
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
	}, "imported_visitors")
	stats, err := visitors.store.SelectVisitorHourStats(filter.Ctx, q, filter.IncludeCR, includeCustomMetric, args...)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// ByMinute returns the visitor count grouped by the minute of the hour.
func (visitors *Visitors) ByMinute(filter *Filter) ([]model.VisitorMinuteStats, error) {
	filter = visitors.analyzer.getFilter(filter)
	fields := []Field{
		FieldMinute,
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
		FieldBounceRate,
	}

	if filter.IncludeCR {
		fields = append(fields, FieldCRPeriod)
	}

	includeCustomMetric := false

	if len(filter.EventName) > 0 && filter.CustomMetricType != "" && filter.CustomMetricKey != "" {
		fields = append(fields, FieldEventMetaCustomMetricAvg, FieldEventMetaCustomMetricTotal)
		includeCustomMetric = true
	}

	q, args := filter.buildQuery(fields, []Field{
		FieldMinute,
	}, []Field{
		FieldMinute,
		FieldVisitors,
	}, nil, "")
	stats, err := visitors.store.SelectVisitorMinuteStats(filter.Ctx, q, filter.IncludeCR, includeCustomMetric, args...)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// ByWeekdayAndHour returns the visitor count grouped by time of day and weekday.
func (visitors *Visitors) ByWeekdayAndHour(filter *Filter) ([]model.VisitorWeekdayHourStats, error) {
	filter = visitors.analyzer.getFilter(filter)
	q, args := filter.buildQuery([]Field{
		FieldWeekday,
		FieldHour,
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
	}, []Field{
		FieldWeekday,
		FieldHour,
	}, []Field{
		FieldWeekday,
		FieldHour,
	}, nil, "")
	stats, err := visitors.store.SelectVisitorWeekdayHourStats(filter.Ctx, q, args...)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// Growth returns the growth rate for visitor count, session count, bounces, views, and average session duration or average time on page (if path is set).
// The growth rate is relative to the previous time range or day.
// The period or day for the filter must be set, else an error is returned.
func (visitors *Visitors) Growth(filter *Filter) (*model.Growth, error) {
	filter = visitors.analyzer.getFilter(filter)

	if filter.From.IsZero() || filter.To.IsZero() {
		return nil, ErrNoPeriodOrDay
	}

	// get current statistics
	fields := []Field{
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
		FieldBounceRate,
	}
	fieldsImported := []Field{
		FieldVisitors,
		FieldSessions,
		FieldViews,
		FieldBounces,
	}

	if filter.IncludeCR {
		fields = append(fields, FieldCR)
	}

	includeCustomMetric := false

	if len(filter.EventName) > 0 && filter.CustomMetricType != "" && filter.CustomMetricKey != "" {
		fields = append(fields, FieldEventMetaCustomMetricAvg, FieldEventMetaCustomMetricTotal)
		includeCustomMetric = true
	}

	q, args := filter.buildQuery(fields, nil, nil, fieldsImported, "imported_visitors")
	current, err := visitors.store.GetGrowthStats(filter.Ctx, q, filter.IncludeCR, includeCustomMetric, args...)

	if err != nil {
		return nil, err
	}

	var currentTimeSpent int

	if len(filter.EventName) > 0 {
		currentTimeSpent, err = visitors.totalEventDuration(filter)
	} else if len(filter.Path) == 0 {
		currentTimeSpent, err = visitors.totalSessionDuration(filter)
	} else {
		currentTimeSpent, err = visitors.totalTimeOnPage(filter)
	}

	if err != nil {
		return nil, err
	}

	visitors.getPreviousPeriod(filter)
	q, args = filter.buildQuery(fields, nil, nil, fieldsImported, "imported_visitors")
	previous, err := visitors.store.GetGrowthStats(filter.Ctx, q, filter.IncludeCR, includeCustomMetric, args...)

	if err != nil {
		return nil, err
	}

	var previousTimeSpent int

	if len(filter.EventName) > 0 {
		previousTimeSpent, err = visitors.totalEventDuration(filter)
	} else if len(filter.Path) == 0 {
		previousTimeSpent, err = visitors.totalSessionDuration(filter)
	} else {
		previousTimeSpent, err = visitors.totalTimeOnPage(filter)
	}

	if err != nil {
		return nil, err
	}

	return &model.Growth{
		VisitorsGrowth:          calculateGrowth(current.Visitors, previous.Visitors),
		ViewsGrowth:             calculateGrowth(current.Views, previous.Views),
		SessionsGrowth:          calculateGrowth(current.Sessions, previous.Sessions),
		BouncesGrowth:           calculateGrowth(current.BounceRate, previous.BounceRate),
		TimeSpentGrowth:         calculateGrowth(currentTimeSpent, previousTimeSpent),
		CRGrowth:                calculateGrowth(current.CR, previous.CR),
		CustomMetricAvgGrowth:   calculateGrowth(current.CustomMetricAvg, previous.CustomMetricAvg),
		CustomMetricTotalGrowth: calculateGrowth(current.CustomMetricTotal, previous.CustomMetricTotal),
	}, nil
}

// Referrer returns the visitor count and bounce rate grouped by referrer.
func (visitors *Visitors) Referrer(filter *Filter) ([]model.ReferrerStats, error) {
	filter = visitors.analyzer.getFilter(filter)
	var fields, groupBy, orderBy, importedFields []Field

	if filter.ImportedUntil.IsZero() {
		fields = []Field{
			FieldReferrerName,
			FieldReferrerIcon,
			FieldVisitors,
			FieldSessions,
			FieldRelativeVisitors,
			FieldBounces,
			FieldBounceRate,
		}
		groupBy = []Field{
			FieldReferrerName,
		}
		orderBy = []Field{
			FieldVisitors,
			FieldReferrerName,
		}

		if len(filter.Referrer) > 0 || len(filter.ReferrerName) > 0 {
			fields = append(fields, FieldReferrer)
			groupBy = append(groupBy, FieldReferrer)
			orderBy = append(orderBy, FieldReferrer)
		} else {
			fields = append(fields, FieldAnyReferrer)
		}
	} else {
		groupBy = []Field{
			FieldReferrerName,
		}
		orderBy = []Field{
			FieldVisitors,
			FieldReferrerName,
		}

		if len(filter.Referrer) > 0 || len(filter.ReferrerName) > 0 || filter.searchContains(FieldReferrer) {
			fields = []Field{
				FieldReferrerNameImported,
				FieldAnyReferrerIcon,
				FieldVisitors,
				FieldSessions,
				FieldRelativeVisitors,
				FieldBounces,
				FieldBounceRate,
				FieldAnyReferrerImported,
			}
			groupBy = append(groupBy, FieldAnyReferrerImported)
			orderBy = append(orderBy, FieldAnyReferrerImported)
			importedFields = []Field{FieldReferrerName}
		} else {
			fields = []Field{
				FieldReferrerName,
				FieldAnyReferrerIcon,
				FieldVisitors,
				FieldSessions,
				FieldRelativeVisitors,
				FieldBounces,
				FieldBounceRate,
				FieldAnyReferrer,
			}
			importedFields = []Field{FieldReferrer}
		}

		importedFields = append(importedFields, []Field{
			FieldVisitors,
			FieldSessions,
			FieldBounces,
		}...)
	}

	q, args := filter.buildQuery(fields, groupBy, orderBy, importedFields, "imported_referrer")
	stats, err := visitors.store.SelectReferrerStats(filter.Ctx, q, args...)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (visitors *Visitors) getPreviousPeriod(filter *Filter) {
	from := filter.From
	to := filter.To

	if !filter.importedFrom.IsZero() && filter.importedFrom.Before(from) {
		from = filter.importedFrom
	}

	if from.Equal(to) {
		if to.Equal(util.Today()) {
			from = from.Add(-time.Hour * 24 * 7)
			to = time.Now().UTC().Add(-time.Hour * 24 * 7)
			filter.IncludeTime = true
		} else {
			from = from.Add(-time.Hour * 24 * 7)
			to = to.Add(-time.Hour * 24 * 7)
		}
	} else {
		days := to.Sub(from)

		if days >= time.Hour*24 {
			to = from.Add(-time.Hour * 24)
			from = to.Add(-days)
		} else {
			from = from.Add(-time.Hour * 24)
			to = to.Add(-time.Hour * 24)
		}
	}

	filter.From = from
	filter.To = to

	if !filter.ImportedUntil.IsZero() {
		filter.validate()
	}
}

func (visitors *Visitors) totalSessionDuration(filter *Filter) (int, error) {
	q := queryBuilder{
		filter: filter,
		from:   sessions,
		search: filter.Search,
	}
	var query strings.Builder
	query.WriteString(`SELECT sum(duration_seconds)
		FROM (
			SELECT sum(duration_seconds*sign) duration_seconds
			FROM session t `)

	if len(filter.Path) > 0 || len(filter.PathPattern) > 0 || len(filter.Tag) > 0 || len(filter.Tags) > 0 {
		q.from = pageViews
		whereTime := q.whereTime()
		q.whereFields()
		query.WriteString(fmt.Sprintf(`INNER JOIN (
			SELECT visitor_id, session_id FROM page_view %s %s
		) v
		ON v.visitor_id = t.visitor_id AND v.session_id = t.session_id `, whereTime, q.q.String()))
		q.from = sessions
		q.q.Reset()
		q.where = nil
	}

	query.WriteString(q.whereTime())
	q.whereFields()
	where := q.q.String()

	if where != "" {
		query.WriteString(where)
	}

	query.WriteString("GROUP BY visitor_id, session_id HAVING sum(sign) > 0)")
	averageTimeSpentSeconds, err := visitors.store.Count(filter.Ctx, query.String(), q.args...)

	if err != nil {
		return 0, err
	}

	if !filter.importedFrom.IsZero() && !filter.importedTo.IsZero() {
		q = queryBuilder{
			filter: filter,
			from:   sessions,
			search: filter.Search,
		}
		query.Reset()
		query.WriteString("SELECT sum(session_duration) FROM imported_visitors ")
		query.WriteString(q.whereTimeImported())
		averageTimeSpentSecondsImported, err := visitors.store.Count(filter.Ctx, query.String(), q.args...)

		if err != nil {
			return 0, err
		}

		averageTimeSpentSeconds += averageTimeSpentSecondsImported
	}

	return averageTimeSpentSeconds, nil
}

func (visitors *Visitors) totalEventDuration(filter *Filter) (int, error) {
	q := queryBuilder{
		filter: filter,
		fields: []Field{FieldEventDurationSeconds},
		from:   events,
		search: filter.Search,
		sample: filter.Sample,
	}
	query, args := q.query()
	averageTimeSpentSeconds, err := visitors.store.Count(filter.Ctx, query, args...)

	if err != nil {
		return 0, err
	}

	return averageTimeSpentSeconds, nil
}

func (visitors *Visitors) totalTimeOnPage(filter *Filter) (int, error) {
	filterCopy := *filter
	filterCopy.Sort = nil
	q := queryBuilder{
		filter: &filterCopy,
		from:   pageViews,
		search: filter.Search,
	}
	fields := q.getFields()
	fieldsQuery := strings.Join(fields, ",")

	if len(fields) > 0 {
		fieldsQuery = "," + fieldsQuery
	}

	var query strings.Builder
	query.WriteString(fmt.Sprintf(`SELECT sum(time_on_page) average_time_spent_seconds
		FROM (
			SELECT nth_value(%s, 2) OVER (PARTITION BY v.visitor_id, v.session_id ORDER BY v."time" ASC Rows BETWEEN CURRENT ROW AND 1 FOLLOWING) AS time_on_page
				%s
			FROM page_view v `, visitors.analyzer.timeOnPageQuery(filter), fieldsQuery))

	if len(filter.EntryPath) > 0 || len(filter.ExitPath) > 0 {
		q.from = sessions
		query.WriteString(fmt.Sprintf(`INNER JOIN (
			SELECT visitor_id,
			session_id,
			entry_path,
			exit_path
			FROM session
			WHERE %s
			GROUP BY visitor_id, session_id, entry_path, exit_path
			HAVING sum(sign) > 0
		) s
		ON v.visitor_id = s.visitor_id AND v.session_id = s.session_id `, q.whereTime()[len("WHERE "):]))
		q.from = pageViews
	}

	query.WriteString(fmt.Sprintf(`WHERE %s)
		WHERE time_on_page > 0
		%s`, q.whereTime()[len("WHERE "):], q.q.String()))
	q.whereFields()
	averageTimeSpentSeconds, err := visitors.store.Count(filter.Ctx, query.String(), q.args...)

	if err != nil {
		return 0, err
	}

	return averageTimeSpentSeconds, nil
}

func calculateGrowth[T int | float64](current, previous T) float64 {
	if current == 0 && previous == 0 {
		return 0
	} else if previous == 0 {
		return 1
	}

	c := float64(current)
	p := float64(previous)
	return (c - p) / p
}
