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

var MonthNumberOfDays = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

var DayOfWeekName = map[int]string{
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
	7: "Sunday",
}
