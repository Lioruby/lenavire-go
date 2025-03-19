package ports

import "time"

type DateProvider interface {
	Now() time.Time
}
