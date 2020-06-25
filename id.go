package jobs

import (
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid"
)

func NewID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	ulid := ulid.MustNew(ulid.Timestamp(t), entropy)
	return strings.ToLower(ulid.String())
}
