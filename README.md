# Jor-El - Master of Scheduling on the Intergalactic Stock Exchange

![](/small-jor-el.jpg)

Jor-El is the service extends Google Pub/Sub with some new features:

- Routing - the message can be republished in the different topics based on their meta information.

- Scheduling - the message can be republished with some delay.

- Debounce - messages with the same *aggregation-id* can be *debounced*. It means, will be republished 
ONLY the last message (or first)  (not implemented yet).

- Aggregation - messages with the same *aggregation-id* can be *joined*. It means, will be republished 
one *aggregated* message contains an array of original messages (not implemented yet).

![](/common-schema-big.png)

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

* JOR_EL_POSTGRES_URL - production database connection string 
(something like *postgresql://user:password@host:port/database*)

* JOR_EL_POSTGRES_TEST_URL - test database connection string

7) run tests *go test ./...*

8) run test receiver *make run-receiver*

9) run jor-el *make run-jorel* (or run several jorel instances in different terminals windows)

10) run test publisher *make run-publisher*

## Jor-El config

**config.yaml**

```
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
```

The first two lines are pretty obvious: GCP project name and jor-el ingress subscription ID. 

**default** section is the configuration for the jor-el default egress topic. jor-el handles all messages 
without message type or with unknown message type according to this section rules. It's a mandatory section.

**routing** - section with routing rules. Optional.

## Ingress message format

The ‘Meta Message’ (or jor-el ingress message) wraps the ‘Target Message’ as jor-el only needs the bytes that will be sent and the name of the queue, the ‘Target Queue’. This way jor-el doesn’t need to know anything about the content of that ‘Target Message’.  

Meta information is passing through message attributes:

* **delay** - message delay in seconds. Mandatory attribute. It can be 0, it means - republish immediately.
* **message-type** - using for the routing (if needed). Optional.
* **aggregation-id** - using for the debounce and aggregation. Optional. Not implemented yet.

## Run in GKE

1) Create GKE cluster
2) Install secrets:
```
    kubectl create secret generic account-credentials --from-file=service-key.json=service-key.json
    kubectl create secret generic db-credentials --from-literal=connection=<DB_CONNECTION_STRING>
```

Build Docker image & generate deployment:
```
    ./scripts/generate_deploy.py
```

Deploy Jor-El on Kubernetes (deployment is being generated on the previous step):
```
    ./run_deployment.sh
```

Check deployment, view logs, scale if needed:
```
    kubectl get pods
    kubectl logs -f jorel-deployment-6f8dd87f4f-msp5q jorel // you have to specify container because two containers deployed in one pod
    kubectl scale deployment.v1.apps/jorel-deployment --replicas=3
```

## Test results

Scale cluster into 5 pods:
```
    kubectl scale deployment.v1.apps/jorel-deployment --replicas=5
```

Run two receivers locally in different terminals:
```
    make run-receiver // receiver for the default topic
    make run-receiver2 // receiver for the cancel topic
```

Open /cmd/publisher/main.go and update tests loop count:
```
    for repeat := 0; repeat < 500; repeat++ { // 100 -> 500
		for indx, d := range delays {
```

Run publisher and wait about 20 minutes:
```
    make run-publisher2 // publish messages with two different types
```

### Test results

Publisher & receivers are being run locally, so local configuration is important:
MacBook Pro 2,3 GHz 8-Core Intel Core i9 32 GB 2400 MHz DDR4

Publisher sent in 5000 messages in 8m 21s.

Receiver 1:
* received messages 2530, max deviation 3.13s, mean deviation 1.53s

Receiver 2:
* received messages 2470, max deviation 5.27s, mean deviation 1.53s

All messages were delivered, the scheduling accuracy is acceptable. 
To be honest it's more Pub/Sub accuracy.