# Go Lambda WebApp Template
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

This repository is a template that can be used to start building a self contained web-application written in Go and deployed on AWS Lambda.

This template makes use of minimal eternal components while providing a ready-to-use template to develop an all-in-one Lambda web application.

This template makes use of:
- [Pkger](https://github.com/markbates/pkger) to statically include templates and static web assets.
- [Apex Gateway](https://github.com/apex/gateway) to map between Lambda requests and standard net.http requests.
- [CHI](github.com/go-chi/chi) to route incoming requests.

## Usage

**Install [pkger](https://github.com/markbates/pkger)**

    go get github.com/markbates/pkger/cmd/pkger

**Clone the Repo**

    git clone https://github.com/ericdaugherty/go-lambda-webapp-template newappname

Since you are creating your own project, remove the remote repo reference

    git remote rm origin

Edit go.mod and change the module name to reflect your new project.

**Run Locally**

    make run

This will start the default webserver on http://localhost:3000. You should be able to modify the HTML templates without recompiling the app or restarting the webserver.

## Deploy To AWS

The template includes a [serverless](https://serverless.com/) file to support easy deployment to AWS Lambda, including creating an API Gateway HTTP Endpoint.

To deploy, first make sure [serverless](https://serverless.com/) is installed and configured with your AWS credentials.

Then deploy:

    make deploy

This should create a new CloudFormation stack on AWS including your all-in-one lambda and an REST API Gateway HTTP Endpoint. 

Serverless should provide the HTTP Endpoint URL for you. It should look something like:

    https://<app-id>.execute-api.<region>.amazonaws.com/<stage> 

You will want to add a trailing / to make sure the links work.

## Deploy with Custom Domain

This basic setup has some limitations.
* You must include the stage name as part of the url (ie /dev)
* AWS API Gateway only serves via HTTPS, which is fine, but you can't easily force http -> https
* You are using a generic AWS Domain

All of this can be solved with a few easy steps. The most complete solution is to use CloudFront with a custom domain.

* Deploy the REST API Gateway Endpoint as Regional (Uncomment "endpointType: regional" in serverless.yml and "make deploy")
* Configure a certificate for your custom domain (https://console.aws.amazon.com/acm/home)
* Configure a CloudFront Distribution
  * Origin Domain Name: https://<app-id>.execute-api.<region>.amazonaws.com
  * Origin Path: /<stage name> (ex /dev)
  * Minimum Origin SSL Protocol: TLSv1.2
  * Origin Protocol Policy: HTTPS Only
  * Viewer Protocol Policy: Redirect HTTP to HTTPS (Not required but best practice)
  * Alternate Domain Names (CNAMEs): yourcustomdomain.com
  * SSL Certificate: Choose Custom Certificate and select the certificate you setup for this domain

Once the distribution is setup, go to your DNS provider and setup your custom domain CNAME to point to the domain name of your newly created CloudFront Distribution.

Once the CloudFront distribution is deployed and your DNS entries have propagated, you should be able to access your Lambda web app via your custom domain, and HTTP requests should auto-upgrade to HTTPS.

Note: This example uses the REST API Gateway by default. You can also use the newer HTTP API Gateway by changing the "- http:" lines in serverless.yml to "- httpApi:". For usage as a webapp this may be a better (cheaper) alternative. Please view the [feature comparison](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-vs-rest.html) for details.