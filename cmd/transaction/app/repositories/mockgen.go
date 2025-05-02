package repositories

//go:generate mockgen -source=repository.go -destination=../../../../test/mocks/transaction_repository.go --package=mocks -mock_names=Repository=TransactionsRepository Repository
