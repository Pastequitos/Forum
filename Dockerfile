FROM golang:latest
COPY . /src
WORKDIR /src
RUN go build -o Forum
CMD ["./Forum"]