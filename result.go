package condorcet

// Result is an immutable snapshot of an election.
//
// A Result must be obtained from an Election.
type Result struct {
	e *Election
}

// Winner returns the winner of the election, if any.
// If there is no winner it returns false.
//
// An election with no vote has no winner.
func (r Result) Winner() (w int, exist bool) {
	// find the winner
	for i := 1; i < r.e.num(); i++ {
		// i is the challenger of w
		if r.e.m[r.e.index(w, i)] < r.e.m[r.e.index(i, w)] {
			w = i // i beats w
		}
	}

	// is w really a winner?
	for i := 0; i < r.e.num(); i++ {
		if w == i {
			continue
		}

		// i is the challenger of w
		if r.e.m[r.e.index(w, i)] <= r.e.m[r.e.index(i, w)] {
			return // w fails to beat i: not a winner finally
		}
	}

	return w, true
}

// NumVoters returns the number of voters.
func (r Result) NumVoters() uint { return r.e.NumVoters() }
