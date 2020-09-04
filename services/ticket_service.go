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

// TicketService is a service implementation of ticket related functionalities.
type TicketService struct {
	logger           *zap.SugaredLogger
	ticketRepository *models.TicketRepository
	natsClient       *nc.Conn
	stop             chan struct{}
}

// NewTicketService returns a newly created and ready to use TicketService.
func NewTicketService(logger *zap.SugaredLogger, db *pgxpool.Pool, natsClient *nc.Conn) *TicketService {
	return &TicketService{
		logger:           logger,
		ticketRepository: models.NewTicketRepository(logger, db),
		natsClient:       natsClient,
		stop:             make(chan struct{}),
	}
}

// Start starts the subscriptions so ready to be notified.
func (s *TicketService) Start() error {
	createTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.create",
		"kiosk.tickets.create_group", s.create)
	if e != nil {
		return e
	}

	loadTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.load",
		"kiosk.tickets.load_group", s.load)
	if e != nil {
		return e
	}

	updateTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.update",
		"kiosk.tickets.update_group", s.update)
	if e != nil {
		return e
	}

	deleteTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.delete",
		"kiosk.tickets.delete_group", s.delete)
	if e != nil {
		return e
	}

	filterTicketsSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.filter",
		"kiosk.tickets.filter_group", s.filter)
	if e != nil {
		return e
	}

	go s.await(createTicketSubscription, loadTicketSubscription, updateTicketSubscription, deleteTicketSubscription,
		filterTicketsSubscription)

	return nil
}

func (s *TicketService) await(ss ...*nc.Subscription) {
	<-s.stop
	s.logger.Debug("TicketService: received stop signal!")

	for _, s := range ss {
		_ = s.Unsubscribe()
	}
}

func (s *TicketService) create(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createTicketRequest := &data.CreateTicketRequest{}
	if e := json.Unmarshal(msg.Data, createTicketRequest); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := createTicketRequest.Validate(); e != nil {
		s.reply(msg, e)
		return
	}

	if e := s.ticketRepository.Insert(ctx, *createTicketRequest.AsTicket()); e != nil {
		s.reply(msg, e)
		return
	}

	s.replyNoContent(msg)
}

func (s *TicketService) load(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := &data.ID{}
	if e := json.Unmarshal(msg.Data, id); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	t, e := s.ticketRepository.LoadByID(ctx, id.ID)
	if e != nil {
		s.reply(msg, e)
		return
	}

	ticketResponse := &data.TicketResponse{}
	ticketResponse.LoadFromTicket(t)
	s.reply(msg, ticketResponse)
}

func (s *TicketService) update(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateTicketRequest := &data.UpdateTicketRequest{}
	if e := json.Unmarshal(msg.Data, updateTicketRequest); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := updateTicketRequest.Validate(); e != nil {
		s.reply(msg, e)
		return
	}

	if e := s.ticketRepository.Update(ctx, updateTicketRequest.AsTicket()); e != nil {
		s.reply(msg, e)
		return
	}

	s.replyNoContent(msg)
}

func (s *TicketService) delete(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := &data.ID{}
	if e := json.Unmarshal(msg.Data, id); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := s.ticketRepository.DeleteByID(ctx, id.ID); e != nil {
		s.reply(msg, e)
		return
	}

	s.replyNoContent(msg)
}

func (s *TicketService) filter(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filterTicketsRequest := &data.FilterTicketsRequest{}
	if e := json.Unmarshal(msg.Data, filterTicketsRequest); e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := filterTicketsRequest.Validate(); e != nil {
		s.reply(msg, e)
		return
	}

	ts, hasNextPage, e := s.ticketRepository.Filter(ctx, filterTicketsRequest.Issuer, filterTicketsRequest.Owner,
		filterTicketsRequest.ImportanceLevel, filterTicketsRequest.Status, filterTicketsRequest.FromDate,
		filterTicketsRequest.ToDate, filterTicketsRequest.PageNumber, filterTicketsRequest.PageSize)
	if e != nil {
		s.reply(msg, e)
		return
	}

	filterTicketsResponse := &data.FilterTicketsResponse{}
	filterTicketsResponse.LoadFromTickets(ts, hasNextPage)
	s.reply(msg, filterTicketsResponse)
}

func (s *TicketService) reply(msg *nc.Msg, t interface{}) {
	reply, _ := json.Marshal(t)
	_ = msg.Respond(reply)
}

func (s *TicketService) replyNoContent(msg *nc.Msg) {
	_ = msg.Respond([]byte(""))
}

// Stop stops the component and it subscriptions.
func (s *TicketService) Stop() {
	s.stop <- struct{}{}
}
