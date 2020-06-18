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

var _ = Describe("Comment", func() {
	var pg testcontainers.Container
	var db *pgxpool.Pool
	var ticketRepository *models.TicketRepository
	var repository *models.CommentRepository

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
			ticketRepository = models.NewTicketRepository(zap.S(), db)
			repository = models.NewCommentRepository(zap.S(), db)
		}
	})

	AfterEach(func() {
		db.Close()
		_ = containers.Stop(pg)
	})

	Describe("CommentRepository", func() {
		Context("When Insert called", func() {
			It("Should insert a comment record in comments table successfully", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := ticketRepository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())

				comment := models.Comment{
					TicketID: 1,
					Owner:    "user@example.com",
					Content:  "Hello, we are working on these.",
					Metadata: `{"ip":"192.168.1.1"}`,
				}

				e = repository.Insert(context.Background(), comment)
				Ω(e).Should(BeNil())
			})

			It("Should return error when ticket does not exists", func() {
				comment := models.Comment{
					TicketID: 1,
					Owner:    "user@example.com",
					Content:  "Hello, we are working on these.",
					Metadata: `{"ip":"192.168.1.1"}`,
				}

				e := repository.Insert(context.Background(), comment)
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeEmpty())
				Ω(e.Errors[0].Code).Should(Equal("ticket.not_exists"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusPreconditionFailed))
			})
		})

		Context("When LoadByID called", func() {
			It("Should load a comment record from comments table successfully", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := ticketRepository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())

				comment := models.Comment{
					TicketID: 1,
					Owner:    "admin@example.com",
					Content:  "Hello, we are working on these.!",
					Metadata: `{"ip":"192.168.1.11"}`,
				}

				e = repository.Insert(context.Background(), comment)
				Ω(e).Should(BeNil())

				t, e := repository.LoadByID(context.Background(), 1)
				Ω(e).Should(BeNil())
				Ω(t.TicketID).Should(Equal(int64(1)))
				Ω(t.Owner).Should(Equal(comment.Owner))
				Ω(t.Content).Should(Equal(comment.Content))
				Ω(t.Metadata).Should(Equal(comment.Metadata))
				Ω(t.CreatedAt).ShouldNot(BeNil())
				Ω(t.ModifiedAt).ShouldNot(BeNil())
			})

			It("Should return error when provided id does not exists", func() {
				t, e := repository.LoadByID(context.Background(), 1)
				Ω(t).Should(BeNil())
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeNil())
				Ω(e.Errors[0].Code).Should(Equal("comment.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusNotFound))
			})
		})

		Context("When Update called", func() {
			It("Should update a comment record in comments table successfully", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := ticketRepository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())

				comment := models.Comment{
					TicketID: 1,
					Owner:    "user@example.com",
					Content:  "Hello, we are working on these.",
					Metadata: `{"ip":"192.168.1.1"}`,
				}

				e = repository.Insert(context.Background(), comment)
				Ω(e).Should(BeNil())

				c, e := repository.LoadByID(context.Background(), 1)
				Ω(e).Should(BeNil())

				c.Metadata = `{"ip":"192.168.1.10"}`

				e = repository.Update(context.Background(), c)
				Ω(e).Should(BeNil())
				Ω(c.Metadata).Should(Equal(`{"ip":"192.168.1.10"}`))
			})

			It("Should return error when comment does not exists", func() {
				comment := models.Comment{
					TicketID: 1,
					Metadata: `{"ip":"192.168.1.1"}`,
				}

				e := repository.Update(context.Background(), &comment)
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeEmpty())
				Ω(e.Errors[0].Code).Should(Equal("comment.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusNotFound))
			})
		})

		Context("When DeleteByID called", func() {
			It("Should delete a comment record from comments table successfully", func() {
				ticket := models.Ticket{
					Issuer:          "Microservice-A",
					Owner:           "user@example.com",
					Subject:         "Technical Problem",
					Content:         "Hello, i have some issues with REST API Docs!",
					Metadata:        `{"ip":"192.168.1.1"}`,
					ImportanceLevel: models.TicketImportanceLevelMedium,
				}

				e := ticketRepository.Insert(context.Background(), ticket)
				Ω(e).Should(BeNil())

				comment := models.Comment{
					TicketID: 1,
					Owner:    "admin@example.com",
					Content:  "Hello, we are working on these.!",
					Metadata: `{"ip":"192.168.1.11"}`,
				}

				e = repository.Insert(context.Background(), comment)
				Ω(e).Should(BeNil())

				e = repository.DeleteByID(context.Background(), 1)
				Ω(e).Should(BeNil())

				t, e := repository.LoadByID(context.Background(), 1)
				Ω(t).Should(BeNil())
				Ω(e).ShouldNot(BeNil())
				Ω(e.FingerPrint).ShouldNot(BeNil())
				Ω(e.Errors[0].Code).Should(Equal("comment.not_found"))
				Ω(e.Errors[0].Message).Should(BeEmpty())
				Ω(e.HTTPStatusCode).Should(Equal(http.StatusNotFound))
			})
		})
	})
})
