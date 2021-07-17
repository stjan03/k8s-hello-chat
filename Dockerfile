# Demo app
# By Stefan
# 20210715

# docker build -f Dockerfile -t k8s-hello-app .


# --- compile the binaries ---------------
FROM golang:latest as builder

# Copy stuff over
WORKDIR /app
COPY k8s-hello-app.go ./

# Compile
ENV GO111MODULE=off
RUN go get -d github.com/prometheus/client_golang/prometheus/promhttp
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-hello-app ${MAIN_PATH}



# --- build the final container ---------------
FROM alpine:latest
LABEL maintainer="Stefan"

# Copy stuff over
WORKDIR /app
COPY --from=builder /app/k8s-hello-app .

# Set default enviroment variables
ENV APP_VERSION="v1.0"
ENV APP_TEXT="Testing..."

# Start the webserver
EXPOSE 8080
CMD ["./k8s-hello-app"]
