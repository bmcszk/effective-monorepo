package e2e

import "testing"

func Test(t *testing.T) {
	given, when, then := NewBlocks(t)

	given.aGame().and().
		aWhiteOpeningMove()

	when.dispatchingMove()

	then.moveIsDispatched().and().
		fetchingBoard().and().
		boardIsFetched()
}
