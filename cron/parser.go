package cron

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(expr string)
	ParseSection(section string, sectionName SectionName) []string
}

type parser struct {
	Expr string
}

func (p *parser) validate() ([]string, error) {
	exprParts := strings.Split(p.Expr, " ")
	if len(exprParts) != 6 {
		err := errors.New("invalid cron expression")
		return nil, err
	}
	for i := range exprParts {
		exprParts[i] = strings.TrimSpace(exprParts[i])
		if exprParts[i] == "" {
			err := errors.New("invalid cron expression")
			return nil, err
		}
	}

	return exprParts, nil
}

func (p *parser) Atoi(section string, sectionName SectionName) (int, error) {
	num, err := strconv.Atoi(section)
	if err != nil {
		// check if this is a week day string
		if sectionName == SectionDayOfWeek {
			val, ok := NameOfTheWeekDay[section]
			if !ok {
				err = errors.New("invalid data")
				return 0, err
			}
			return val, nil
		}
	}
	return num, nil
}

func (p *parser) ParseSection(section string, sectionName SectionName) []string {
	var (
		min = Ranges[sectionName][0]
		max = Ranges[sectionName][1]
		res = make([]string, 0)
	)

	if section == "?" {
		res = append(res, "no specific meaning")
		return res
	}

	num, err := p.Atoi(section, sectionName)
	if err == nil {
		// number type data
		if num < min || num > max {
			res = append(res, fmt.Sprintf("%d is out of range (%d-%d) for %v", num, min, max, sectionName))
			return res
		} else {
			res = append(res, fmt.Sprintf("%d", num))
			return res
		}
	}

	if section == "*" {
		for i := min; i <= max; i++ {
			res = append(res, strconv.Itoa(i))
		}
		return res
	}

	part := strings.TrimPrefix(section, "*/")
	if part != section {
		interval, err := p.Atoi(part, sectionName)
		if err != nil {
			res = append(res, err.Error())
			return res
		}
		for i := min; i <= max; i += interval {
			res = append(res, strconv.Itoa(i))
		}
		return res
	}

	if strings.Contains(section, ",") {
		parts := strings.Split(section, ",")
		for _, part := range parts {
			num, err := p.Atoi(part, sectionName)
			if err != nil {
				res = append(res, err.Error())
				return res
			} else if num < min || num > max {
				res = append(res, fmt.Sprintf("%d is out of range (%d-%d) for %v", num, min, max, sectionName))
				return res
			}
			res = append(res, part)
		}
		return res
	}

	if strings.Contains(section, "-") {
		parts := strings.Split(section, "-")
		if len(parts) != 2 {
			res = append(res, fmt.Sprintf("invalid range %v", section))
			return res
		}

		p1, err := p.Atoi(parts[0], sectionName)
		if err != nil {
			res = append(res, err.Error())
			return res
		} else if p1 < min {
			res = append(res, fmt.Sprintf("%d must be less than %d", p1, min))
			return res
		}

		p2, err := p.Atoi(parts[1], sectionName)
		if err != nil {
			res = append(res, err.Error())
			return res
		} else if p2 > max {
			res = append(res, fmt.Sprintf("%d must be greater than %d", p1, min))
			return res
		}

		end := p2
		isMultipleRange := false
		if p1 > p2 {
			end = MaxDayOfWeek
			isMultipleRange = true
		}

		for i := p1; i <= end; i++ {
			res = append(res, strconv.Itoa(i))
		}
		if isMultipleRange {
			for i := MinDayOfWeek; i <= p2; i++ {
				res = append(res, strconv.Itoa(i))
			}
		}

		return res
	}

	// special case of "L"
	if strings.Contains(section, "L") {
		num, err := strconv.ParseInt(string(section[0]), 10, 64)
		if err != nil {
			res = append(res, err.Error())
			return res
		}
		if sectionName == SectionDayOfMonth {
			date := MonthNumberOfDays[int(num)]
			res = append(res, strconv.Itoa(date))
			return res
		}
	}

	// special case of "W"
	if strings.Contains(section, "W") {
		if sectionName == SectionDayOfWeek {
			res = append(res, strconv.Itoa(int(max)))
			return res
		}
	}

	res = append(res, fmt.Sprintf("invalid section %v", section))
	return res
}

func (p *parser) prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (p *parser) Parse(expr string) {
	p.Expr = expr
	exprParts, err := p.validate()
	if err != nil {
		panic(err)
	}

	type result struct {
		Minute     []string
		Hour       []string
		DayOfMonth []string
		Month      []string
		DayOfWeek  []string
		Command    string
	}

	r := result{
		Minute:     p.ParseSection(exprParts[0], SectionMinute),
		Hour:       p.ParseSection(exprParts[1], SectionHour),
		DayOfMonth: p.ParseSection(exprParts[2], SectionDayOfMonth),
		Month:      p.ParseSection(exprParts[3], SectionMonth),
		DayOfWeek:  p.ParseSection(exprParts[4], SectionDayOfWeek),
		Command:    exprParts[5],
	}

	fmt.Printf(
		"minute\t\t%s"+
			"\nhour\t\t%s"+
			"\nday of month\t%s"+
			"\nmonth\t\t%s"+
			"\nday of week\t%s"+
			"\ncommand\t\t%s",
		strings.Join(r.Minute, " "),
		strings.Join(r.Hour, " "),
		strings.Join(r.DayOfMonth, " "),
		strings.Join(r.Month, " "),
		strings.Join(r.DayOfWeek, " "),
		r.Command,
	)
}

func NewParser() Parser {
	return &parser{}
}
