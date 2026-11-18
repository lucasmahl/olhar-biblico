# make buildBackend: atualiza os pacotes e builda a api
buildBackend:
	cd backend/ &&\
	cp .env.local .env &&\
	go mod tidy &&\
	go build

buildFrontend:
	cd frontend/ &&\
	cp .env.example .env &&\
	npm i --force &&\
	npm run build
