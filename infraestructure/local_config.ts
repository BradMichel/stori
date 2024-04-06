import * as aws from "@pulumi/aws";

export const localStack = "local";
const providerName = "local-provider";
export const localEndpoint = "http://localhost:4566";
export const Region = "us-east-1";
const credential = "local";

export const GetProvider = (stack: string): aws.Provider => {
  let provider = new aws.Provider("my-provider", { region: Region });
  if (stack === localStack) {
    // @ts-ignore
    provider = new aws.Provider(providerName, {
      region: Region,
      accessKey: credential,
      secretKey: credential,
      skipCredentialsValidation: true,
      skipRequestingAccountId: true,
      s3UsePathStyle: true,
      endpoints: [
        {
          dynamodb: localEndpoint,
          s3: localEndpoint,
          sns: localEndpoint,
          sqs: localEndpoint,
          ses: localEndpoint,
        },
      ],
    });
  }

  // @ts-ignore
  return provider;
};
