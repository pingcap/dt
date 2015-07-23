#!/bin/bash
echo "start a probe"
nohup ../bin/probe &
echo "done"

echo "sleep 3s"
sleep 3
ps -ef|grep probe 
