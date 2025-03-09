build-local-image:
	KO_DEFAULTBASEIMAGE=alpine:3 KO_DOCKER_REPO=transaction-mapper ko build . --bare --push=false --local --platform=linux/arm64
