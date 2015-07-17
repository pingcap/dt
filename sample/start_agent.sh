#!/bin/bash
echo "start a agent"
nohup ../bin/cmd -role=agent -cfg=../etc/agent_cfg.toml &
echo "done"

echo "sleep 3s"
sleep 3
tail agent.log
