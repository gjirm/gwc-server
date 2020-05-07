############################
# STEP 1 build executable binary
############################
# golang debian buster 1.14 linux/amd64
# https://github.com/docker-library/golang/blob/master/1.14/buster/Dockerfile
#FROM golang@sha256:eee8c0a92bc950ecb20d2dffe46546da12147e3146f1b4ed55072c10cacf4f4c as builder
FROM golang@sha256:09b04534495af5148e4cc67c8ac55408307c2d7b9e6ce70f6e05f7f02e427f68 as builder

# Ensure ca-certficates are up to date
RUN update-ca-certificates

WORKDIR $GOPATH/src/gwc/

# use modules
COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
RUN mkdir /gwc
RUN mkdir /gwc/log
RUN mkdir /gwc/db
#RUN mkdir /gwc-server/db
RUN go mod download
RUN go mod verify

COPY app/ app/
COPY internal/ internal/
#COPY vendor/ .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /gwc/gwc-server app/main.go

############################
# STEP 2 build a small image
############################
# using base nonroot image
# user:group is nobody:nobody, uid:gid = 65534:65534
#FROM gcr.io/distroless/base@sha256:2b0a8e9a13dcc168b126778d9e947a7081b4d2ee1ee122830d835f176d0e2a70
FROM gcr.io/distroless/base@sha256:54c459100e9d420e023b0aecc43f7010d2731b6163dd8e060906e2dec4c59890

# Copy our static executable
COPY --from=builder --chown=nonroot:nonroot /gwc/gwc-server /gwc/gwc-server
COPY --from=builder --chown=nonroot:nonroot /gwc/log /gwc/log
COPY --from=builder --chown=nonroot:nonroot /gwc/db /gwc/db

# VOLUME /gwc/log
# VOLUME /gwc/db

WORKDIR /gwc/

EXPOSE 8080

# Run the hello binary.
ENTRYPOINT ["/gwc/gwc-server"]
