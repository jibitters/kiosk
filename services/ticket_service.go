package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
	"github.com/jibitters/kiosk/web/data"
	nc "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// TicketService is a service implementation of ticket related functionalities.
type TicketService struct {
	logger           *zap.SugaredLogger
	natsClient       *nc.Conn
	ticketRepository models.TicketRepository
	stop             chan struct{}
}

// NewTicketService returns a newly created and ready to use TicketService.
func NewTicketService(logger *zap.SugaredLogger, natsClient *nc.Conn) *TicketService {
	return &TicketService{
		logger:     logger,
		natsClient: natsClient,
		stop:       make(chan struct{}),
	}
}

// Start starts the subscription so ready to be notified.
func (s *TicketService) Start() error {
	createTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.create",
		"kiosk.tickets.create_group", s.createTicket)
	if e != nil {
		return e
	}

	loadTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.load",
		"kiosk.tickets.load_group", s.loadTicket)
	if e != nil {
		return e
	}

	updateTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.update",
		"kiosk.tickets.update_group", s.updateTicket)
	if e != nil {
		return e
	}

	deleteTicketSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.delete",
		"kiosk.tickets.delete_group", s.deleteTicket)
	if e != nil {
		return e
	}

	filterTicketsSubscription, e := s.natsClient.QueueSubscribe("kiosk.tickets.filter",
		"kiosk.tickets.filter_group", s.filterTickets)
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

func (s *TicketService) createTicket(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createTicketRequest := &data.CreateTicketRequest{}
	e := json.Unmarshal(msg.Data, createTicketRequest)
	if e != nil {
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

func (s *TicketService) loadTicket(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := &data.ID{}
	e := json.Unmarshal(msg.Data, id)
	if e != nil {
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

func (s *TicketService) updateTicket(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateTicketRequest := &data.UpdateTicketRequest{}
	e := json.Unmarshal(msg.Data, updateTicketRequest)
	if e != nil {
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

func (s *TicketService) deleteTicket(msg *nc.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := &data.ID{}
	e := json.Unmarshal(msg.Data, id)
	if e != nil {
		s.reply(msg, errors.InvalidRequestBody())
		return
	}

	if e := s.ticketRepository.DeleteByID(ctx, id.ID); e != nil {
		s.reply(msg, e)
		return
	}

	s.replyNoContent(msg)
}

func (s *TicketService) filterTickets(msg *nc.Msg) {
	s.reply(msg, errors.NotImplemented())
}

func (s *TicketService) reply(msg *nc.Msg, t interface{}) {
	reply, _ := json.Marshal(t)
	_ = msg.Respond(reply)
}

func (s *TicketService) replyNoContent(msg *nc.Msg) {
	_ = msg.Respond([]byte(""))
}

// Stop stops the component and it subscription.
func (s *TicketService) Stop() {
	s.stop <- struct{}{}
}
