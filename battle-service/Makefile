format:
	buf format -w proto

lint:
	buf lint proto --config buf.yaml

generate:
	buf generate proto

build_dev:
	docker build --build-arg ENV_FILE=dev -t battle-service .