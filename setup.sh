#!/bin/sh

go build Sywer.go
mv Sywer /usr/sbin/sywer
cp ressources/sywer.service /etc/systemd/system/.
mkdir -p /var/log/sywer/
mkdir -p /var/lib/sywer/
