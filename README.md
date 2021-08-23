[![Integration Test Status](https://github.com/jrkhan/qframe_integration_test/actions/workflows/integration-test.yml/badge.svg)](https://github.com/jrkhan/qframe_integration_test/actions/workflows/integration-test.yml
)

Uses Github Actions to perform integration tests for [QFrame](https://github.com/tobgu/qframe).

Presently, this consists of ensuring that QFrame will be able to build dataframes from each supported database (for now MySQL, Postgres, and SQLite).

### Local Usage
Update go.mod to the version of the QFrame library you would like to test.


```
module github.com/jrkhan/qframe-integration-test

go 1.16

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/lib/pq v1.10.2
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/tobgu/qframe v0.3.6
)

replace github.com/tobgu/qframe => ../../tobgu/qframe // local changes
```
Then `docker-compose up` and run tests through the CLI or IDE of your choice.