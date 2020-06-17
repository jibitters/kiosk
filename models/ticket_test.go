package models_test

import (
	"context"
	"net/http"
	"time"

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
	var commentRepository *models.CommentRepository

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
			commentRepository = models.NewCommentRepository(zap.S(), db)
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
				Ω(t.Comments).Should(BeEmpty())
			})

			It("Should load a ticket record and all of its comments successfully", func() {
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

				comment := models.Comment{
					TicketID: 1,
					Owner:    "user@example.com",
					Content:  "Hello, we are working on these.",
					Metadata: `{"ip":"192.168.1.11"}`,
				}

				e = commentRepository.Insert(context.Background(), comment)
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
				Ω(len(t.Comments)).Should(Equal(1))
				Ω(t.Comments[0].ID).Should(Equal(int64(1)))
				Ω(t.Comments[0].TicketID).Should(Equal(comment.TicketID))
				Ω(t.Comments[0].Owner).Should(Equal(comment.Owner))
				Ω(t.Comments[0].Content).Should(Equal(comment.Content))
				Ω(t.Comments[0].Metadata).Should(Equal(comment.Metadata))
				Ω(t.Comments[0].CreatedAt).ShouldNot(BeNil())
				Ω(t.Comments[0].ModifiedAt).ShouldNot(BeNil())
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

		Context("When DeleteByID called", func() {
			It("Should delete a ticket record from tickets table successfully", func() {
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

				e = repository.DeleteByID(context.Background(), 1)
				Ω(e).Should(BeNil())

				t, e := repository.LoadByID(context.Background(), 1)
				Ω(t).Should(BeNil())
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeNil())
				Ω(e.Errors[0].Code).Should(Equal("ticket.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusNotFound))
			})

			It("Should delete a ticket record and all of its comments successfully", func() {
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

				comment := models.Comment{
					TicketID: 1,
					Owner:    "user@example.com",
					Content:  "Hello, we are working on these.",
					Metadata: `{"ip":"192.168.1.11"}`,
				}

				e = commentRepository.Insert(context.Background(), comment)
				Ω(e).Should(BeNil())

				e = repository.DeleteByID(context.Background(), 1)
				Ω(e).Should(BeNil())

				t, e := repository.LoadByID(context.Background(), 1)
				Ω(t).Should(BeNil())
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeNil())
				Ω(e.Errors[0].Code).Should(Equal("ticket.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusNotFound))

				c, e := commentRepository.LoadByID(context.Background(), 1)
				Ω(c).Should(BeNil())
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeNil())
				Ω(e.Errors[0].Code).Should(Equal("comment.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusNotFound))
			})
		})

		Context("When Filter called", func() {
			It("Should filter tickets successfully", func() {
				ticket1 := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user1@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket1)
				Ω(e).Should(BeNil())

				comment1 := models.Comment{
					TicketID: 1,
					Owner:    "admin@example.com",
					Content:  "Ok, we will check!",
					Metadata: `{"ip":"192.168.1.11"}`,
				}

				e = commentRepository.Insert(context.Background(), comment1)
				Ω(e).Should(BeNil())

				comment2 := models.Comment{
					TicketID: 1,
					Owner:    "user1@example.com",
					Content:  "Ok",
					Metadata: `{"ip":"192.168.1.1"}`,
				}

				e = commentRepository.Insert(context.Background(), comment2)
				Ω(e).Should(BeNil())

				ticket2 := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user2@example.com",
					Subject:         "UI Problem",
					Content:         "Hello, i have some issues with panel!",
					Metadata:        `{"ip":"192.168.1.2"}`,
					ImportanceLevel: models.TicketImportanceLevelLow,
				}

				e = repository.Insert(context.Background(), ticket2)
				Ω(e).Should(BeNil())

				comment3 := models.Comment{
					TicketID: 2,
					Owner:    "admin@example.com",
					Content:  "Ok, we will check!",
					Metadata: `{"ip":"192.168.1.11"}`,
				}

				e = commentRepository.Insert(context.Background(), comment3)
				Ω(e).Should(BeNil())

				ts, hasNextPage, e := repository.Filter(context.Background(), "", "", "",
					"", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano), time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano),
					1, 10)

				Ω(e).Should(BeNil())
				Ω(len(ts)).Should(Equal(2))
				Ω(hasNextPage).Should(Equal(false))
				Ω(len(ts[0].Comments)).Should(Equal(1))
				Ω(len(ts[1].Comments)).Should(Equal(2))
				Ω(ts[1].Comments[0].ModifiedAt.After(ts[1].Comments[1].ModifiedAt)).Should(BeTrue())
			})

			It("Should filter tickets by issuer", func() {
				ticket1 := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user1@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket1)
				Ω(e).Should(BeNil())

				ticket2 := models.Ticket{
					Issuer:          "Microservice-B",
					Owner:           "user2@example.com",
					Subject:         "UI Problem",
					Content:         "Hello, i have some issues with panel!",
					Metadata:        `{"ip":"192.168.1.2"}`,
					ImportanceLevel: models.TicketImportanceLevelLow,
				}

				e = repository.Insert(context.Background(), ticket2)
				Ω(e).Should(BeNil())

				ts, hasNextPage, e := repository.Filter(context.Background(), "Microservice-A", "", "",
					"", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano), time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano),
					1, 10)

				Ω(e).Should(BeNil())
				Ω(len(ts)).Should(Equal(1))
				Ω(hasNextPage).Should(Equal(false))
			})

			It("Should filter tickets by owner", func() {
				ticket1 := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user1@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket1)
				Ω(e).Should(BeNil())

				ticket2 := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user2@example.com",
					Subject:         "UI Problem",
					Content:         "Hello, i have some issues with panel!",
					Metadata:        `{"ip":"192.168.1.2"}`,
					ImportanceLevel: models.TicketImportanceLevelLow,
				}

				e = repository.Insert(context.Background(), ticket2)
				Ω(e).Should(BeNil())

				ts, hasNextPage, e := repository.Filter(context.Background(), "Microservice-A", "user1@example.com", "",
					"", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano), time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano),
					1, 10)

				Ω(e).Should(BeNil())
				Ω(len(ts)).Should(Equal(1))
				Ω(hasNextPage).Should(Equal(false))
			})

			It("Should return paginated response correctly", func() {
				ticket1 := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user1@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := repository.Insert(context.Background(), ticket1)
				Ω(e).Should(BeNil())

				ticket2 := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user2@example.com",
					Subject:         "UI Problem",
					Content:         "Hello, i have some issues with panel!",
					Metadata:        `{"ip":"192.168.1.2"}`,
					ImportanceLevel: models.TicketImportanceLevelLow,
				}

				e = repository.Insert(context.Background(), ticket2)
				Ω(e).Should(BeNil())

				ts, hasNextPage, e := repository.Filter(context.Background(), "", "", "",
					"", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano), time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano),
					1, 1)

				Ω(e).Should(BeNil())
				Ω(len(ts)).Should(Equal(1))
				Ω(hasNextPage).Should(Equal(true))

				ts, hasNextPage, e = repository.Filter(context.Background(), "", "", "",
					"", time.Now().UTC().Add(-time.Hour).Format(time.RFC3339Nano), time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano),
					2, 1)

				Ω(e).Should(BeNil())
				Ω(len(ts)).Should(Equal(1))
				Ω(hasNextPage).Should(Equal(false))
			})
		})
	})
})
