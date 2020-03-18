package condorcet_test

import (
	"strconv"
	"testing"

	"github.com/batiazinga/condorcet"
)

func TestElection_Winner(t *testing.T) {
	testcases := []struct {
		label      string
		num        int     // number of candidates
		ballots    [][]int // ballots prefixed by the number of times this ballot appears
		hasWinnter bool
		winner     int
	}{
		{
			// example from Condorcet described here
			// https://fr.wikipedia.org/wiki/M%C3%A9thode_de_Condorcet
			label: "Condorcet's example",
			num:   3,
			ballots: [][]int{
				[]int{
					23,
					0, 2, 1,
				},
				[]int{
					19,
					1, 2, 0,
				},
				[]int{
					16,
					2, 1, 0,
				},
				[]int{
					2,
					2, 0, 1,
				},
			},
			hasWinnter: true,
			winner:     2,
		},
		{
			// example from
			// https://en.wikipedia.org/wiki/Condorcet_method
			label: "4 candidates",
			num:   4,
			ballots: [][]int{
				[]int{
					42,
					2, 3, 0, 1,
				},
				[]int{
					26,
					3, 0, 1, 2,
				},
				[]int{
					15,
					0, 1, 3, 2,
				},
				[]int{
					17,
					1, 0, 3, 2,
				},
			},
			hasWinnter: true,
			winner:     3,
		},
		{
			// example from
			// https://fr.wikipedia.org/wiki/Paradoxe_de_Condorcet
			label: "paradoxe",
			num:   3,
			ballots: [][]int{
				[]int{
					23,
					0, 1, 2,
				},
				[]int{
					17,
					1, 2, 0,
				},
				[]int{
					2,
					1, 0, 2,
				},
				[]int{
					10,
					2, 0, 1,
				},
				[]int{
					8,
					2, 1, 0,
				},
			},
			hasWinnter: false,
		},
	}

	for i, tc := range testcases {
		t.Run(
			strconv.Itoa(i),
			func(t *testing.T) {
				e, err := condorcet.New(tc.num)
				if err != nil {
					t.Errorf("invalid testcase: %v", err)
					return
				}

				for _, ballot := range tc.ballots {
					for i := 0; i < ballot[0]; i++ {
						valid := e.Vote(ballot[1:]...)
						if !valid {
							t.Errorf("invalid testcase: invalid ballot %v", ballot[1:])
							return
						}
					}
				}

				w, exist := e.Winner()
				if exist && !tc.hasWinnter {
					t.Error("no winner expected")
					return
				}
				if !exist && tc.hasWinnter {
					t.Error("a winner was expected")
					return
				}
				if exist && w != tc.winner {
					t.Errorf("wrong winner: %d instead of %d", w, tc.winner)
				}
			},
		)
	}
}
