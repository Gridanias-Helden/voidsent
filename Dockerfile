FROM golang:1.20-alpine AS builder

COPY ./ /app

RUN cd /app && \
    CGO_ENABLED=0 go build -ldflags="-s -w" .

# ----------------------------------------

FROM scratch

WORKDIR /app

COPY --from=builder /app /
COPY ./static /static

CMD ["/voidsent"]

EXPOSE 80
