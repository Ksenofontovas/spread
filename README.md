######
###  Docker
###### 
sudo docker run --name=pg_spread -e POSTGRES_PASSWORD="qwerty" -p 5432:5432 -d --rm postgres
docker compose --env-file ../.env.dev up
######
###  Migrations
######
migrate create -ext sql -dir ./schema -seq init

migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/symbols?sslmode=disable' up
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/symbols?sslmode=disable' down