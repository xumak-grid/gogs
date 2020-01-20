#!/bin/sh
echo "Running exposer"
if EXPORTED_ENVS=$(/bin/exposer)
then
printenv
eval $EXPORTED_ENVS
sed -i "s/{{domain}}/$ELB_HOSTNAME/g" /data/gogs/conf/app.ini
sed -i "s/{{appname}}/$APP_NAME/g" /data/gogs/conf/app.ini
cat /data/gogs/conf/app.ini
/data/gogs/docker/initial -port 3000 &
/bin/sh /app/gogs/docker/start.sh
else
  echo "Error running exposer" $EXPORTED_ENVS
  exit 1
fi
 