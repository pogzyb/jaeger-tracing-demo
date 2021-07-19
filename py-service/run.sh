#!/usr/bin/env sh

uvicorn --reload --factory main:factory --host 0.0.0.0 --port 8081
