# cloudevents-feed-notifier

cloudevents-feed-notifier watches the update of XML feeds and notify the new entry of feed as CloudEvents.

## Install

Download the latest binary from the [Releases](https://github.com/summerwind/cloudevents-feed-notifier/releases) page.

Docker image is also available. Running with Docker is as follows.

```
$ docker run -it -v $PWD/config.yml:/config.yml summerwind/cloudevents-feed-notifier:latest
```

## Usage

cloudevents-feed-notifier can be started from the command line as follows.

```
$ cloudevents-feed-notifier -c config.yml
```

To start cloudevents-feed-notifier, specify the configuration file using the `-c` option. The configuration format is in YAML. Please see `example/config.yml` for the full configuration file format.

## Event

Events are sent as CloudEvents in the following headers and body.

### Request headers

| Header | Value |
| --- | --- |
| CE-SpecVersion | 0.2 |
| CE-Time        | (Published time) |
| CE-ID          | (Auto generated UUID) |
| CE-Type        | feed.publish |
| CE-Source      | (Entry URL) |
| Content-Type   | application/json |

### Request body

```
{
  "feed": {
    "title": "XML feed title",
    "url": "https://xml.feed",
  },
  "entry": {
    "title": "Entry title",
    "url": "https://xml.feed/entry",
  }
}
```
