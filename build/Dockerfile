FROM --platform=$BUILDPLATFORM golang:1.22-alpine as builder

ARG VERSION

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -ldflags="-X github.com/qase-tms/qasectl/internal.Version=$VERSION" -o /qli ./main.go

FROM --platform=${TARGETPLATFORM} alpine:3.20 as final

COPY --from=builder /qli /usr/local/bin/qli

ENTRYPOINT ["qli"]
