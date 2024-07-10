# Define image name
IMAGE_NAME = cornelius

# Define default build arguments (no CMD)
DEFAULT_CMD ?= ""

.PHONY: all build run clean

all: build

build:
	docker build -t $(IMAGE_NAME):latest .

run:
	# Allow user to override CMD via $(CMD) argument
	CMD ?= $(DEFAULT_CMD)
	docker run -it --rm $(IMAGE_NAME):latest $(CMD)

clean:
	docker rmi $(IMAGE_NAME):latest