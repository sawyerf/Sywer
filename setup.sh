#!/bin/sh


go build Sywer.go
mv Sywer /usr/sbin/sywer
cp ressources/sywer.service /etc/systemd/system/.