FROM --platform=$BUILDPLATFORM golang:1.26-alpine as builder

ARG VERSION

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -ldflags="-s -w -X github.com/qase-tms/qasectl/internal.Version=$VERSION" -o /qasectl ./main.go

FROM --platform=${TARGETPLATFORM} alpine:3.20 as final

COPY --from=builder /qasectl /usr/local/bin/qasectl

ENTRYPOINT ["qasectl"]
