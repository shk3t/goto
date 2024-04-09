#!/bin/bash

BASE_DIR=$PWD
DB_NAME="wb"
USER_NAME="wb_admin"
APPS=(main)

source $BASE_DIR/.venv/bin/activate

sudo -u postgres psql -d $DB_NAME << EOF
    DROP SCHEMA public CASCADE;
    CREATE SCHEMA public;
    GRANT ALL ON SCHEMA public TO $USER_NAME;
EOF

rm -rf $BASE_DIR/*/migrations
python $BASE_DIR/manage.py makemigrations "${APPS[@]}"
python $BASE_DIR/manage.py migrate