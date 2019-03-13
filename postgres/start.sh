#!/bin/bash

sudo docker run --net=host dpsql postgres -c config_file=postgresql.conf
