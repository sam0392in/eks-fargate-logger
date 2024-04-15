#---------------------------STAGE 1---------------------------------#

FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .
COPY Dockerfile Dockerfile

# Build code
RUN go mod tidy
RUN go build -o app && \
    mv app /home

#---------------------------STAGE 2---------------------------------#

FROM golang:1.22-alpine

WORKDIR /home

COPY --from=builder /home/app ./app
COPY config config

EXPOSE 3000

CMD ["./app"]