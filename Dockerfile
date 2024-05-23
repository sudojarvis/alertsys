
FROM ubuntu:latest AS build


RUN apt-get update && apt-get install -y \
    golang-go \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app


COPY . .


RUN go mod tidy


RUN go build -o my-web-server main.go


FROM ubuntu:latest


WORKDIR /root/


COPY --from=build /app/my-web-server .


EXPOSE 8080


CMD ["./my-web-server"]
