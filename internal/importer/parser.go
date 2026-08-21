package importer

import (
	"aroma-maintenance/internal/domain"
	"bufio"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type Parser struct{}

func New() *Parser { return &Parser{} }
func (p *Parser) ParseCSV(r io.Reader) ([]domain.Record, []string, error) {
	cr := csv.NewReader(bufio.NewReader(r))
	cr.TrimLeadingSpace = true
	headers, err := cr.Read()
	if err != nil {
		return nil, nil, err
	}
	index := map[string]int{}
	for i, h := range headers {
		index[strings.ToLower(strings.TrimSpace(h))] = i
	}
	rows := []domain.Record{}
	errors := []string{}
	line := 1
	for {
		line++
		values, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return rows, errors, err
		}
		get := func(name string) string {
			if i, ok := index[name]; ok && i < len(values) {
				return strings.TrimSpace(values[i])
			}
			return ""
		}
		record := domain.NewRecord(get("id"), get("batch"), get("name"), get("scent"), get("material"), get("owner"))
		record.Tags = domain.SanitizeTags(strings.Split(get("tags"), "|"))
		if err := record.Validate(); err != nil {
			errors = append(errors, fmtLine(line, err))
			continue
		}
		rows = append(rows, record)
	}
	return rows, errors, nil
}
func fmtLine(line int, err error) string { return "line " + strconv.Itoa(line) + ": " + err.Error() }
