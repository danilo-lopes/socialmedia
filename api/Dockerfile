FROM golang:1.18.1-alpine3.15 as base

COPY . /go/src/api

WORKDIR /go/src/api

ENV GROUP_ID 10101
ENV USER_ID 10100
ENV OS_USER_NAME sm
ENV OS_GROUP_NAME sm

RUN set -x \
    && addgroup -g ${GROUP_ID} -S ${OS_GROUP_NAME} \
    && adduser -S -D -H -u ${USER_ID} -h /go/src/api -s /sbin/nologin -G ${OS_GROUP_NAME} -g ${OS_GROUP_NAME} ${OS_GROUP_NAME} \
    && chown -R ${USER_ID}.${GROUP_ID} /go/src/api \
    && go mod download all \
    && go mod verify \
    && CGO_ENABLED=0 go build -o /bin/sm

FROM scratch

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=base /bin/sm /bin/sm

EXPOSE 8080

STOPSIGNAL SIGQUIT

USER ${USER_ID}

ENTRYPOINT [ "/bin/sm"]
