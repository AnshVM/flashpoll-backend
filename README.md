# Backend for Flashpoll
Check it out [here](http://ec2-13-49-102-133.eu-north-1.compute.amazonaws.com/) 

Check out the [frontend repo](https://github.com/AnshVM/flashpoll)

## Setup on Local

```
$ git clone https://github.com/AnshVM/flashpoll-backend.git
$ cd flashpoll-backend
```
#### Set these variables accordingly in your .env file
```
BCRYPT_COST
DB_PASSWORD
ACCESS_TOKENS_SECRET_KEY
REFRESH_TOKENS_SECRET_KEY
DB_HOST
DB_USER
DB_NAME
DB_PORT
```

```
$ go mod tidy
$ go run main.go
```
