#!/bin/bash
echo "start a agent"
nohup ../bin/cmd -role=controller -cfg=../etc/ctrl_cockroach_cfg.toml &
echo "done"

echo "sleep 3s"
sleep 3
tail controller.log -f 
