#!/bin/bash
echo `service mysql status`
echo '1.启动mysql....'
service mysql start
sleep 3
echo `service mysql status`
echo '2.开始导入数据....'
mysql < /mysql/es_view.sql
echo '3.导入数据完毕....'
sleep 3
echo `service mysql status`
echo '4.开始修改密码....'
mysql < /mysql/privileges.sql
echo '5.修改密码完毕....'
echo `service mysql status`
echo `mysql容器启动完毕,且数据导入成功`
echo '6.启动ElasticView...'
cd /usr/local && chmod +x ElasticViewLinux && nohup ./ElasticViewLinux &
echo '6.启动ElasticView完毕！！！'
tail -f /dev/null
