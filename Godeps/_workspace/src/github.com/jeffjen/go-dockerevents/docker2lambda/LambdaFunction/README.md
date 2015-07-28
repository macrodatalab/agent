# How to setup AWS Lambda endpoint for go-dockerevents

AWS Lambda allows backend to be setup without a server.  Uses include:

- Data transformation
- Filter
- Routing to backend storage

Without the need to setup a dedicated server.  Further, it claims (not tested throughly) that it scales up computing node as hit rate increases (for a price of course).

# DockerReport

DockerReport is an sample setup for shipping docker container events to Google Spreadsheet for later analysis.  DockerReport assumes that the incoming AWS Lambda event is a json encoded array of objects, each object reprents a docker event.

Reference the instructions [here](http://docs.aws.amazon.com/lambda/latest/dg/walkthrough-s3-events-adminuser-create-test-function-create-function.html) to create a deployment package.
