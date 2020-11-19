############################
# STEP 1 build executable binary
############################
FROM golang:1.14-alpine as builder

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates pkgconfig

WORKDIR /app
COPY svc  /go/bin/svc

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /go/bin/svc /svc

# Port on which the service will be exposed.
EXPOSE 3000

# Run the svc binary.
CMD ["./svc"]