#!/bin/bash

while true
do
  curl http://gateway
  curl http://gateway/propagate1
  curl http://gateway/propagate2
  curl http://gateway/propagate3

  sleep 30
done
