FRONTEND_BINARY=frontEndApp
BROKER_BINARY=brokerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

up_build_push: up_build push_images

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_frontend: builds the frontend binary as a linux executable
build_frontend:
	@echo "Building frontend binary..."
	cd ../frontend && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONTEND_BINARY} ./cmd/web
	@echo "Done!"


## push_images: pushes all images to the docker hub
push_images:
	@echo "Pushing images to docker hub..."
	docker tag project-frontend-service calvarado2004/portworx-client-frontend:latest
	docker push calvarado2004/portworx-client-frontend

	docker tag project-broker-service calvarado2004/portworx-client-broker:latest
	docker push calvarado2004/portworx-client-broker

	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../frontend && ./${FRONTEND_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONTEND_BINARY}"
	@echo "Stopped front end!"
