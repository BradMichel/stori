import * as fs from "fs";
import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import { GetProvider, localEndpoint, localStack, Region } from "./local_config";
import aws_sdk = require("aws-sdk");

const env = pulumi.getStack();
const provider = GetProvider(env);
const transactionsFolder = "records";

const emailSender = new aws.ses.EmailIdentity(
  "Michel Barrera",
  { email: "michel.en@hotmail.com" },
  { provider },
);

const emailReceiver = new aws.ses.EmailIdentity(
  "Michel",
  { email: "bradmichel10@gmail.com" },
  { provider },
);

const deploymentBucket = new aws.s3.Bucket(
  `${env}-stori-account-summarizer-bucket`,
  {
    bucket: `${env}-stori-account-summarizer-bucket`,
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

const accountDynamoDB = new aws.dynamodb.Table(
  `${env}-stori-accounts`,
  {
    name: `${env}-stori-accounts`,
    attributes: [{ name: "id", type: "S" }],
    hashKey: "id",
    readCapacity: 5,
    writeCapacity: 5,
  },
  { provider },
);

const accountTransactionsBucket = new aws.s3.Bucket(
  `${env}-stori-account-transactions-bucket`,
  {
    bucket: `${env}-stori-account-transactions-bucket`,
  },
  { provider },
);

new aws.s3.BucketPublicAccessBlock(
  `${env}-stori-account-transactions-bucket`,
  {
    bucket: accountTransactionsBucket.bucket,
    blockPublicAcls: true,
    blockPublicPolicy: true,
    ignorePublicAcls: true,
    restrictPublicBuckets: true,
  },
  { provider },
);

const transactionSeed = new aws.s3.BucketObject(
  `${env}-stori-transaction-seed`,
  {
    bucket: accountTransactionsBucket.id,
    source: new pulumi.asset.FileAsset("../data/1_txns.csv"),
    key: transactionsFolder + "/1_txns.csv",
  },
  { provider },
);

const accountTransactionsDynamoDB = new aws.dynamodb.Table(
  `${env}-stori-account-transactions-table`,
  {
    name: `${env}-stori-account-transactions-table`,
    attributes: [{ name: "id", type: "S" }],
    hashKey: "id",
    readCapacity: 5,
    writeCapacity: 5,
  },
  { provider },
);

const summarizeInitiatorTopic = new aws.sns.Topic(
  `${env}-stori-summarize-initiator-topic`,
  {
    name: `${env}-stori-summarize-initiator-topic`,
  },
  { provider },
);

const summarizeInitiatorSQS = new aws.sqs.Queue(
  `${env}-stori-summarize-initiator-queue`,
  {
    name: `${env}-stori-summarize-initiator-queue`,
    visibilityTimeoutSeconds: 30,
  },
  { provider },
);

const summarizeInitiatorSQSPolicy = new aws.sqs.QueuePolicy(
  `${env}-stori-summarize-initiator-queue-policy`,
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

const summarizeInitiatorTopicSubscription = new aws.sns.TopicSubscription(
  `${env}-stori-summarize-initiator-topic-subscription`,
  {
    topic: summarizeInitiatorTopic.arn,
    protocol: "sqs",
    endpoint: summarizeInitiatorSQS.arn,
  },
  { provider },
);

const summarizeFinishedTopic = new aws.sns.Topic(
  `${env}-stori-summarize-finished-topic`,
  {
    name: `${env}-stori-summarize-finished-topic`,
  },
  { provider },
);

const summarizeFinishedSQS = new aws.sqs.Queue(
  `${env}-stori-summarize-finished-queue`,
  {
    name: `${env}-stori-summarize-finished-queue`,
    visibilityTimeoutSeconds: 30,
  },
  { provider },
);

const summarizeFinishedSQSPolicy = new aws.sqs.QueuePolicy(
  `${env}-stori-summarize-finished-queue-policy`,
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

const summarizeFinishedTopicSubscription = new aws.sns.TopicSubscription(
  `${env}-stori-summarize-finished-topic-subscription`,
  {
    topic: summarizeFinishedTopic.arn,
    protocol: "sqs",
    endpoint: summarizeFinishedSQS.arn,
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
  },
}).apply((config: any) => {
  if (env == localStack) {
    config.conf.endpoint = localEndpoint;
  }
  const yamlConfig = yaml.dump(config);
  fs.writeFile(`../conf.${env}.yml`, yamlConfig, (err) => {
    if (err) {
      console.log(err);
    } else {
      console.log("File created successfully!");
    }
  });
});
