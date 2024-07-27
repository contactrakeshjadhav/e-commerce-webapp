# Builder
# -----------------------------------------------------------------
FROM golang:1.19.4-alpine3.16 as builder

RUN apk add --update make bash git build-base

ADD ./ /metzeal/ecommerce/product

WORKDIR /metzeal/ecommerce/product
RUN make tools
RUN GOBIN=/metzeal/ecommerce/product/bin time make build

# -----------------------------------------------------------------
# Runner
# -----------------------------------------------------------------
FROM alpine:3.16

RUN apk add --no-cache --update ca-certificates curl coreutils bash

COPY --from=builder /metzeal/ecommerce/product/bin/product /usr/bin/
ADD database/migrations /migrations

EXPOSE 8080

CMD ["product"]
