FROM golang:1.15-alpine3.12

ARG HTTPS_PREFIX

ENV GOPRIVATE "github.com/findy-network"

RUN apk update && \
    apk add git && \
    git config --global url."https://"${HTTPS_PREFIX}"github.com/".insteadOf "https://github.com/"

WORKDIR /work

COPY go.* .
RUN go mod download

COPY . .

RUN go build -o /go/bin/findy-agent-vault

FROM alpine:3.12

# override when running
ENV FAV_JWT_KEY "mySuperSecretKeyLol"
ENV FAV_DB_HOST "vault-db"
ENV FAV_DB_PASSWORD "my-secret-password"

COPY --from=0 /go/bin/findy-agent-vault /findy-agent-vault

ENTRYPOINT ["/findy-agent-vault"]
