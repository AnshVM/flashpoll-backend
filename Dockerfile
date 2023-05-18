FROM golang:alpine
ENV SERVER_DOMAIN=localhost
ENV PORT=:8080
WORKDIR /app
COPY . .
RUN go build -o flashpoll_build
EXPOSE 80
CMD ["./flashpoll_build"]