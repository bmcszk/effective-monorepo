//go:build e2e

package e2e_test

import (
	"testing"

	"github.com/bmcszk/effective-monorepo/e2e"
)

func Test(t *testing.T) {
	given, when, then := e2e.NewBlocks(t)

	given.AGame().And().
		AWhiteOpeningMove()

	when.DispatchingMove()

	then.MoveIsDispatched().And().
		FetchingBoard().And().
		BoardIsFetched()
}
