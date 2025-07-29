package translation

//go:generate mockgen -source=service.go -destination=service_mock_test.go -package=translation -mock_names=Service=MockService Service
