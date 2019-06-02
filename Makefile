all:

build:
	go build -v -o ./bin/artists ./cmd/artists/...
	go build -v -o ./bin/artistsctl ./cmd/artistsctl/...

install:
	go install -v ./cmd/artists/...
	go install -v ./cmd/artistsctl/...

run:
	go run ./cmd/artists/...

runctl:
	go run ./cmd/artistsctl/...

tests t:
	go test -v ./internal/...

lint-all l:
	bash ./scripts/golangci-lint.sh
	bash ./scripts/consistent.sh

add-ssh-key:
	openssl aes-256-cbc -K $(encrypted_f6f9818801b5_key) -iv $(encrypted_f6f9818801b5_iv) -in travis_key.enc -out /tmp/travis_key -d
	chmod 600 /tmp/travis_key
	ssh-add /tmp/travis_key

docker-login:
	docker login -u $(REGISTRY_USER) -p $(REGISTRY_PASS)

docker-build:
	docker build -t $(REGISTRY_REPO):$(VERSION) .

docker-push: docker-login
	docker push $(REGISTRY_REPO):$(VERSION)

deploy:
	ssh -o "StrictHostKeyChecking no" $(HOST_USER)@$(HOST) make run-music-artists

deploy-staging:
	ssh -o "StrictHostKeyChecking no" $(STAGING_USER)@$(STAGING_HOST) make run-music-artists
