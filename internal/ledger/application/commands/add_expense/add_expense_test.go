package commands

import (
	"lenavire/internal/ledger/application/ports"
	"lenavire/internal/ledger/infrastructure/adapters"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Add Expense Command Suite")
}

var _ = Describe("AddExpenseCommandHandler", func() {
	var expenseRepository adapters.InMemoryExpenseRepository
	var idProvider ports.IdProvider
	var dateProvider ports.DateProvider
	var ledgerActivityChannel adapters.FakeLedgerActivityChannel
	var command AddExpenseCommand
	var handler *AddExpenseCommandHandler

	BeforeEach(func() {
		expenseRepository = *adapters.NewInMemoryExpenseRepository()
		idProvider = adapters.NewStubIdProvider("xxxx")
		dateProvider = adapters.NewStubDateProvider("2006-01-05")
		ledgerActivityChannel = *adapters.NewFakeLedgerActivityChannel()

		command = NewAddExpenseCommand(100)
		handler = NewAddExpenseCommandHandler(&expenseRepository, idProvider, dateProvider, &ledgerActivityChannel)
	})

	Context("when the command is executed", func() {
		It("should create a new expense", func() {
			handler.Execute(command)

			Expect(expenseRepository.Expenses).To(HaveLen(1))
			Expect(expenseRepository.Expenses[0].Amount.Value).To(Equal(100))
		})

		It("should generate an id for the expense", func() {
			handler.Execute(command)

			Expect(expenseRepository.Expenses[0].Id).To(Equal("xxxx"))
		})

		It("should mark the date as now", func() {
			handler.Execute(command)

			Expect(expenseRepository.Expenses[0].Date).To(Equal("2006-01-05"))
		})

		It("should send a notification to the ledger activity channel", func() {
			handler.Execute(command)

			Expect(ledgerActivityChannel.Messages).To(HaveLen(1))
			Expect(ledgerActivityChannel.Messages[0]).To(Equal("expense-added"))
		})
	})

	Context("when the command fails", func() {
		It("should failed for a negative amount", func() {
			command = NewAddExpenseCommand(-100)
			err := handler.Execute(command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("amount must be positive: -100"))
		})

	})
})
