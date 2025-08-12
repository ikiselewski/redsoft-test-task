oapi-generate: 
	go install oapi-codegen
	oapi-codegen -generate types,gin -package srv -o internal/srv/generated.go docs/openapi.yaml
run:
	docker-compose up --build