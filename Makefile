docs:
	rm -rf api
	mkdir api
	swag init --parseDependency --generalInfo internal/api/app.go --output ./api