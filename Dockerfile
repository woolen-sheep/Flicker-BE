FROM alpine:latest
COPY ./app /bin/app
WORKDIR /
ENTRYPOINT [ "/bin/app" ]