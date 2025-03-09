build-local-image:
	KO_DOCKER_REPO=transaction-mapper ko build . --bare --push=false --local --platform=linux/arm64
