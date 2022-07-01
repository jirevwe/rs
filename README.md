# rs

rs is a zero config mongodb replica set runner

## Install

```
go install github.com/jirevwe/rs
```

## Usage

Start a single node [replica set](https://www.mongodb.com/docs/manual/tutorial/convert-standalone-to-replica-set/) running MongoDB 4.2.21

```sh
$ rs download
$ rs run
```

You can download and run different versions of MongoDB

```sh
$ rs download 4.2.0
$ rs run 4.2.0
```

## Production Use

This tool was not designed to be used to run production database replica set. It was designed for local development and testing alone. If you want to run MongoDB in production and don't want to manage a replica set yourself, use [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).
