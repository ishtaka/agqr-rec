FROM golang:1.15.3-alpine AS build
COPY ./ ./

ENV GO111MODULE=on
ENV GOPATH=''
RUN mkdir -p /build
RUN go build -o=/build/agqrrec cmd/agqrrec/main.go

FROM alpine:latest
RUN apk --update --no-cache add ffmpeg tzdata

COPY --from=build /build/agqrrec agqrrec
RUN chmod u+x agqrrec
COPY ./configs ./configs

ENTRYPOINT ["/agqrrec"]
