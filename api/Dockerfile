FROM golang:1.18.1-alpine3.15 as base

COPY . /go/src/api

WORKDIR /go/src/api

RUN set -x \
    && go mod download all \
    && go mod verify \
    && CGO_ENABLED=0 go build -o /bin/sm

FROM alpine:3.15.5

ENV CONTAINER_USER sm
ENV CONTAINER_GROUP sm

ENV CONTAINER_GROUP_ID 1001
ENV CONTAINER_USER_ID 1002

COPY --from=base /bin/sm /bin/sm

RUN set -x \
    && apk add --no-cache tini \
    && addgroup -g ${CONTAINER_GROUP_ID} -S ${CONTAINER_GROUP} \
    && adduser -H -S -D -u ${CONTAINER_USER_ID} -G ${CONTAINER_GROUP} ${CONTAINER_USER} \
    && chown ${CONTAINER_USER_ID}.${CONTAINER_GROUP_ID} /bin/sm && chmod +x /bin/sm

EXPOSE 8080

STOPSIGNAL SIGQUIT

USER ${CONTAINER_USER_ID}

ENTRYPOINT ["/sbin/tini", "-s", "--"]
CMD ["sm"]