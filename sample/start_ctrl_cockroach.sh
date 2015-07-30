#!/bin/bash
count=$1

i=0
while ((i<$count))
do
	echo "start a controller, no."$i
	nohup ../bin/ctrl -role=controller -cfg=../etc/ctrl_cockroach_cfg.toml &
	sleep 3
	tail controller.log -f 

	sleep 100
	echo "done"
	let ++i
done

echo "end"
