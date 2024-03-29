# Go Lambda WebApp Template
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

This repository is a template that can be used to start building a self contained web-application written in Go and deployed on AWS Lambda.

This template makes use of minimal external components while providing a ready-to-use template to develop an all-in-one Lambda web application.

This template makes use of:
- [Apex Gateway](https://github.com/apex/gateway) to map between Lambda requests and standard net.http requests.
- [CHI](github.com/go-chi/chi) to route incoming requests.

## Usage

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

## Deploy using CloudFormation with a Custom Domain

This basic setup has some limitations.
* You must include the stage name as part of the url (ie /dev)
* AWS API Gateway only serves via HTTPS, which is fine, but you can't easily force http -> https
* You are using a generic AWS Domain

All of this can be solved with a few easy steps. The most complete solution is to use CloudFront with a custom domain.

* Deploy the REST API Gateway Endpoint as Regional (Uncomment "endpointType: regional" in serverless.yml and "make deploy")
* Configure a certificate for your custom domain (https://console.aws.amazon.com/acm/home)
* Configure a [CloudFront Distribution](https://console.aws.amazon.com/cloudfront/home)
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

## Deploy using API Custom Domain

Note: This approach only supports https, not http or http->https

* Configure a certificate for your custom domain (https://console.aws.amazon.com/acm/home)
* Navigate to the [API Gateway Console](https://console.aws.amazon.com/apigateway/)
* Select 'Custom domain names' and click "Create"
  * Enter the full name of the domain you created a certificate for
  * Leave the Endpoint and TLS Defaults
  * Select your certificate from the drop-down
  * Select 'Create'
  * On the next page, select 'Configure API Mappings'
  * Select 'Add new mapping'
  * Map your API Endpoint to the domain
* Update your DNS records to point your custom CNAME to the 'API Gateway domain name' value.