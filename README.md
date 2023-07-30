# Bullshape API
The Bullshape Web Service is a RESTful API that allows users to generate and retrieve company names. It provides endpoints to create new companies, fetch , update and delete existing ones. Only authenticated users have access to create, update and delete companies.


## Get started
# Prerequisites
- Kafka/Zookeeper (v2.12-3.5.1)
- Mysql(Ver 8.0.33 or higher)

## Project structure
### cmd
It containes the main function
### confs
It reads the bullshape-api.conf file to load the needed parameters for starting the service. Config file is a Toml file and it should be placed at the same directory as the execution of the service
### ctrls
It is responsible for the transport level, such as request validation, marshalling a request into an object or unmashalling an object to feed to models.
### db
It is the permanent store and communicates with the mysql  database for storing the users and companies data.

### models
This folder contains all the systems models. All bussiness logic is permormed in the models. This folder contains the integration tests integration.

### router
Authentication and service routes are created. JWT auth is implemented.

# utils 
A set of utilities is placed here. Kafka for sending events ,  logging mechanism and http wrappers are created.

# docker-compose.yml
All services (Kafka, Zookeeper, Mysql and bullshape ) are dockerized. To start all services run:
```
sudo docker-compose up -d --build
```
To stop all services run:
```
sudo docker-compose down
```
# Feature work
 - Use prometheus for monitoring and alerting 
 - Change project structure to DDD
 - Better use of Kafka event producer. Aim the performance
 - Write test for users. 
 - write e2e tests
 - Create a makefile 
 - Change configuration to yaml. Use go standards libraries. Parameterize Docker and bullshape service with the same yaml file.
 

