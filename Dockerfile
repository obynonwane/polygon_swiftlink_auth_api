# base go image - for production 
# FROM alpine:latest
FROM --platform=linux/amd64 alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp" ]