#!/bin/bash

cd PROJECT_DIR

sqlite3 myfood.db ".backup 'myfood_$(date +%Y%m%d).db'"
find . -type f -mtime +5 -name 'myfood_*.db' -execdir rm -- '{}' \;
