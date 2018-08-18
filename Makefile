build:
	npm install
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/record record/main.go
	cp ./record/serviceAccount.json .

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	npm run deploy
