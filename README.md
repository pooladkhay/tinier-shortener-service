
# Tinier Shortener Service
Tinier is a URL Shortener REST API written in go.

## Installation
Open ```docker-compose.yml``` and set env variables for cassandra and redis addresses properly.

```
$ git clone https://github.com/pooladkhay/tinier-shortener-service.git
$ cd tinier-shortener-service
$ docker-compose up --build
```
## Design

![Ttinier](tinier-design.png?raw=true)