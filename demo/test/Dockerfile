FROM findy-agent-cli

FROM optechlab/indy-golang:1.16.0

ARG HTTPS_PREFIX

ENV GOPRIVATE "github.com/findy-network"

RUN git config --global url."https://"${HTTPS_PREFIX}"github.com/".insteadOf "https://github.com/"

COPY --from=0 /findy-agent-cli /findy-agent-cli

WORKDIR /work

COPY go* ./

RUN go mod download

COPY *.go ./

RUN go build -o /jwt-extractor

RUN apt-get update && apt-get install -y curl

COPY test.sh /

COPY .docker/cert /go/src/github.com/findy-network/findy-common-go/cert

ENTRYPOINT ["/test.sh"]
