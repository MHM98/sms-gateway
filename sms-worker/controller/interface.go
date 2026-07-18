package controller

import "context"

type IMessageRepository interface {
	MarkMessageStatuSsubmitted(ctx context.Context, messageId uint64) error
}


