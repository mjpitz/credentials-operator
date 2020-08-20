FROM golang:1.14-alpine3.12 AS builder

ADD bin/ /usr/bin/

RUN GOOS=$(go env GOOS) && \
    GOARCH=$(go env GOARCH) && \
    mv /usr/bin/credentials-operator_${GOOS}_${GOARCH} /usr/bin/credentials-operator

FROM alpine:3.12

RUN apk -U update && apk add ca-certificates

COPY --from=builder /usr/bin/credentials-operator /usr/bin/credentials-operator

USER 13490:13490

ENTRYPOINT [ "credentials-operator" ]
