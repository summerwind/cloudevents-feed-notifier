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

Event are sent as CloudEvents using HTTP Transport Binding (Binary Content Mode).

### Request headers

| Header | Value |
| --- | --- |
| CE-SpecVersion | 0.2 |
| CE-Time        | 2019-01-16T19:41:51+09:00 (Published time) |
| CE-ID          | cffa4fa7-095b-485c-aaa8-a28e98a5f897 (Auto generated UUID) |
| CE-Type        | feed.publish |
| CE-Source      | https://github.com/kubernetes/kubernetes/releases/tag/v1.13.2 (Entry URL) |
| Content-Type   | application/json |

### Request body

```
{
  "feed":{
    "title": "Release notes from kubernetes",
    "url": "https://github.com/kubernetes/kubernetes/releases.atom"
  },
  "entry":{
    "title": "v1.13.2",
    "url": "https://github.com/kubernetes/kubernetes/releases/tag/v1.13.2",
    "published_at": "2019-01-11T02:18:07Z"
  }
}
```
