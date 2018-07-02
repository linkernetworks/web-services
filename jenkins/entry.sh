#!/bin/bash

mkdir -p /tmp/mongo
nohup mongod --dbpath=/tmp/mongo 2>&1 > /dev/null &

bash
