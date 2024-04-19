#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <server_ip_addr>"
    exit 1
fi

IPADDR=$1

openssl req \
        -x509 \
        -nodes \
        -days 3650 \
        -newkey rsa:2048 \
        -keyout ./myfood-nginx-selfsigned.key \
        -out ./myfood-nginx-selfsigned.crt \
        -subj "/C=RU/ST=Moscow/L=Moscow/O=MyFood/OU=MyFood/CN=${IPADDR}"

openssl dhparam -out ./myfood-dhparam.pem 4096

echo "Run sudo cp ./myfood-nginx-selfsigned.key /etc/ssl/private/myfood-nginx-selfsigned.key"
echo "Run sudo cp ./myfood-nginx-selfsigned.crt /etc/ssl/certs/myfood-nginx-selfsigned.crt"
echo "Run sudo cp ./myfood-dhparam.pem /etc/nginx/myfood-dhparam.pem"
echo "Run sudo cp ./myfood.conf /etc/nginx/conf.d"