FROM golang:1.12.4-stretch as builder
LABEL maintainer="Darren Clarke <darren@redaranj.com>"
RUN mkdir /build
WORKDIR /build
COPY src .
RUN go build

FROM golang:1.12.4-stretch
LABEL maintainer="Darren Clarke <darren@redaranj.com>"
ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL
ARG VERSION

LABEL org.label-schema.schema-version="1.1"
LABEL org.label-schema.name="redaranj/quepasa"
LABEL org.label-schema.description="Write WhatsApp bots with https"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.vcs-url=$VCS_URL
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.version=$VERSION

RUN mkdir /app
WORKDIR /app
COPY --from=builder /build/quepasa ./quepasa
COPY --from=builder /build/views ./views
COPY --from=builder /build/assets ./assets
COPY --from=builder /build/migrations ./migrations
COPY docker-entrypoint.sh /usr/local/bin/
EXPOSE 3000
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/app/quepasa"]
