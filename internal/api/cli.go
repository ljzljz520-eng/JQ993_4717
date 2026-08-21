package api

import (
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/report"
	"aroma-maintenance/internal/review"
	"aroma-maintenance/internal/search"
	"encoding/json"
	"fmt"
	"io"
)

type CLI struct {
	Catalog *catalog.Service
	Review  *review.Service
	Search  *search.Service
	Report  *report.Service
}

func (c CLI) Run(args []string, in io.Reader, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("command required")
	}
	switch args[0] {
	case "summary":
		sum, err := c.Report.Summary()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, report.RenderSummary(sum))
		return err
	case "list":
		rows, err := c.Search.Search(domain.Filter{})
		if err != nil {
			return err
		}
		return json.NewEncoder(out).Encode(rows)
	default:
		return fmt.Errorf("unknown command %s", args[0])
	}
}
