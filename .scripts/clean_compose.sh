#!/bin/bash

systemctl start docker
docker compose down
docker volume rm website-health-checker-go_app_data
docker image rm website-health-checker-go-app:latest
