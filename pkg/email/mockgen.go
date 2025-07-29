package email

//go:generate mockgen -source=service.go -destination=service_mock_test.go -package=email -mock_names=Service=MockService Service
