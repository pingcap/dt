#!/bin/bash
echo "start a controller"
nohup ../bin/ctrl -role=controller -cfg=../etc/ctrl_example_cfg.toml &
echo "done"

echo "sleep 3s"
sleep 3
tail controller.log
