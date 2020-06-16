package models_test

import (
	"context"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/models"
	"github.com/jibitters/kiosk/test"
	"github.com/jibitters/kiosk/test/containers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

var _ = Describe("Ticket", func() {
	var pg testcontainers.Container
	var db *pgxpool.Pool
	var repository *models.TicketRepository

	BeforeEach(func() {
		container, port, e := containers.RunPostgres()
		if e != nil {
			Fail(e.Error())
		} else {
			pg = container
		}

		if pool, e := test.ConnectToDatabase(pgHost, port); e != nil {
			Fail(e.Error())
		} else {
			db = pool
			repository = models.NewTicketRepository(zap.S(), db)
		}
	})

	AfterEach(func() {
		db.Close()
		_ = containers.Stop(pg)
	})

	Describe("TicketRepository", func() {
		Context("When Insert called", func() {
			It("Should insert a ticket record in tickets table successfully", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())
			})
		})

		Context("When LoadByID called", func() {
			It("Should load a ticket record from tickets table successfully", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())

				t, e := repository.LoadByID(context.Background(), 1)
				Ω(e).Should(BeNil())
				Ω(t.Issuer).Should(Equal(ticket.Issuer))
				Ω(t.Owner).Should(Equal(ticket.Owner))
				Ω(t.Subject).Should(Equal(ticket.Subject))
				Ω(t.Content).Should(Equal(ticket.Content))
				Ω(t.Metadata).Should(Equal(ticket.Metadata))
				Ω(t.ImportanceLevel).Should(Equal(ticket.ImportanceLevel))
				Ω(t.Status).Should(Equal(models.TicketStatusNew))
				Ω(t.CreatedAt).ShouldNot(BeNil())
				Ω(t.ModifiedAt).ShouldNot(BeNil())
			})

			It("Should return error when provided id does not exists", func() {
				t, e := repository.LoadByID(context.Background(), 1)
				Ω(t).Should(BeNil())
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeNil())
				Ω(e.Errors[0].Code).Should(Equal("ticket.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusNotFound))
			})
		})

		Context("When Update called", func() {
			It("Should update a ticket successfully", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())

				t, e := repository.LoadByID(context.Background(), 1)
				Ω(e).Should(BeNil())

				t.Subject = "Technical Documentation Problem"
				t.Metadata = `{"ip":"192.168.1.10"}`
				t.ImportanceLevel = models.TicketImportanceLevelHigh
				t.Status = models.TicketStatusClosed

				e = repository.Update(context.Background(), t)
				Ω(e).Should(BeNil())

				t, e = repository.LoadByID(context.Background(), 1)
				Ω(e).Should(BeNil())
				Ω(t.Subject).Should(Equal("Technical Documentation Problem"))
				Ω(t.Metadata).Should(Equal(`{"ip":"192.168.1.10"}`))
				Ω(t.ImportanceLevel).Should(Equal(models.TicketImportanceLevelHigh))
				Ω(t.Status).Should(Equal(models.TicketStatusClosed))
			})

			It("Should return error when provided id does not exists", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())

				t, e := repository.LoadByID(context.Background(), 1)
				Ω(e).Should(BeNil())
				t.ID = 100

				e = repository.Update(context.Background(), t)
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeNil())
				Ω(e.Errors[0].Code).Should(Equal("ticket.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusPreconditionFailed))
			})
		})
	})
})
