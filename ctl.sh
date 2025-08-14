#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

    fyne package -os android -app-id com.example.kaso -icon Icon.png

elif [[ "$1" == "-c" ]]; then

    rm -rf settings/ **/**/settings/

fi
