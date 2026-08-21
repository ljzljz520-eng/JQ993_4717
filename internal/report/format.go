package report

import (
	"aroma-maintenance/internal/domain"
	"fmt"
	"strings"
)

func RenderSummary(sum domain.Summary) string {
	return fmt.Sprintf("total=%d draft=%d pending=%d confirmed=%d withdrawn=%d archived=%d", sum.Total, sum.Draft, sum.Pending, sum.Confirmed, sum.Withdrawn, sum.Archived)
}
func RenderRecord(r domain.Record) string {
	return strings.Join([]string{r.ID, r.Batch, r.Name, r.Scent, string(r.Status), r.Owner}, " | ")
}
func RenderQueue(rows []domain.Record) string {
	parts := make([]string, 0, len(rows))
	for _, r := range rows {
		parts = append(parts, RenderRecord(r))
	}
	return strings.Join(parts, "\n")
}
func StatusLabel(status domain.Status) string {
	switch status {
	case domain.Draft:
		return "草稿"
	case domain.PendingReview:
		return "待审核"
	case domain.Confirmed:
		return "已确认"
	case domain.Withdrawn:
		return "已撤回"
	case domain.Archived:
		return "已归档"
	default:
		return "未知"
	}
}
func Progress(sum domain.Summary) string {
	if sum.Total == 0 {
		return "0%"
	}
	return fmt.Sprintf("%.0f%%", sum.CompletionRate()*100)
}
