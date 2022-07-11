FROM golang:1.18.1-alpine3.15

COPY . /go/src/api

WORKDIR /go/src/api

ENV TINI_VERSION v0.19.0
ENV GROUP_ID 10101
ENV USER_ID 10100
ENV OS_USER_NAME sm
ENV OS_GROUP_NAME sm

RUN set -x \
    && addgroup -g ${GROUP_ID} -S ${OS_GROUP_NAME} \
    && adduser -S -D -H -u ${USER_ID} -h /go/src/api -s /sbin/nologin -G ${OS_GROUP_NAME} -g ${OS_GROUP_NAME} ${OS_GROUP_NAME} \
    && go mod download all \
    && go build -o /bin/sm \
    && chown -R ${USER_ID}.${GROUP_ID} /go/src/api \
    && apk add --no-cache tini

EXPOSE 8080

STOPSIGNAL SIGQUIT

USER ${USER_ID}

ENTRYPOINT [ "/sbin/tini", "-s", "--"]

CMD ["sm"]