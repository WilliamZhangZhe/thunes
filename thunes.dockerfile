FROM golang:1.15 as golang
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
COPY ./ /data/
WORKDIR /data
RUN go build -o thunes -mod=vendor main.go

FROM ubuntu
WORKDIR /data
EXPOSE 8099
COPY --from=golang /data/thunes /data
COPY ./thunes.toml /data

CMD [ "./thunes", "-c", "./thunes.toml" ]