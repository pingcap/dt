#!/bin/bash
echo "start an agent"
nohup ../bin/agent -role=agent -cfg=../etc/agent_cfg.toml &
echo "done"

echo "sleep 3s"
sleep 3
tail agent.log
