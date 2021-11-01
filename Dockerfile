# build
FROM golang:1.17.2-alpine AS build
COPY . /src
WORKDIR /src
ENV "GOPROXY" "https://goproxy.io,direct"
RUN go build -o /build/app

# iamge
FROM alpine:latest
COPY --from=build /build/app /bin/app
WORKDIR /
ENTRYPOINT [ "/bin/app" ]