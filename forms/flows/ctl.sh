#!/usr/bin/env bash

if [[ "$1" == "-c" ]]; then

	for f in $(ls *.yml); do

		root="$(echo $f | sed 's/.yml//')"
		
		yq -o=json eval "${root}.yml" > "${root}.json"

	done

fi

