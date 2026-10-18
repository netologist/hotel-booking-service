package domain

import "context"

type AdminRepository interface {
	TruncateAll(ctx context.Context) error
}
