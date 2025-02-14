### 빌더 ###
FROM golang:1.23-alpine AS builder
LABEL maintainer="devphilip21 <philip21.dev@gmail.com>"

WORKDIR /build

# 빌드 env 셋팅
ENV GOOS=linux GOARCH=amd64
ENV PATH /go/bin:$PATH

# github private repo 에 go 패키지모듈이 접근하기 위한 설정
RUN apk update && apk upgrade && apk add openssh
RUN apk add --no-cache git

RUN mkdir ~/.ssh
RUN touch ~/.ssh/known_hosts
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
RUN git config --global url."https://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:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# 복사 및 빌드
COPY . .
RUN go mod tidy
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sbec ./application/cmd/main.go

# 몽고 Atlas x509 에러 핸들링을 위한 패키지 설치 (스크래치에 포함시키기 위함)
# @see: https://gist.github.com/michaelboke/564bf96f7331f35f1716b59984befc50
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates

# timezone 적용
RUN apk add tzdata

### 메인 ###
FROM scratch AS runner

ARG PHASE

ENV APP_PHASE prod

COPY --from=builder /build/sbec /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 실행 엔트리
ENTRYPOINT ["/sbec"]
