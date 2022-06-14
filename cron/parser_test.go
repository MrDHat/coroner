package cron

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite

	mockExpr        string
	mockSection     string
	mockSectionName SectionName
	p               Parser
}

func (t *ParserTestSuite) TestParseSectionHappyPath() {
	// t.mockExpr = "* * * * *"

	t.p = NewParser()

	t.mockSectionName = SectionMinute
	t.mockSection = "?"

	res := t.p.ParseSection(t.mockSection, t.mockSectionName)

	t.Equal(1, len(res))
	t.Equal([]string{"no specific meaning"}, res)
}

func (t *ParserTestSuite) TestParseSectionWeekDayString() {
	// t.mockExpr = "* * * * *"

	t.p = NewParser()

	t.mockSectionName = SectionDayOfWeek
	t.mockSection = "Mon-Fri"

	res := t.p.ParseSection(t.mockSection, t.mockSectionName)

	t.Equal(5, len(res))
	t.Equal([]string{"1", "2", "3", "4", "5"}, res)
}

func (t *ParserTestSuite) TestParseSectionWeekDayMultipleRangeString() {
	// t.mockExpr = "* * * * *"

	t.p = NewParser()

	t.mockSectionName = SectionDayOfWeek
	t.mockSection = "Fri-Tue"

	res := t.p.ParseSection(t.mockSection, t.mockSectionName)

	t.Equal(5, len(res))
	t.Equal([]string{"5", "6", "7", "1", "2"}, res)
}

func TestParserTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, &ParserTestSuite{})
}
