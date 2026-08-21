package domain

type Summary struct {
	Total     int
	Draft     int
	Pending   int
	Confirmed int
	Withdrawn int
	Archived  int
}

func (s *Summary) Add(r Record) {
	s.Total++
	switch r.Status {
	case Draft:
		s.Draft++
	case PendingReview:
		s.Pending++
	case Confirmed:
		s.Confirmed++
	case Withdrawn:
		s.Withdrawn++
	case Archived:
		s.Archived++
	}
}
func (s Summary) CompletionRate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Confirmed+s.Archived) / float64(s.Total)
}
func (s Summary) NeedsAttention() int { return s.Draft + s.Pending + s.Withdrawn }
