FROM python:3.11-slim

WORKDIR /app

RUN pip install -U pip poetry

EXPOSE 8000
