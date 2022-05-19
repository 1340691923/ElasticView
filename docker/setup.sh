#!/bin/bash

# echo `service mysql status`
# echo '1.启动mysql....'
# service mysql start
# sleep 3
# echo `service mysql status`
# echo '2.开始导入数据....'
# mysql < /mysql/es_view.sql
# echo '3.导入数据完毕....'
# sleep 3
# echo `service mysql status`
# echo '4.开始修改密码....'
# mysql < /mysql/privileges.sql
# echo '5.修改密码完毕....'
# echo `service mysql status`
# echo `mysql容器启动完毕,且数据导入成功`
# echo '6.启动ElasticView...'
# cd /usr/local && chmod +x ElasticViewLinux && nohup ./ElasticViewLinux &
# echo '6.启动ElasticView完毕！！！'
# tail -f /dev/null

db_lock_file="/db_init.lock"
es_view_path="/usr/local/ElasticViewLinux"

echo '[ElasticView] mysql service starting...';
service mysql start;
sleep 1;

if [ ! -f "$db_lock_file" ]; then
    echo '[ElasticView] init starting...';
    sleep 3;
    mysql < /mysql/es_view.sql;
    sleep 1;
    mysql < /mysql/privileges.sql;
    echo '[ElasticView] mysql init success!';
    echo 'LOCK' > "$db_lock_file";
    chmod +x "$es_view_path";
    echo '[ElasticView] init success!';
fi

nohup "$es_view_path" &

tail -f /dev/null
