build:
	go build -o bin/unifi-captive-portal
docker:
	docker build -t unifi-captive-portal .