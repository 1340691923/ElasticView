use mysql;
select host, user from user;
-- 因为mysql版本是5.7，因此新建用户为如下命令：
create user worker identified by '1340691923';
-- 将docker_mysql数据库的权限授权给创建的docker用户，密码为123456：
grant all on es_view.* to worker@'%' identified by '1340691923' with grant option;
-- 这一条命令一定要有：
flush privileges;