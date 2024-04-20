#!/bin/bash

psql << EOF
    DROP DATABASE goto;
    CREATE DATABASE goto;
EOF
