package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(expr string)
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

func (p *parser) parseSection(section string, sectionName SectionName) ([]string, error) {
	var (
		min = Ranges[sectionName][0]
		max = Ranges[sectionName][1]
		res = make([]string, 0)
	)

	if section == "?" {
		res = append(res, "no specific meaning")
		return res, nil
	}

	num, err := strconv.Atoi(section)
	if err == nil {
		// number type data
		if num < min || num > max {
			return nil, fmt.Errorf("%d is out of range (%d-%d) for %v", num, min, max, sectionName)
		} else {
			res = append(res, section)
			return res, nil
		}
	}

	if section == "*" {
		for i := min; i <= max; i++ {
			res = append(res, strconv.Itoa(i))
		}
		return res, nil
	}

	part := strings.TrimPrefix(section, "*/")
	if part != section {
		interval, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		for i := min; i <= max; i += interval {
			res = append(res, strconv.Itoa(i))
		}
		return res, nil
	}

	if strings.Contains(section, ",") {
		parts := strings.Split(section, ",")
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			} else if num < min || num > max {
				return nil, fmt.Errorf("%d is out of range (%d-%d) for %v", num, min, max, sectionName)
			}
			res = append(res, part)
		}
		return res, nil
	}

	if strings.Contains(section, "-") {
		parts := strings.Split(section, "-")
		if len(parts) != 2 {
			return nil, errors.New("invalid range")
		}

		p1, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		} else if p1 < min {
			e := fmt.Sprintf("%d must be less than %d", p1, min)
			return nil, errors.New(e)
		}

		p2, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		} else if p2 > max {
			e := fmt.Sprintf("%d must be greater than %d", p1, min)
			return nil, errors.New(e)
		}

		if p1 > p2 {
			return nil, errors.New("invalid input")
		}

		for i := p1; i <= p2; i++ {
			res = append(res, strconv.Itoa(i))
		}

		return res, nil
	}

	// special case of "L"
	if strings.Contains(section, "L") {
		num, err := strconv.ParseInt(string(section[0]), 10, 64)
		if err != nil {
			return nil, err
		}
		fmt.Println(num)
		if sectionName == SectionDayOfMonth {
			date := MonthNumberOfDays[int(num)]
			res = append(res, strconv.Itoa(date))
			return res, nil
		}
	}

	// special case of "W"
	if strings.Contains(section, "W") {
		if sectionName == SectionDayOfWeek {
			res = append(res, strconv.Itoa(int(max)))
			return res, nil
		}
	}

	return nil, errors.New("invalid section input")
}

func (p *parser) Parse(expr string) {
	p.Expr = expr

	exprParts, err := p.validate()
	if err != nil {
		panic(err)
	}
	cmd := exprParts[5]

	exprParts = exprParts[:5]
	for i := range exprParts {
		section := PartPositions[i]
		sectionStr := exprParts[i]

		r, err := p.parseSection(sectionStr, section)
		if err != nil {
			fmt.Printf("%v  %v\n", PartOutput[section], err)
		} else {
			fmt.Printf("%v  %v\n", PartOutput[section], r)
		}
	}

	fmt.Println("command  ", cmd)
}

func NewParser() Parser {
	return &parser{}
}
