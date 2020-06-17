package services

import (
	"context"

	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/web/data"
)

// CommentService is the interface definition of comment related functionalities.
type CommentService interface {
	Create(ctx context.Context, request *data.CreateCommentRequest) *errors.Type
	Update(ctx context.Context, request *data.UpdateCommentRequest) *errors.Type
	Delete(ctx context.Context, commentID *data.ID) *errors.Type
}
