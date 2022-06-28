#!/bin/bash

FILE=banco.zip
if [ -f "$FILE" ]; then
    rm -rf banco.zip ./masterInput/* ./masterOutput/*
    echo "Master limpo para execucao"
else 
    echo "Master limpo para execucao"
fi
