FROM golang:1.18-alpine3.15 AS go-builder

SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git


WORKDIR /code
COPY . /code/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 728993b91b35037ae8d9933c3a9ee018e49a7926571ce4109f55d9954efcbe9a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep d06607db7bda6d3981f0717133584dd5480a6bca7b1e208b4526e68f3ccf3b31

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
# then log output of file /code/bin/eved
# then ensure static linking
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build \
  && file /code/bin/eved \
  && echo "Ensuring binary is statically linked ..." \
  && (file /code/bin/eved | grep "statically linked")

# --------------------------------------------------------
FROM golang:1.18-alpine3.16

COPY --from=go-builder /code/bin/eved /usr/bin/eved
ENV HOME /eve
WORKDIR $HOME

# rest server
EXPOSE 1317
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# gRPC address
EXPOSE 9090

ENTRYPOINT ["eved"]