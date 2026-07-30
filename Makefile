# Build the Docker image
build:
	docker build --platform linux/arm64 -t smol-container:latest .

# Run the container
run:
	docker run -it --rm --privileged \
	    --platform linux/arm64 \
		--name smol-container \
		-p 8888:8888 \
		smol-container:latest $(ARGS)

build-and-run:build run