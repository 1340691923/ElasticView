#!/bin/bash
db_lock_file="/usr/local/db_init.lock"
echo `service mysql status`
echo '启动mysql....'
service mysql restart
sleep 3
echo `service mysql status`
if [ ! -f "$db_lock_file" ]; then
    echo '[ElasticView] init starting...';
    echo '开始导入数据....'
    mysql < /mysql/es_view.sql
    echo '导入数据完毕....'
    sleep 3
    echo `service mysql status`
    echo '开始修改密码....'
    mysql < /mysql/privileges.sql
    echo '修改密码完毕....'
    echo `service mysql status`
    echo `mysql容器启动完毕,且数据导入成功`
    echo '[ElasticView] mysql init success!';
    echo 'LOCK' > "$db_lock_file";
fi


echo '启动ElasticView...'
cd /usr/local && chmod +x ElasticViewLinux && nohup ./ElasticViewLinux &
echo '启动ElasticView完毕！！！'
tail -f /dev/null
