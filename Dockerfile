FROM golang:1.18-alpine3.15

WORKDIR /work

RUN apk update && apk add git

COPY go.* ./
RUN go mod download

COPY . ./

RUN VERSION=$(cat ./VERSION) && \
    go build \
    -ldflags "-X 'github.com/findy-network/findy-agent-vault/utils.Version=$VERSION'"\
    -o /go/bin/findy-agent-vault

FROM alpine:3.15

LABEL org.opencontainers.image.source https://github.com/findy-network/findy-agent-vault

EXPOSE 8085

# override when running
ENV FAV_JWT_KEY "mySuperSecretKeyLol"
ENV FAV_DB_HOST "vault-db"
ENV FAV_DB_PASSWORD "my-secret-password"
ENV FAV_AGENCY_HOST "localhost"
ENV FAV_AGENCY_PORT "50051"
ENV FAV_AGENCY_CERT_PATH "/grpc"
ENV FAV_AGENCY_ADMIN_ID "findy-root"
ENV FAV_AGENCY_INSECURE "false"

COPY --from=0 /work/db/migrations /db/migrations
COPY --from=0 /go/bin/findy-agent-vault /findy-agent-vault

RUN echo '#!/bin/sh' > /start.sh && \
    echo '/findy-agent-vault' >> /start.sh && chmod a+x /start.sh

ENTRYPOINT ["/bin/sh", "-c", "/start.sh"]
