# githubservice

A small service in GO exercising github.com rest API
This is a small service written in Go. 
it uses the net/http library, it marshall/unmarshall json documents, implements a servicelog and it uses unit tests and packages to enforce modularity.
I did it to exercise myself with the Go constructors: net/http, json, routines, packages, files, unit tests.

The goal of this exercise is to write an HTTP service against GitHub’s API with the following
specs:
Service specs
● Given a city name (e.g. Barcelona) the service returns a list of the top contributors
(sorted by number of repositories) in GitHub.
● The service should give the possibility to choose between the Top 50, Top 100 or Top
150 contributors.

Run in your machine:

Assuming the GOPATH, GOROOT are set correctly the project can be compiled in this way:

go install newrelic

and find the newrelic binary inside GOPATH/bin and run from there.

./newrelic

You can test with

curl http://localhost:8080/location=Barcelona+50

Unit tests are provided you can just run:

go test -v ./...

to run the test suite.

A Dockerfile is provided, so you can just:

sudo  docker build -t newrelic .
docker run --publish 6060:8080 --name test --rm newrelic

The process is then running on external port 6060

You can test with

curl http://localhost:6060/location=Barcelona+50

To avoid this is necessary to be authenticated, so I provided a conf file inside this directory where you need to put
a valid github username/password. YOU SHOULD ALSO SPECIFY IF YOU WANT THE APPLICATION BE RUN WITH OR WITHOUT THREADING,
putting yes or not. The conf file is like this
USERNAME:
PASSWD:
MULTITHREADING:no
you can put your GITHUB username and passwd here. Without this information the service will send requests without authentication.
In this way you can do more request to the github api but always limited to 5000 per hour
