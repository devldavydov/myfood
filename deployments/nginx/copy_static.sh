#!/bin/bash

sudo rm -rf /var/www/myfood/static
sudo cp -r ./static /var/www/myfood/
sudo chown -R nginx:nginx /var/www/myfood
