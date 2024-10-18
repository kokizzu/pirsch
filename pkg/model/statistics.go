package model

import (
	"github.com/emvi/null"
)

// ActiveVisitorStats is the result type for active visitor statistics.
type ActiveVisitorStats struct {
	Path     string `json:"path"`
	Title    string `json:"title"`
	Visitors int    `json:"visitors"`
}

// TotalVisitorStats is the result type for total visitor statistics.
type TotalVisitorStats struct {
	Visitors          int     `json:"visitors"`
	Views             int     `json:"views"`
	Sessions          int     `json:"sessions"`
	Bounces           int     `json:"bounces"`
	BounceRate        float64 `db:"bounce_rate" json:"bounce_rate"`
	CR                float64 `json:"cr"`
	CustomMetricAvg   float64 `db:"custom_metric_avg" json:"custom_metric_avg"`
	CustomMetricTotal float64 `db:"custom_metric_total" json:"custom_metric_total"`
}

// TotalVisitorsPageViewsStats is the result type for total visitor count and number of page views statistics.
type TotalVisitorsPageViewsStats struct {
	Visitors       int     `json:"visitors"`
	Views          int     `json:"views"`
	VisitorsGrowth float64 `json:"visitors_growth"`
	ViewsGrowth    float64 `json:"views_growth"`
}

// VisitorStats is the result type for visitor statistics.
type VisitorStats struct {
	Day               null.Time `json:"day"`
	Week              null.Time `json:"week"`
	Month             null.Time `json:"month"`
	Year              null.Time `json:"year"`
	Minute            null.Time `json:"minute"`
	Visitors          int       `json:"visitors"`
	Views             int       `json:"views"`
	Sessions          int       `json:"sessions"`
	Bounces           int       `json:"bounces"`
	BounceRate        float64   `db:"bounce_rate" json:"bounce_rate"`
	CR                float64   `json:"cr"`
	CustomMetricAvg   float64   `db:"custom_metric_avg" json:"custom_metric_avg"`
	CustomMetricTotal float64   `db:"custom_metric_total" json:"custom_metric_total"`
}

// Growth represents the visitors, views, sessions, bounces, and average session duration growth between two time periods.
type Growth struct {
	VisitorsGrowth          float64 `json:"visitors_growth"`
	ViewsGrowth             float64 `json:"views_growth"`
	SessionsGrowth          float64 `json:"sessions_growth"`
	BouncesGrowth           float64 `json:"bounces_growth"`
	TimeSpentGrowth         float64 `json:"time_spent_growth"`
	CRGrowth                float64 `json:"cr_growth"`
	CustomMetricAvgGrowth   float64 `json:"custom_metric_avg_growth"`
	CustomMetricTotalGrowth float64 `json:"custom_metric_total_growth"`
}

// VisitorHourStats is the result type for visitor statistics grouped by time of day.
type VisitorHourStats struct {
	Hour              int     `json:"hour"`
	Visitors          int     `json:"visitors"`
	Views             int     `json:"views"`
	Sessions          int     `json:"sessions"`
	Bounces           int     `json:"bounces"`
	BounceRate        float64 `db:"bounce_rate" json:"bounce_rate"`
	CR                float64 `json:"cr"`
	CustomMetricAvg   float64 `db:"custom_metric_avg" json:"custom_metric_avg"`
	CustomMetricTotal float64 `db:"custom_metric_total" json:"custom_metric_total"`
}

// VisitorMinuteStats is the result type for visitor statistics grouped by the minute of the hour.
type VisitorMinuteStats struct {
	Minute            int     `json:"minute"`
	Visitors          int     `json:"visitors"`
	Views             int     `json:"views"`
	Sessions          int     `json:"sessions"`
	Bounces           int     `json:"bounces"`
	BounceRate        float64 `db:"bounce_rate" json:"bounce_rate"`
	CR                float64 `json:"cr"`
	CustomMetricAvg   float64 `db:"custom_metric_avg" json:"custom_metric_avg"`
	CustomMetricTotal float64 `db:"custom_metric_total" json:"custom_metric_total"`
}

// VisitorWeekdayHourStats is the result type for visitor statistics grouped by time of day and weekday.
type VisitorWeekdayHourStats struct {
	Weekday  int `json:"weekday"`
	Hour     int `json:"hour"`
	Visitors int `json:"visitors"`
	Views    int `json:"views"`
	Sessions int `json:"sessions"`
	Bounces  int `json:"bounces"`
}

// HostnameStats is the result type for hostname statistics.
type HostnameStats struct {
	Hostname         string  `json:"hostname"`
	Visitors         int     `json:"visitors"`
	Views            int     `json:"views"`
	Sessions         int     `json:"sessions"`
	Bounces          int     `json:"bounces"`
	RelativeVisitors float64 `db:"relative_visitors" json:"relative_visitors"`
	RelativeViews    float64 `db:"relative_views" json:"relative_views"`
	BounceRate       float64 `db:"bounce_rate" json:"bounce_rate"`
}

// PageStats is the result type for page statistics.
type PageStats struct {
	Path                    string  `json:"path"`
	Title                   string  `json:"title"`
	Visitors                int     `json:"visitors"`
	Views                   int     `json:"views"`
	Sessions                int     `json:"sessions"`
	Bounces                 int     `json:"bounces"`
	RelativeVisitors        float64 `db:"relative_visitors" json:"relative_visitors"`
	RelativeViews           float64 `db:"relative_views" json:"relative_views"`
	BounceRate              float64 `db:"bounce_rate" json:"bounce_rate"`
	AverageTimeSpentSeconds int     `db:"average_time_spent_seconds" json:"average_time_spent_seconds"`
}

func (stats PageStats) GetPath() string {
	return stats.Path
}

// EntryStats is the result type for entry page statistics.
type EntryStats struct {
	Path                    string  `db:"entry_path" json:"path"`
	Title                   string  `json:"title"`
	Visitors                int     `json:"visitors"`
	Sessions                int     `json:"sessions"`
	Entries                 int     `json:"entries"`
	EntryRate               float64 `db:"entry_rate" json:"entry_rate"`
	AverageTimeSpentSeconds int     `db:"average_time_spent_seconds" json:"average_time_spent_seconds"`
}

func (stats EntryStats) GetPath() string {
	return stats.Path
}

// ExitStats is the result type for exit page statistics.
type ExitStats struct {
	Path     string  `db:"exit_path" json:"path"`
	Title    string  `json:"title"`
	Visitors int     `json:"visitors"`
	Sessions int     `json:"sessions"`
	Exits    int     `json:"exits"`
	ExitRate float64 `db:"exit_rate" json:"exit_rate"`
}

func (stats ExitStats) GetPath() string {
	return stats.Path
}

// ConversionsStats is the result type for page conversions.
type ConversionsStats struct {
	Visitors          int     `json:"visitors"`
	Views             int     `json:"views"`
	CR                float64 `json:"cr"`
	CustomMetricAvg   float64 `db:"custom_metric_avg" json:"custom_metric_avg"`
	CustomMetricTotal float64 `db:"custom_metric_total" json:"custom_metric_total"`
}

// EventStats is the result type for custom events.
type EventStats struct {
	Name                   string   `db:"event_name" json:"name"`
	Count                  int      `json:"count"`
	Visitors               int      `json:"visitors"`
	Views                  int      `json:"views"`
	CR                     float64  `json:"cr"`
	AverageDurationSeconds int      `db:"average_time_spent_seconds" json:"average_duration_seconds"`
	MetaKeys               []string `db:"meta_keys" json:"meta_keys"`
	MetaValue              string   `db:"meta_value" json:"meta_value"`
}

// EventListStats is the result type for a custom event list.
type EventListStats struct {
	Name     string            `db:"event_name" json:"name"`
	Meta     map[string]string `json:"meta"`
	Visitors int               `json:"visitors"`
	Count    int               `json:"count"`
}

// ReferrerStats is the result type for referrer statistics.
type ReferrerStats struct {
	Referrer         string  `json:"referrer"`
	ReferrerName     string  `db:"referrer_name" json:"referrer_name"`
	ReferrerIcon     string  `db:"referrer_icon" json:"referrer_icon"`
	Visitors         int     `json:"visitors"`
	Sessions         int     `json:"sessions"`
	RelativeVisitors float64 `db:"relative_visitors" json:"relative_visitors"`
	Bounces          int     `json:"bounces"`
	BounceRate       float64 `db:"bounce_rate" json:"bounce_rate"`
}

// PlatformStats is the result type for platform statistics.
type PlatformStats struct {
	PlatformDesktop         int     `db:"platform_desktop" json:"platform_desktop"`
	PlatformMobile          int     `db:"platform_mobile" json:"platform_mobile"`
	PlatformUnknown         int     `db:"platform_unknown" json:"platform_unknown"`
	RelativePlatformDesktop float64 `db:"relative_platform_desktop" json:"relative_platform_desktop"`
	RelativePlatformMobile  float64 `db:"relative_platform_mobile" json:"relative_platform_mobile"`
	RelativePlatformUnknown float64 `db:"relative_platform_unknown" json:"relative_platform_unknown"`
}

// TimeSpentStats is the result type for average time spent statistics (sessions, time on page).
type TimeSpentStats struct {
	Day                     null.Time `json:"day"`
	Week                    null.Time `json:"week"`
	Month                   null.Time `json:"month"`
	Year                    null.Time `json:"year"`
	Minute                  null.Time `json:"minute"`
	Path                    string    `json:"path"`
	Title                   string    `json:"title"`
	AverageTimeSpentSeconds int       `db:"average_time_spent_seconds" json:"average_time_spent_seconds"`
}

// MetaStats is the base for meta result types (languages, countries, ...).
type MetaStats struct {
	Visitors         int     `json:"visitors"`
	RelativeVisitors float64 `db:"relative_visitors" json:"relative_visitors"`
}

// LanguageStats is the result type for language statistics.
type LanguageStats struct {
	MetaStats
	Language string `json:"language"`
}

// CountryStats is the result type for country statistics.
type CountryStats struct {
	MetaStats
	CountryCode string `db:"country_code" json:"country_code"`
}

// RegionStats is the result type for region statistics.
type RegionStats struct {
	MetaStats
	CountryCode string `db:"country_code" json:"country_code"`
	Region      string `json:"region"`
}

// CityStats is the result type for city statistics.
type CityStats struct {
	MetaStats
	CountryCode string `db:"country_code" json:"country_code"`
	Region      string `json:"region"`
	City        string `json:"city"`
}

// BrowserStats is the result type for browser statistics.
type BrowserStats struct {
	MetaStats
	Browser string `json:"browser"`
}

// BrowserVersionStats is the result type for browser version statistics.
type BrowserVersionStats struct {
	MetaStats
	Browser        string `json:"browser"`
	BrowserVersion string `db:"browser_version" json:"browser_version"`
}

// OSStats is the result type for operating system statistics.
type OSStats struct {
	MetaStats
	OS string `json:"os"`
}

// OSVersionStats is the result type for operating system version statistics.
type OSVersionStats struct {
	MetaStats
	OS        string `json:"os"`
	OSVersion string `db:"os_version" json:"os_version"`
}

// ScreenClassStats is the result type for screen class statistics.
type ScreenClassStats struct {
	MetaStats
	ScreenClass string `db:"screen_class" json:"screen_class"`
}

// UTMSourceStats is the result type for utm source statistics.
type UTMSourceStats struct {
	MetaStats
	UTMSource string `db:"utm_source" json:"utm_source"`
}

// UTMMediumStats is the result type for utm medium statistics.
type UTMMediumStats struct {
	MetaStats
	UTMMedium string `db:"utm_medium" json:"utm_medium"`
}

// UTMCampaignStats is the result type for utm campaign statistics.
type UTMCampaignStats struct {
	MetaStats
	UTMCampaign string `db:"utm_campaign" json:"utm_campaign"`
}

// UTMContentStats is the result type for utm content statistics.
type UTMContentStats struct {
	MetaStats
	UTMContent string `db:"utm_content" json:"utm_content"`
}

// UTMTermStats is the result type for utm term statistics.
type UTMTermStats struct {
	MetaStats
	UTMTerm string `db:"utm_term" json:"utm_term"`
}

// GrowthStats is the sum to calculate the growth rate.
type GrowthStats struct {
	Visitors          int
	Views             int
	Sessions          int
	Bounces           int
	BounceRate        float64 `db:"bounce_rate"`
	CR                float64
	CustomMetricAvg   float64 `db:"custom_metric_avg" json:"custom_metric_avg"`
	CustomMetricTotal float64 `db:"custom_metric_total" json:"custom_metric_total"`
}

// TotalVisitorSessionStats are the total amount of visitors, views, and sessions for a page.
type TotalVisitorSessionStats struct {
	Path     string
	Visitors int
	Views    int
	Sessions int
}

// AvgTimeSpentStats is the average time spent on a page.
type AvgTimeSpentStats struct {
	Path                    string
	AverageTimeSpentSeconds int `db:"average_time_spent_seconds"`
}

// TagStats is the result type for tags.
type TagStats struct {
	Key              string  `json:"key"`
	Value            string  `json:"value"`
	Visitors         int     `json:"visitors"`
	Views            int     `json:"views"`
	RelativeVisitors float64 `db:"relative_visitors" json:"relative_visitors"`
	RelativeViews    float64 `db:"relative_views" json:"relative_views"`
}

// SessionStep is the result type combining page views and events for a single session.
type SessionStep struct {
	PageView *PageView `json:"page_view"`
	Event    *Event    `json:"event"`
}

// FunnelStep is the result type for a funnel step.
type FunnelStep struct {
	Step                     int     `json:"step"`
	Visitors                 int     `json:"visitors"`
	RelativeVisitors         float64 `db:"relative_visitors" json:"relative_visitors"`
	PreviousVisitors         int     `db:"previous_visitors" json:"previous_visitors"`
	RelativePreviousVisitors float64 `db:"relative_previous_visitors" json:"relative_previous_visitors"`
	Dropped                  int     `json:"dropped"`
	DropOff                  float64 `db:"drop_off" json:"drop_off"`
}
