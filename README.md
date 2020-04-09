# Jor-El - Master of Scheduling on the Intergalactic Stock Exchange

Jor-El is the service extends Google Pub/Sub with some new features:

- Routing - the message can be republished in the different topics based on their meta information.

- Scheduling - the message can be republished with some delay.

- Debounce - messages with the same *aggregation-id* can be *debounced*. It means, will be republished 
ONLY the last message (or first)  (not implemented yet).

- Aggregation - messages with the same *aggregation-id* can be *joined*. It means, will be republished 
one *aggregated* message contains an array of original messages (not implemented yet).

## How to run the test project

1) create GCP project (or use existing) and download the service account key.

2) you need the following topics and subscriptions:
* **ingress-topic** topic with **ingress-subs** subscription;
* **default-topic** topic with **default-topic-subs** subscription;
* **cancel-topic** topic with **cancel-topic-subs** subscription;

3) activate Cloud SQL wth PostgresSQL 11 (documentation https://cloud.google.com/sql/docs/postgres/quickstart)

4) create two databases: one for the production service, one for the tests

5) create tables using **create_tables.sql** script

6) set environment variables:

* GOOGLE_APPLICATION_CREDENTIALS = path_to_the_account_key

* JOR_EL_POSTGRE_URL - production database connection string 
(something like *postgresql://user:password@host:port/database*)

* JOR_EL_POSTGRES_TEST_URL - test database connection string

7) run tests *go test ./...*

8) run test receiver *make run-receiver*

9) run jor-el *make run-jorel* (or run several jorel instances in different terminals windows)

10) run test publisher *make run-publisher*

## Jor-El config

**config.yaml**

`
---
project-id: jorel-test-project
ingress-subscription: ingress-subs

default:
  topic-name: default-topic
  aggregation: no

routing:
  - message-type: cancel
    topic-name: cancel-topic
    aggregation: no
`

The first two lines are pretty obvious: GCP project name and jor-el ingress subscription ID. 

**default** section is the configuration for the jor-el default egress topic. jor-el handles all messages 
without message type or with unknown message type according to this section rules. It's a mandatory section.

**routing** - section with routing rules. Optional.

## Ingress message format

The ‘Meta Message’ (or jor-el ingress message) wraps the ‘Target Message’ as jor-el only needs 
the bytes that will be sent and the name of the queue, the ‘Target Queue’. 
This way jor-el doesn’t need to know anything about the content of that ‘Target Message’.  

Meta information is passing through message attributes:

* **delay** - message delay in seconds. Mandatory attribute. It can be 0, it means - republish immediately.
* **message-type** - using for the routing (if needed). Optional.
* **aggregation-id** - using for the debounce and aggregation. Optional. Not implemented yet.



