package commands

import (
	"lenavire/internal/ledger/application/ports"
	"lenavire/internal/ledger/domain/valuesobjects"
	"lenavire/internal/ledger/infrastructure/adapters"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Receive Payment Command Suite")
}

var _ = Describe("ReceivePaymentCommandHandler", func() {
	var command ReceivePaymentCommand
	var paymentRepository *adapters.InMemoryPaymentRepository
	var expenseRepository *adapters.InMemoryExpenseRepository
	var idProvider ports.IdProvider
	var dateProvider ports.DateProvider
	var handler *ReceivedPaymentCommandHandler
	var ledgerActivityChannel *adapters.FakeLedgerActivityChannel

	BeforeEach(func() {
		paymentRepository = adapters.NewInMemoryPaymentRepository()
		expenseRepository = adapters.NewInMemoryExpenseRepository()
		idProvider = adapters.NewStubIdProvider("xxxxx")
		dateProvider = adapters.NewStubDateProvider("2021-01-04")
		command = NewReceivePaymentCommand(100, "John Doe", "john.doe@example.com", valuesobjects.OneTime)
		ledgerActivityChannel = adapters.NewFakeLedgerActivityChannel()
		handler = NewReceivedPaymentCommandHandler(paymentRepository, idProvider, dateProvider, expenseRepository, ledgerActivityChannel)
	})

	Context("when the command is executed", func() {
		It("should create a new payment", func() {
			handler.Execute(command)
			Expect(paymentRepository.Payments).To(HaveLen(1))
			Expect(paymentRepository.Payments[0].Amount.Value).To(Equal(100))
			Expect(paymentRepository.Payments[0].Name).To(Equal("John Doe"))
			Expect(paymentRepository.Payments[0].Email).To(Equal("john.doe@example.com"))
		})

		It("should assign a payment type", func() {
			handler.Execute(command)
			Expect(paymentRepository.Payments[0].PaymentType).To(Equal(valuesobjects.OneTime))

			command2 := NewReceivePaymentCommand(100, "John Doe", "john.doe@example.com", "recurring")
			handler.Execute(command2)
			Expect(paymentRepository.Payments[1].PaymentType).To(Equal(valuesobjects.Recurring))
		})

		It("should generate an id for the payment", func() {
			handler.Execute(command)
			Expect(paymentRepository.Payments[0].Id).To(Equal("xxxxx"))
		})

		It("should mark the payment date as now", func() {
			handler.Execute(command)
			Expect(paymentRepository.Payments[0].Date).To(Equal("2021-01-04"))
		})

		It("should generate an expense of 20% for the TVA taxes", func() {
			handler.Execute(command)
			Expect(expenseRepository.Expenses).To(HaveLen(1))
			Expect(expenseRepository.Expenses[0].Amount.Value).To(Equal(20))
			Expect(expenseRepository.Expenses[0].Date).To(Equal("2021-01-04"))
			Expect(expenseRepository.Expenses[0].Id).To(Equal("xxxxx"))
		})

		It("should return an error if the amount is not a valid integer", func() {
			command := NewReceivePaymentCommand(-100, "John Doe", "john.doe@example.com", valuesobjects.OneTime)
			err := handler.Execute(command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("amount must be positive: -100"))
		})

		It("should send a notification to the ledger activity channel", func() {
			handler.Execute(command)
			Expect(ledgerActivityChannel.Messages).To(HaveLen(1))
			Expect(ledgerActivityChannel.Messages[0]).To(Equal("payment-received"))
		})
	})
})
