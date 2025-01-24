run-postgres:
	docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:17.2

run-mysql:
	docker run -d \
		--name mysql-container \
		-p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=admin \
		-e MYSQL_USER=admin \
		-e MYSQL_PASSWORD=admin \
		-e MYSQL_DATABASE=mydb \
		mysql:9.1

run-redis:
	docker run --name redis -p 6379:6379 -d redis

build-img:
	docker build -t go-server-template .

run-container:
	docker run -d --name go-server-template go-server-template

backup:
	# ./scripts/add_gitkeep.sh
	# golines . -w -m 99
	
	# gofumpt -w .
	
	git add .
	git commit -m "backup"
	git push

gen-rest-doc:
	swag init --output ./docs/openapi  -g main.go internal/controller/*.go