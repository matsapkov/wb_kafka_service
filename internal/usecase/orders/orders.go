package orders

import (
	"context"

	"github.com/google/uuid"
)

type Usecase interface {
	GetOrder(ctx context.Context, Id uuid.UUID)
}
