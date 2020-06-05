#!/bin/bash
docker-compose down -v
docker-compose up -d db
docker build -t orders-api .
docker-compose up -d app