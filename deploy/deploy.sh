#!/bin/bash


SSH_OPTS="-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -oLogLevel=ERROR"

if [ -z $1 ]; then
	echo "Usage: ./deploy [user@hostname]"
	exit
fi

echo "deploy on $1"
HOST=$1

if [ ! -f bin/rdockerskt ];then
  echo "no binary file at bin/rdockerskt"
  exit
fi

ssh $SSH_OPTS $HOST "mkdir -p /tmp/rdockerskt/bin"
scp -r $SSH_OPTS bin rdockerskt rdockerskt.service $HOST:/tmp/rdockerskt

ssh $SSH_OPTS -t $HOST "\
sudo -p '[sudo] password to deploy:' mkdir -p /opt/bin && sudo cp /tmp/rdockerskt/bin/* /opt/bin/ ;\
sudo cp /tmp/rdockerskt/rdockerskt.service /lib/systemd/system/; \
sudo cp /tmp/rdockerskt/rdockerskt /etc/default/; \
sudo systemctl enable rdockerskt.service
sudo service rdockerskt start
"

#test 
echo "Try to test it with curl host:4244/containers/json"




