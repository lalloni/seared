package location

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationString(t *testing.T) {
	l := &Location{10, 30, 40}
	assert.EqualValues(t, "position 40 (line 10, column 30)", l.String())
}

func TestLocationStringShort(t *testing.T) {
	l := &Location{10, 30, 40}
	assert.EqualValues(t, "40/10:30", l.ShortString())
}
