# youpassbutter [![Go Report Card](https://goreportcard.com/badge/github.com/stuartaroth/youpassbutter)](https://goreportcard.com/report/github.com/stuartaroth/youpassbutter)

**youpassbutter** is PostgreSQL data service with a simple premise: Write all your queries in one place and get the results anywhere as JSON by using HTTP POST calls.

There are several reasons why you should consider using **youpassbutter**

###### Rewrite costs

When you inevitably need to do a rewrite into another language or framework you won't have to rewrite all your queries as well. They'll be available right away through HTTP calls.


###### ORM avoidance

ORMs love to come up with "clever" syntax for their users. While it may be interesting to learn these there's one API that will work with any language or framework: pure SQL. Your queries will be more efficient and concise.

###### Single source of knowledge and access

By keeping your queries in one place you will better be able to control your data. **youpassbutter** will act as the gatekeeper. It also eases testing when something goes wrong. Users can quickly copy and run the queries instead of taking time to figure out what's being run.

###### Microservice friendly

Microservice architecture is getting more and more popular. You can cut down on development time by offloading the data access to **youpassbutter**. Worry about what your service does with the data instead of how it gets the data.


#### Get the server started

Before you get started you will need to install PostgreSQL and go.

You will also need to install the following dependencie(s):

```shell
go get github.com/lib/pq
```

Once compiled youpass butter is easy to use.

```shell
youpassbutter -c config.json -q queries.json
```

Your config.json and queries.json should follow the examples provided in the project.

To clone the repo and get the examples set up run the following commands:

```shell
git clone https://github.com/stuartaroth/youpassbutter.git
cd youpassbutter/
psql < example_schema.sql
```

Modify the example_config.json according to your database access needs.

Then run the following commands to get the server running:

```shell
./build
./youpassbutter -c example_config.json -q example_queries.json
```

In another terminal run the setup:
```shell
./example_setup.sh
```

Now you can use the queries in your example_queries.json to access the data.

#### Access the data

Access calls follow a simple format:

POST http://localhost:8080?q=queryNameInJson
JSON [FirstParameter, SecondParameter, ThirdParameter]

Queries are written as PostgreSQL prepared statments. The first item in the JSON array paramter will replace the $1, the second will replace the $2 and so on.

Select results will come back as an array of JSON objects.

Other requests will either have a JSON object with a "Success" or an "Error" and an associated message.

#### Questions / issues

This project is still in development and therefore not perfect. Please report any issues or questions you have through the issues tab of the repo. I would love your feedback!
