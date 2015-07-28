# go-dockerevents

Utility for monitoring docker events

## Project Goal

This is a simple project to enable user to quickly plug and output destination
for docker events.

Control behavior by setting the following environment variables

- DOCKER\_HOST and DOCKER\_CERT\_PATH to connect to docker daemon
- DOCKER\_EVENT\_FILTER as a yaml encoded string to filter docker events

DOCKER\_EVENT\_FILTER should be structued as followed

```yaml
image:
    ["ubuntu:latest", "quay.io/coreos/etcd:latest"]
status:
    [create, die, destroy]
container:
    [a925b0d4690c]
```

Fields are options.  The filter follows the following order:

- If the *image* field is not empty, is the event from this image in the *image* field?
- If the *status* field is not empty, is the event type included in the *status* field?
- If the *contianer* field is not empty, Is the contianer ID of this event in the *contianer* field?

## docker2lambda

docker2lambda provides a packaged solution for shipping docker events to **AWS
Lambda** for additional processing.

Examples such as

- Shipping the aggregated events to append to a Google Spreadsheet
- Filter unwanted events by other criteria
- Ship events to database

## docker2kineiss

docker2kineiss behaves much like docker2lambda except events are shipped to
**AWS Kinesis**.

Kinesis is a paid service (no free tier as of current release) and well suited
for mass data ingestion from multiple event source.

If you have a cluster of docker containers, this might suit your need.  For
instance, setup all cluster node to report container events through Kinesis,
then process in batches using **AWS Lambda** setting Kinesis stream as trigger.
