# **Group-24-ECE461**

Group 24's implementation of a package rater written in Golang for ECE 461.

Group:

- Tony Ni @ ni86@purdue.edu
- Yigitkan Balci @ yblaci@purdue.edu
- Varun Dengi @ vdengi@purdue.edu

# **Components**

# CMD

# API

# Internal

## CLI

## Config

- For configuration files we used the cobra package.
- All environmental variables are accessed through a cofig struct that get initialized using the NewConfig() function.
- In main.go there is a check to see if GITHUB_TOKEN is set if it not set the program will terminate with an error.

## Error

- Our project consists of 2 implementations error structs for handling errors: Request Errors and General Errors.

- Request Error will only appear in the API when sending out HTTP request to the package host.

- General Error will appear through out the code base whenever an error is thrown.

### Request Error

- Request Error = {
  StatusCode: int,
  Message: string,
  RequestType: string
  }
- The Status Code represents the returned HTTP status code.
- The Message contains the Error message.
- The Request Type is what kind of API request type i.e. (RESTful or GraphQL).

### General Error

- General Error = {
  Function: string,
  Message: string
  }
- The Function represents when function caused the error.
- The message represents the Error message.

## Logger

- We used the Zap library to build our logger.
- Zap provides fast peformance while also giving the option for log levels.
- The logger will be initialized in main.go and a single logger object will be passed around.

### LOG_LEVEL

- The default LOG_LEVEL is set to silent, but you can change this by setting LOG_LEVEL in bashrc.
- LOG_LEVEL=1 means debug.
- LOG_LEVEL=2 means info.
- default is silent.

### LOG_FILE

- the LOG_FILE enviromental variable determines where the log file will be output to.
- Without setting a value for LOG_FILE, logger will create a "mylog.log" in the cmd/ directory.

## Models

## Parser

## Scorer

# Test
