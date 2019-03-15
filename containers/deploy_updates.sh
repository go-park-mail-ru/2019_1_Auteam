#!/bin/sh

git pull origin RouterDev && \
sudo docker-compose stop && \
sudo docker-compose build --no-cache && \
sudo docker-compose up
