#!/bin/sh

mkdir -p /var/log/sywer/
mkdir -p /var/lib/sywer/

go build Sywer.go
chmod +x ressources/sywer


mv Sywer /usr/sbin/sywer
# cp ressources/sywer /usr/sbin/sywer
cp ressources/sywer.service /etc/systemd/system/.
