package domain

type Withdrawal struct {
	Record     Record
	Reason     string
	RetryToken string
}

func PrepareWithdrawal(r Record, reason, token string) (Withdrawal, error) {
	if r.IsTerminal() || reason == "" || token == "" {
		return Withdrawal{}, ErrTransition
	}
	next := r.Clone()
	if err := next.Transition(Withdrawn); err != nil {
		return Withdrawal{}, err
	}
	return Withdrawal{Record: next, Reason: reason, RetryToken: token}, nil
}
func ConfirmRetry(w Withdrawal) (Record, error) {
	if w.Record.Status != Withdrawn {
		return Record{}, ErrTransition
	}
	next := w.Record.Clone()
	if err := next.Transition(PendingReview); err != nil {
		return Record{}, err
	}
	return next, nil
}
func IndependentRetry(r Record, token string) (Record, error) {
	w, err := PrepareWithdrawal(r, "retry", token)
	if err != nil {
		return Record{}, err
	}
	return ConfirmRetry(w)
}
