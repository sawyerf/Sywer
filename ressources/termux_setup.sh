#!/bin/sh

mkdir -p /data/data/com.termux/files/usr/var/log/sywer/
mkdir -p /data/data/com.termux/files/usr/var/lib/sywer/

go build ../Sywer.go

mv Sywer /data/data/com.termux/files/usr/bin/sywer
cp termux_settings.swy /data/data/com.termux/files/usr/var/lib/sywer/settings.swy

