#Go Lambda WebApp Template

This repository is a template that can be used to start building a self contained web-application written in Go and deployed on AWS Lambda.

This template attempts to make use of minimal key eternal components while providing a ready-to-use template to develop an all-in-one Lambda web application.

This template makes use of:
- Mark Bates [Pkger](https://github.com/markbates/pkger) to statically include templates and static web assets.
- [Apex Gateway](https://github.com/apex/gateway) to map between Lambda requests and standard net.http requests.
- [CHI](github.com/go-chi/chi) to route incoming requests.

## Usage

Install [pkger](https://github.com/markbates/pkger)
```go get github.com/markbates/pkger/cmd/pkger
```

Clone the Repo
```git clone https://github.com/ericdaugherty/go-lambda-webapp-template
```

Since you are creating your own project, remove the remote Repo Reference
```git remote rm origin
```

Edit go.mod and change the module name to reflect your new project.

Run Locally
```make run
```
This will start the default webserver on http://localhost:3000

## Deploy To Lambda

The template includes a [serverless](https://serverless.com/) file to support easy deployment to AWS Lambda, including creating an API Gateway HTTP Endpoint.

To deploy, first make sure [serverless](https://serverless.com/) is installed and configured with your AWS credentials.

Then deploy:
```make deploy
```

This should create a new CloudFormation stack on AWS including your all-in-one lambda and an API Gateway HTTP Endpoint. 
