package email

//go:generate mockgen -source=service.go -destination=service_mock.go -package=email -mock_names=Service=MockService Service
