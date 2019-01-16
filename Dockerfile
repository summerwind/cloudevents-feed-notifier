FROM ubuntu:18.04 AS base

RUN apt update \
  && apt install -y ca-certificates

#############################

FROM scratch

COPY ./cloudevents-feed-notifier /bin/cloudevents-feed-notifier
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/bin/cloudevents-feed-notifier"]
