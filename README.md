<h1 align="center">query2metric</h1>
<p align="center">A tool to run db queries in defined frequency and expose the count as prometheus metrics.</p>
<p align="center">
    <img src="https://github.com/yolossn/query2metric/blob/master/images/gopher.png" height="200px"/>
</p>

## Why ?

Product metrics play an important role in understanding product adoption and historic metrics helps answer many questions about a product (for eg: which day of the week do I get the most signups). One common thing is that most of these metrics are extracted by querying the databases. The tool takes queries and time frequency as configuration and runs the queries in the specified intervals and exposes the output as prometheus metrics.

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

`query2metric_success_count` - No of successful queries coverted to metrics.

`query2metric_error_count` - No of errors when converting query to metric.

Note: Errors can occur due to invalid queries or connection issues to the db, one can use the logs to debug the issues.

## How to use ?

At present the tool supports mongo and sql queries. Just create a config.yaml file and run the code.

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

metrics output: [localhost:8090/metrics](http://localhost:8090/metrics).

prometheus dashboard: [localhost:9090/graph](http://localhost:9090/graph).

grafana dashboard: [localhost:3000/d/qqTN2unMk/example?orgId=1](http://localhost:3000/d/qqTN2unMk/example?orgId=1).

Example Output:

<p align="center">
  <img width="720" height="354" src="https://github.com/yolossn/query2metric/blob/master/images/grafana.png">
</p>

## Credits

- Logo credit [gopherize.me](gopherize.me)
