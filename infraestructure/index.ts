import * as fs from "fs";
import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import { GetProvider, localEndpoint, localStack, Region } from "./local_config";
import aws_sdk = require("aws-sdk");

const stack = pulumi.getStack();
const provider = GetProvider(stack);

const deploymentBucket = new aws.s3.Bucket(
  `${stack}-stori-account-summarizer-bucket`,
  {
    bucket: `${stack}-stori-account-summarizer-bucket`,
  },
  { provider },
);

new aws.s3.BucketPublicAccessBlock(
  `stori-account-summarizer-bucket`,
  {
    bucket: deploymentBucket.bucket,
    blockPublicAcls: true,
    blockPublicPolicy: true,
    ignorePublicAcls: true,
    restrictPublicBuckets: true,
  },
  { provider },
);

export const transactionsFolder = "records";

export const accountTransactionsBucket = new aws.s3.Bucket(
  `${stack}-stori-account-transactions-bucket`,
  {
    bucket: `${stack}-stori-account-transactions-bucket`,
  },
  { provider },
);

new aws.s3.BucketPublicAccessBlock(
  `${stack}-stori-account-transactions-bucket`,
  {
    bucket: accountTransactionsBucket.bucket,
    blockPublicAcls: true,
    blockPublicPolicy: true,
    ignorePublicAcls: true,
    restrictPublicBuckets: true,
  },
  { provider },
);

export const transactionSeed = new aws.s3.BucketObject(
  `${stack}-stori-transaction-seed`,
  {
    bucket: accountTransactionsBucket.id,
    source: new pulumi.asset.FileAsset("../data/1_txns.csv"),
    key: transactionsFolder + "/1_txns.csv",
  },
  { provider },
);

export const accountTransactionsDynamoDB = new aws.dynamodb.Table(
  `${stack}-stori-account-transactions-table`,
  {
    name: `${stack}-stori-account-transactions-table`,
    attributes: [{ name: "id", type: "S" }],
    hashKey: "id",
    readCapacity: 5,
    writeCapacity: 5,
  },
  { provider },
);

export const accountDynamoDB = new aws.dynamodb.Table(
  `${stack}-stori-accounts`,
  {
    name: `${stack}-stori-accounts`,
    attributes: [{ name: "id", type: "S" }],
    hashKey: "id",
    readCapacity: 5,
    writeCapacity: 5,
  },
  { provider },
);

export const summarizeInitiatorTopic = new aws.sns.Topic(
  `${stack}-stori-summarize-initiator-topic`,
  {
    name: `${stack}-stori-summarize-initiator-topic`,
  },
  { provider },
);

export const summarizeInitiatorSQS = new aws.sqs.Queue(
  `${stack}-stori-summarize-initiator-queue`,
  {
    name: `${stack}-stori-summarize-initiator-queue`,
    visibilityTimeoutSeconds: 30,
  },
  { provider },
);

export const summarizeInitiatorSQSPolicy = new aws.sqs.QueuePolicy(
  `${stack}-stori-summarize-initiator-queue-policy`,
  {
    queueUrl: summarizeInitiatorSQS.id,
    policy: {
      Version: "2012-10-17",
      Id: `${summarizeInitiatorSQS.arn}/policy`,
      Statement: [
        {
          Effect: "Allow",
          Principal: "*",
          Action: "sqs:SendMessage",
          Resource: summarizeInitiatorSQS.arn,
          Condition: {
            ArnEquals: {
              "aws:SourceArn": summarizeInitiatorTopic.arn,
            },
          },
        },
      ],
    },
  },
  { provider },
);

export const summarizeInitiatorTopicSubscription =
  new aws.sns.TopicSubscription(
    `${stack}-stori-summarize-initiator-topic-subscription`,
    {
      topic: summarizeInitiatorTopic.arn,
      protocol: "sqs",
      endpoint: summarizeInitiatorSQS.arn,
    },
    { provider },
  );

export const summarizeFinishedTopic = new aws.sns.Topic(
  `${stack}-stori-summarize-finished-topic`,
  {
    name: `${stack}-stori-summarize-finished-topic`,
  },
  { provider },
);

export const summarizeFinishedSQS = new aws.sqs.Queue(
  `${stack}-stori-summarize-finished-queue`,
  {
    name: `${stack}-stori-summarize-finished-queue`,
    visibilityTimeoutSeconds: 30,
  },
  { provider },
);

export const summarizeFinishedSQSPolicy = new aws.sqs.QueuePolicy(
  `${stack}-stori-summarize-finished-queue-policy`,
  {
    queueUrl: summarizeFinishedSQS.id,
    policy: {
      Version: "2012-10-17",
      Id: `${summarizeFinishedSQS.arn}/policy`,
      Statement: [
        {
          Effect: "Allow",
          Principal: "*",
          Action: "sqs:SendMessage",
          Resource: summarizeFinishedSQS.arn,
          Condition: {
            ArnEquals: {
              "aws:SourceArn": summarizeFinishedTopic.arn,
            },
          },
        },
      ],
    },
  },
  { provider },
);

export const summarizeFinishedTopicSubscription = new aws.sns.TopicSubscription(
  `${stack}-stori-summarize-finished-topic-subscription`,
  {
    topic: summarizeFinishedTopic.arn,
    protocol: "sqs",
    endpoint: summarizeFinishedSQS.arn,
  },
  { provider },
);

accountDynamoDB.name.apply(async (tableName) => {
  const client = new aws_sdk.DynamoDB.DocumentClient({
    region: Region,
  });
  const params = {
    region: Region,
    TableName: tableName,
    Item: {
      id: "1",
      accounts: [
        {
          id: "1",
          email: "bradmichel10@gmail.com",
          type: "blue_account",
        },
        {
          id: "1",
          email: "bradmichel10@gmail.com",
          type: "black_card",
        },
        {
          id: "1",
          email: "bradmichel10@gmail.com",
          type: "green_card",
        },
      ],
    },
  };

  try {
    const data = await client.put(params).promise();
    console.log("PutItem succeeded:", data);
  } catch (err) {
    console.error("Unable to add item. Error JSON:", err);
  }
});

export const emailSender = new aws.ses.EmailIdentity(
  "Michel Barrera",
  { email: "michel.en@hotmail.com" },
  { provider },
);

export const emailReceiver = new aws.ses.EmailIdentity(
  "Michel",
  { email: "bradmichel10@gmail.com" },
  { provider },
);

export const templateFolder = "templates";

export const templateBucket = new aws.s3.Bucket(
  `${stack}-stori-templates-bucket`,
  {
    bucket: `${stack}-stori-templates-bucket`,
  },
  { provider },
);

new aws.s3.BucketPublicAccessBlock(
  `${stack}-stori-templates-bucket`,
  {
    bucket: templateBucket.bucket,
    blockPublicAcls: true,
    blockPublicPolicy: true,
    ignorePublicAcls: true,
    restrictPublicBuckets: true,
  },
  { provider },
);

const blueAccountSummaryKey = "stori_account_summary_template.html";
export const blueAccountSummaryTemplateSeed = new aws.s3.BucketObject(
  `${stack}-stori-blue-account-summary-template-seed`,
  {
    bucket: templateBucket.id,
    source: new pulumi.asset.FileAsset("../data/" + blueAccountSummaryKey),
    key: templateFolder + "/" + blueAccountSummaryKey,
  },
  { provider },
);

const blackAccountSummaryKey = "stori_black_summary_template.html";
export const blackAccountSummaryTemplateSeed = new aws.s3.BucketObject(
  `${stack}-stori-black-account-summary-template-seed`,
  {
    bucket: templateBucket.id,
    source: new pulumi.asset.FileAsset("../data/" + blackAccountSummaryKey),
    key: templateFolder + "/" + blackAccountSummaryKey,
  },
  { provider },
);

const greenAccountSummaryKey = "stori_green_summary_template.html";
export const greenAccountSummaryTemplateSeed = new aws.s3.BucketObject(
  `${stack}-stori-green-account-summary-template-seed`,
  {
    bucket: templateBucket.id,
    source: new pulumi.asset.FileAsset("../data/" + greenAccountSummaryKey),
    key: templateFolder + "/" + greenAccountSummaryKey,
  },
  { provider },
);

export const bucketName = deploymentBucket.id;
export const accountDynamoDBName = accountDynamoDB.id;
export const transactionsBucketName = accountTransactionsBucket.id;
export const transactionsDynamoDBName = accountTransactionsDynamoDB.id;
export const summarizeInitiatorName = summarizeInitiatorTopic.id;
export const summarizeInitiatorQueueName = summarizeInitiatorSQS.id;
export const summarizeFinishedName = summarizeFinishedTopic.id;
export const summarizeFinishedQueueName = summarizeFinishedSQS.id;

import yaml = require("js-yaml");

pulumi.Output.create({
  conf: {
    region: Region,
    log_retention: 5,
    email_sender: emailSender.email,
    deployment_bucket: deploymentBucket.id,
    account_dynamodb_table: accountDynamoDB.arn,
    account_dynamodb_table_name: accountDynamoDB.name,
    account_transactions_bucket: accountTransactionsBucket.id,
    account_transactions_folder: transactionsFolder,
    account_transactions_dynamodb_table: accountTransactionsDynamoDB.arn,
    account_transactions_dynamodb_table_name: accountTransactionsDynamoDB.name,
    account_summarize_initiator_topic: summarizeInitiatorTopic.arn,
    account_summarize_initiator_topic_name: summarizeInitiatorTopic.name,
    account_summarize_initiator_queue: summarizeInitiatorSQS.arn,
    account_summarize_finished_topic: summarizeFinishedTopic.arn,
    account_summarize_finished_topic_name: summarizeFinishedTopic.name,
    account_summarize_finished_queue: summarizeFinishedSQS.arn,
    templates_bucket: templateBucket.id,
    templates_folder: templateFolder,
  },
}).apply((config: any) => {
  if (stack == localStack) {
    config.conf.endpoint = localEndpoint;
  }
  const yamlConfig = yaml.dump(config);
  fs.writeFile(`../conf.${stack}.yml`, yamlConfig, (err) => {
    if (err) {
      console.log(err);
    } else {
      console.log("File created successfully!");
    }
  });
});
