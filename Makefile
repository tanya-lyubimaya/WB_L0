clean:
	rm -f l0

build: clean
	go build -o l0 main.go

run: build
	./l0

nats:
	sudo docker container stop nats || true && \
    sudo docker container rm nats || true && \
    sudo docker-compose -f docker-compose.yml up --detach && \
    go install github.com/nats-io/natscli/nats@latest && \
    go run github.com/nats-io/natscli/nats@latest context add nats --server 127.0.0.1:4222 --description "nats server"

pub:
	./pub_to_nats.sh