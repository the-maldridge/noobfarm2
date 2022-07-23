FROM golang:alpine as build

RUN mkdir -p /go/noobfarm2
COPY ./ /go/noobfarm2
RUN cd /go/noobfarm2 && \
        go mod vendor && \
        CGO_ENABLED=0 go build -o /noobfarm2 ./cmd/noobfarm/noobfarm.go

FROM scratch
ARG theme=sample
COPY --from=build /noobfarm2 /
COPY --from=build /go/noobfarm2/themes/$theme/ /theme
CMD ["/noobfarm2"]
VOLUME /data
EXPOSE 8080/tcp
ENV NF_BIND=:8080
ENV NF_QDB=json
ENV NF_JSONROOT=/data
ENV NF_AUTH=file
ENV NF_USER_FILE=/data/accounts.txt
