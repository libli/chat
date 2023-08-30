.DEFAULT_GOAL := build
.PHONY: build clean push
build:
	docker build --platform linux/amd64 -t libli/chat:1.5 -t libli/chat:latest .
push:
	@echo "Pushing to docker hub"
	docker login -u libli -p $(DOCKER_PASSWORD)
	docker push libli/chat -a
clean:
	docker rmi libli/chat:latest