#!/usr/bin/env bash

if [[ "$1" == "-b" ]]; then

    fyne package -os android -app-id com.example.kaso -icon Icon.png

elif [[ "$1" == "-bc" ]]; then

    go build -o cli cmd/wscli/*.go

elif [[ "$1" == "-c" ]]; then

    rm -rf settings/ **/**/settings/ **/**/data/ **/**/*.db **/tmp*/

elif [[ "$1" == "-cov" ]]; then

    go test -coverprofile=coverage.out ./...

    go tool cover -html=coverage.out

fi
