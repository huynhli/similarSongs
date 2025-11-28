# Instruction go here

# start from this command --> declares go image
FROM golang:1.23-alpine AS builder

# process:
    # make folder + cd
    # copy dependencies (go.mod + go.sum, or package.json etc...), download dependencies
    # copy overything else

# make a /app folder to be the container + cd into
WORKDIR /app

# doing first caches
# copy package.json + lock if have
COPY go.mod go.sum ./
# downloading dependencies
RUN go mod download

# copy over other folder
COPY . .

# compile backend
RUN go build -o main

# final stage: take compiled backend(binary) and run it

FROM alpine:latest
# make + cd into /app
WORKDIR /app
# from builder stage, take main file
COPY --from=builder /app/main .

# Set port
ENV PORT=8080
EXPOSE 8080

# when container starts, docker runs this
CMD ["./main"] 

# docker build -t name_of_build/fiber_backend .
# docker run --env-file .env -p 8080:8080 fiber_backend
# OR
# if yml, run:
# docker-compose up --build