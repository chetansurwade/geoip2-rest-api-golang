# geoip2 rest api golang
GeoIP2 rest api in golang (Gin) just fork and docker-compose up and star it if you like the work.


## To get started :

1. clone the repository

2. prerequisites - docker and docker-compose or go 1.10 or later installed

3. easy way to get up and running is :

`docker-compose up -d --build`

or

`go build main.go`

and run the executable


## Usage:

->  http://localhost:8080/geoip?ip=113.193.190.207

->  http://localhost:8080/geoip?ip=113.193.190.207&ip=2607:f238:2::5
