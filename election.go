package condorcet

import "errors"

// Election follows the Condorcet method (see https://en.wikipedia.org/wiki/Condorcet_method).
//
// The (pointer to) default zero value is an election with 2 candidates.
type Election struct {
	n int    // number of candidates - 2
	m []uint // sum matrix (row major order)
}

// New returns an election with n candidates.
// There must be at least 2 candidates.
//
// Candidates are identified by an index such that 0 <= index < n.
func New(n int) (*Election, error) {
	if n < 2 {
		return nil, errors.New("expecting at least 2 candidates")
	}

	return &Election{n: n - 2}, nil
}

// num returns the number of candidates.
func (e *Election) num() int { return e.n + 2 }

// is the sum matrix initialized?
func (e *Election) initialized() bool { return e.m != nil }

// init the sum matrix
// it is an n*n matrix with no value on the diagonal
func (e *Election) init() {
	n := e.num()
	e.m = make([]uint, n*n)
}

// index of the (i,j) pair in the sum matrix
// the sum matrix is stored in row major order
// no check is done on the values of i and j:
//  - i!=j
//  - 0 <= i,j < n
func (e *Election) index(i, j int) int { return e.num()*i + j }

// Vote registers the ballot.
// First item is the prefered candidate, second is the second choice, and so on.
//
// The ballot must be a total order preference over all the candidates.
// Otherwise the ballot is ignored and false is returned.
func (e *Election) Vote(ballot ...int) bool {
	// check that ballot is a total preference
	if len(ballot) != e.num() {
		return false
	}
	candidates := make([]int, e.num())
	for _, candidate := range ballot {
		if candidate < 0 || candidate >= e.num() {
			return false
		}
		candidates[candidate]++
	}
	for _, count := range candidates {
		if count != 1 {
			return false
		}
	}

	if !e.initialized() {
		e.init()
	}

	// fill the sum matrix
	for i := range ballot {
		for j := i + 1; j < len(ballot); j++ {
			// candidate i is prefered to candidate j
			e.m[e.index(ballot[i], ballot[j])]++
		}
	}

	return true
}

// Result returns the a snapshot of the election.
// The election can continue receiving votes without
// impacting previously created results.
func (e *Election) Result() Result {
	if !e.initialized() {
		e.init()
	}

	// copy the content of the election into the result
	cp := &Election{}
	cp.n = e.n
	cp.m = make([]uint, len(e.m))
	copy(cp.m, e.m)

	return Result{cp}
}
