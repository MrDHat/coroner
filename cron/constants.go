package cron

type SectionName string

var SectionMinute SectionName = "minute"
var SectionHour SectionName = "hour"
var SectionDayOfMonth SectionName = "dayOfMonth"
var SectionMonth SectionName = "month"
var SectionDayOfWeek SectionName = "dayOfWeek"

var Ranges = map[SectionName][]int{
	SectionMinute:     {0, 59},
	SectionHour:       {0, 23},
	SectionDayOfMonth: {1, 31},
	SectionMonth:      {1, 12},
	SectionDayOfWeek:  {1, 7},
}

var PartPositions = []SectionName{
	SectionMinute,
	SectionHour,
	SectionDayOfMonth,
	SectionMonth,
	SectionDayOfWeek,
}

var PartOutput = map[SectionName]string{
	SectionMinute:     "minute",
	SectionHour:       "hour",
	SectionDayOfMonth: "day of month",
	SectionMonth:      "month",
	SectionDayOfWeek:  "day of week",
}