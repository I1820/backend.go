# WebBackend
[![Travis branch](https://img.shields.io/travis/com/I1820/backend/master.svg?style=flat-square)](https://travis-ci.com/I1820/backend)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4407720668f84f119b15c653c84fa1f2)](https://www.codacy.com/app/i1820/backend?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=I1820/backend&amp;utm_campaign=Badge_Grade)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/backend?style=flat-square)](https://goreportcard.com/report/github.com/I1820/backend)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/backend)


## Introduction

Each platform with microservice architecture has a web backend component. It just like a glue and tries to put everything
together in a simple fashion.
I1820 backend, it provides a single interface to all services.
It adds user management and authorization into the platform.
And at the end it hosts our web application that is written in Angular with :heart: and memory of someone who is not among us.

## Notes that are worth taking
Users must have an **unique username** in I1820 platform.
Backend heavily uses JWT tokens for each of its works.

### Operations need token refreshing
1. Create project

## Up and Running
To build this module from source do the following steps

1. Make sure MongoDB is up and running.

2. Install the required dependencies (Please note that we use [dep](https://github.com/golang/dep) as our go package manager)
```sh
dep ensure
```

3. Check the configuration in `.env` file. (You can use `.env.example` as an example configuration).

4. Run :runner:
```sh
go build
./backend
```

5. Create MongoDB indexes
```sh
buffalo task mongo
```

