FROM golang:1.16-alpine3.13

ARG HTTPS_PREFIX

ENV GOPRIVATE "github.com/findy-network"

RUN apk update && \
    apk add git && \
    git config --global url."https://"${HTTPS_PREFIX}"github.com/".insteadOf "https://github.com/"

WORKDIR /work

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o /go/bin/findy-agent-vault

FROM alpine:3.13

# override when running
ENV FAV_JWT_KEY "mySuperSecretKeyLol"
ENV FAV_DB_HOST "vault-db"
ENV FAV_DB_PASSWORD "my-secret-password"

COPY --from=0 /work/db/migrations /db/migrations
COPY --from=0 /go/bin/findy-agent-vault /findy-agent-vault

RUN echo 'sleep 20 && /findy-agent-vault' > /start.sh && chmod a+x /start.sh

ENTRYPOINT ["/bin/sh", "-c", "/start.sh"]
