FROM golang:1.18-alpine3.16 AS go-builder

RUN set -eux

RUN apk add --no-cache ca-certificates git build-base linux-headers

WORKDIR /code
COPY . /code/

# Install babyd binary
RUN echo "Installing eved binary"
RUN make build

#-------------------------------------------
FROM golang:1.18-alpine3.16

RUN apk add --no-cache bash py3-pip jq curl
RUN pip install toml-cli

WORKDIR /

COPY --from=go-builder /code/bin/eved /usr/bin/eved
COPY --from=go-builder /code/bin/eved /


# rest server
EXPOSE 1317
# tendermint rpc
EXPOSE 26657
# p2p address
EXPOSE 26656
# gRPC address
EXPOSE 9090

# wrong ENTRYPOINT can lead to executable not running
ENTRYPOINT ["/bin/bash", "-c"]