#!/usr/bin/env bash
for input in `ls *.json`
do
	name=`echo ${input} | awk -F '.' '{print $1}'`
	output="../${name}.go"
	echo ${input} to ${output}
	go run maker.go -c ${input} -o ${output}
done