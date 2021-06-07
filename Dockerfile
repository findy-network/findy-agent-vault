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

RUN VERSION=$(cat ./VERSION) && \
    go build \
    -ldflags "-X 'github.com/findy-network/findy-agent-vault/utils.Version=$VERSION'"\
    -o /go/bin/findy-agent-vault

FROM ghcr.io/findy-network/findy-base:alpine-3.13

EXPOSE 8085

# override when running
ENV FAV_JWT_KEY "mySuperSecretKeyLol"
ENV FAV_DB_HOST "vault-db"
ENV FAV_DB_PASSWORD "my-secret-password"
ENV FAV_AGENCY_HOST "localhost"
ENV FAV_AGENCY_PORT "50051"
ENV FAV_AGENCY_CERT_PATH "/grpc"
ENV FAV_AGENCY_ADMIN_ID "findy-root"

COPY --from=0 /work/db/migrations /db/migrations
COPY --from=0 /go/bin/findy-agent-vault /findy-agent-vault

RUN echo '/s3-copy $STARTUP_FILE_STORAGE_S3 grpc /' > /start.sh && \
    echo '/findy-agent-vault' >> /start.sh && chmod a+x /start.sh

ENTRYPOINT ["/bin/sh", "-c", "/start.sh"]
