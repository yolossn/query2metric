# query2metric

A tool to run db queries in defined frequency and expose the value as prometheus metric

## Why ?

Product metrics play an important role in understanding product adoption and historic metrics helps answer many questions about a product (for eg: which week of the day do I get the most signups). One common things is that most of these metrics are extracted by querying the databases. The tool takes queries and time frequency as configuration and runs the queries in the specified intervals and exposes the output as prometheus metrics.

## Example

Create a config.yaml file.

> config.yaml

```yml
connections:
  - name: mongodb1
    type: MONGO
    connectionStringFromEnv: MONGO_CONN
    metrics:
      - name: active_user_count
        helpString: users in the product
        database: test
        collection: test
        query: '{"is_active":true}'
        time: 10
      - name: total_user_count
        helpString: users in the product
        database: test
        collection: test
        query: ""
        time: 120
  - name: postgres1
    type: SQL
    connectionStringFromEnv: POSTGRES_CONN
    metrics:
      - name: template_count
        helpString: products in the db
        query: select * from templates
        time: 2
      - name: active_template_count
        helpString: products in the db
        query: error
        time: 4
```

Along with the metrics defined, the success and failure count of queries are also exposed as prometheus metrics.
`query2metric_success_count` - No of successful queries coverted to metrics
`query2metric_error_count` - No of errors when converting query to metric

## How to use ?

At present the tool supports mongo and sql queries. Just create a config.yaml

### Mongo

set `type` as `MONGO` and metrics as given in example with `query`,`time` (in seconds) etc.

```yml
connections:
- name: mongodb1
    type: MONGO
    connectionStringFromEnv: MONGO_CONN
    metrics:
      - name: active_user_count
        helpString: users in the product
        database: test
        collection: test
        query: '{"is_active":true}'
        time: 10
```

### SQL

set `type` as `SQL` and metrics as give in example.

```yml
connections:
  - name: postgres1
    type: SQL
    connectionStringFromEnv: POSTGRES_CONN
    metrics:
      - name: template_count
        helpString: products in the db
        query: select * from templates
        time: 2
```

## Run example using docker

You can run the example along with prometheus and grafana using docker.

docker-compose.yaml

> docker-compose up

Output

metrics output: [localhost:8090/metrics](localhost:8090/metrics)
prometheus dashboard: [localhost:9090/graph](localhost:9090/graph)
grafana dashboard: [http://localhost:3000/d/qqTN2unMk/example?orgId=1](http://localhost:3000/d/qqTN2unMk/example?orgId=1)

<p align="center">
  <img width="720" height="354" src="https://github.com/yolossn/query2metric/blob/docker_example/images/grafana.png">
</p>
