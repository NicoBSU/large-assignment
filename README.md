# Large Assignment (Minio Gateway)

This project is a gateway to multiple minio instances which is used
to load files into minio instance and read them. 
Files are loaded evenly between instances, application is stateless.
All minio instances, their access and secret keys, IP addresses are read dynamically from docker daemon

ID's are being validated to be aphanumeric. Existing file is being overwritten

This README provides instructions on how to run and interact with the application.

## Table of Contents

- [Project structure](#project-structure)
- [Building the Application](#building-the-application)
- [Making Requests](#making-requests)


## Project structure

Project is divided into two layers: handlers and services.
Router has two endpoinds (**[GET] /object/{id}** and **[PUT] /object/{id}**). Both of them are wrapped with validation middleware which assures that there is an id passed as a parameter and that it is alphanumeric

Docker service is responsible for initializing Docker client, getting Minio Containers, their environment variables, their IPs and getting full config list for all of minio instances

Minio service is responsible for Creating Buckets on app initialization, putting and retrieving objects

Minio service manager is responsible for creating minio service instances, creating minio clients and getting appropriate minio
service by object id

Application configs are stored in config.yaml and read by viper library

Zap is used for logging
Gorilla mux is used as a router

## Building the Application

Run this command to rebuild app image and run all containers

```bash
docker compose up --build
```

## Making requests

To interact with the API, you need to make HTTP requests to the following endpoints. I've added postman jsons, which you might import
to your postman, to run requests.  Those jsons are located in .extras/postman folder
