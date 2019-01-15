FROM scratch

COPY ./cloudevents-feed-notifier /bin/cloudevents-feed-notifier

ENTRYPOINT ["/bin/cloudevents-feed-notifier"]
