# Cloudhus

[Cloudhus](https://auth.cloudhus.com/auth)
central authentication, cloud service for Project-Hus which provides various features for lots of fields.

```bash
.
├── makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── services                    <-- Microservices
│   └── hus-auth                <-- Auth service
└── template.yaml              <-- SAM template
```

## Fully Automated Single Sign-On Feature
<div align="center">
  
### [SSO demo page](https://auth.cloudhus.com/auth/demo/sso)

    1. Before you visit any of our services, visit "SSO demo page" above.
    and you're gonna see the page below. 

<img width="300" alt="image" src="https://github.com/Project-Hus/cloudhus/assets/108582413/955bc599-0dd9-4d10-9a05-2e035e534b6d">

    2. Now click the link and visit "Lifthus", and then come back immediately to SSO demo.
    if you refresh the page, you're gonna see this page.
    Which means that your sessions are established and that they are connected to each other.

<img width="300" alt="image" src="https://github.com/Project-Hus/cloudhus/assets/108582413/743a7e23-38df-4180-920f-7b518cb415fe">

    3. Go to Lifthus, and sign in. you will see this, when you come back and refresh.
    Which means that you signed in to two different domains at the same time.

<img width="300" alt="image" src="https://github.com/Project-Hus/cloudhus/assets/108582413/c800d7cf-9037-465b-88aa-7fefddff1026">

    4. Finally, try signing out from Lifthus. you will see the 2nd pic again.
    Which means that you are signed out from two different domains at the same time.

</div>

## Requirements

- [Docker](https://www.docker.com/community-edition)
- AWS CLI with Administrator permission
- SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [Golang](https://golang.org)
- Makefile

## Setup process

### Installing dependencies & building the target

1. "swag init" from hus-auth
2. "make build" from root

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://localhost:3000/`

```bash
make start
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invke your function 'http://localhost:9090'

**Running in native language runtime**

```bash
# hus-auth
make go
```

## Packaging and deployment

Every build and packaging is done by SAM CLI but NestJS must be transpiled to JS before packaging.
Currently this job is done by Makefile. We are migrating to esbuild transpiling which SAM supports natively.

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

- **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
- **AWS Region**: The AWS region you want to deploy your app to.
- **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
- **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modifies IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
- **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

You can find your API Gateway Endpoint URL in the output values displayed after deployment.

### Testing

No tests yet.
