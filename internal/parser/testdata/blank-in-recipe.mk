build: setup
	echo step1

	docker run -v $(DIR):$(DIR):Z \
		/bin/bash -c "go build -a -installsuffix cgo"

setup:
	echo setup
