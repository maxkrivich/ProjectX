FROM ubuntu:16.04
MAINTAINER "Maxim Krivich"
RUN mkdir /application
WORKDIR application
COPY main_linux /application
CMD exec ./main_linux