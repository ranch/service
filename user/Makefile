PROJECT=$(shell gcloud config list project | awk 'FNR ==2 { print $$3 }')
NAME=user-svc
TAG=gcr.io/$(PROJECT)/$(NAME)

all: clean build run
	
build:
	docker build -t $(TAG) .

run:
	docker run -d -p 9111:9111 --name=$(NAME) $(TAG)
	
clean:
	-docker stop $(NAME)
	-docker rm $(NAME)
	
push:
	gcloud docker -- push $(TAG)
