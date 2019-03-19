#!/bin/bash

if [ -d "./nginx/frontend_dir" ]; then
    echo "[Delete previous directory]"
    sudo rm -rf nginx/frontend_dir
fi

if [ "$1" == "-u" ]; then
    arg="--no-cache"
    echo "[Mode: Full update build]"

    echo "[Building frontend part]"


    git clone -b dev https://github.com/frontend-park-mail-ru/2019_1_Auteam.git ./nginx/frontend_dir
    sudo docker run -v $(pwd)/nginx/frontend_dir:/public node bash -c "cd /public && npm link cors && npm install && npm run hand"

    git pull origin RouterDev 
else
    arg=""
    echo "[Mode: Cache data build]"
fi

echo "[Building backend part]"
sudo docker-compose stop && \
sudo docker-compose build $arg && \
sudo docker-compose up