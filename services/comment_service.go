package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
	"github.com/jibitters/kiosk/web/data"
	nc "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// CommentService is a service implementation of comment related functionalities.
type CommentService struct {
	logger            *zap.SugaredLogger
	natsClient        *nc.Conn
	commentRepository *models.CommentRepository
	stop              chan struct{}
}

// NewCommentService returns a newly created and ready to use CommentService.
func NewCommentService(logger *zap.SugaredLogger, db *pgxpool.Pool, natsClient *nc.Conn) *CommentService {
	return &CommentService{
		logger:            logger,
		natsClient:        natsClient,
		commentRepository: models.NewCommentRepository(logger, db),
		stop:              make(chan struct{}),
	}
}

// Start starts the subscription so ready to be notified.
func (s *CommentService) Start() error {
	createCommentSubscription, e := s.natsClient.QueueSubscribe("kiosk.comments.create",
		"kiosk.comments.create_group", s.createComment)
	if e != nil {
		return e
	}

	loadCommentSubscription, e := s.natsClient.QueueSubscribe("kiosk.comments.load",
		"kiosk.comments.load_group", s.loadComment)
	if e != nil {
		return e
	}

	updateCommentSubscription, e := s.natsClient.QueueSubscribe("kiosk.comments.update",
		"kiosk.comments.update_group", s.updateComment)
	if e != nil {
		return e
	}

	deleteCommentSubscription, e := s.natsClient.QueueSubscribe("kiosk.comments.delete",
		"kiosk.comments.delete_group", s.deleteComment)
	if e != nil {
		return e
	}

	go s.await(createCommentSubscription, loadCommentSubscription, updateCommentSubscription, deleteCommentSubscription)

	return nil
}

func (s *CommentService) await(ss ...*nc.Subscription) {
	<-s.stop
	s.logger.Debug("CommentService: received stop signal!")

	for _, s := range ss {
		_ = s.Unsubscribe()
	}
}

func (s *CommentService) createComment(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createCommentRequest := &data.CreateCommentRequest{}
	if e := json.Unmarshal(msg.Data, createCommentRequest); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := createCommentRequest.Validate(); e != nil {
		s.reply(msg, e)
		return
	}

	if e := s.commentRepository.Insert(ctx, *createCommentRequest.AsComment()); e != nil {
		s.reply(msg, e)
		return
	}

	s.replyNoContent(msg)
}

func (s *CommentService) loadComment(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := &data.ID{}
	if e := json.Unmarshal(msg.Data, id); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	c, e := s.commentRepository.LoadByID(ctx, id.ID)
	if e != nil {
		s.reply(msg, e)
		return
	}

	commentResponse := &data.CommentResponse{}
	commentResponse.LoadFromComment(c)
	s.reply(msg, commentResponse)
}

func (s *CommentService) updateComment(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateCommentRequest := &data.UpdateCommentRequest{}
	if e := json.Unmarshal(msg.Data, updateCommentRequest); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := updateCommentRequest.Validate(); e != nil {
		s.reply(msg, e)
		return
	}

	if e := s.commentRepository.Update(ctx, updateCommentRequest.AsComment()); e != nil {
		s.reply(msg, e)
		return
	}

	s.replyNoContent(msg)
}

func (s *CommentService) deleteComment(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := &data.ID{}
	if e := json.Unmarshal(msg.Data, id); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := s.commentRepository.DeleteByID(ctx, id.ID); e != nil {
		s.reply(msg, e)
		return
	}

	s.replyNoContent(msg)
}

func (s *CommentService) reply(msg *nc.Msg, t interface{}) {
	reply, _ := json.Marshal(t)
	_ = msg.Respond(reply)
}

func (s *CommentService) replyNoContent(msg *nc.Msg) {
	_ = msg.Respond([]byte(""))
}

// Stop stops the component and it subscription.
func (s *CommentService) Stop() {
	s.stop <- struct{}{}
}
