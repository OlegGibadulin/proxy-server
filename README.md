# Proxy server

Proxy server for handling HTTP and HTTPS requests.


## Run

* Init database

	`psql -U postgres -d proxy -a -f scripts/init.sql`

* Generate keys

	`./scripts/gen_ca.sh`

* Run server

	`go run cmd/app/main.go`

 

Proxy server [http://localhost:8080](http://localhost:8080)

Web server   [http://localhost:8000](http://localhost:8000)
