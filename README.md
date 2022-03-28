  

技术栈主要用： Vue + golang 

欢迎大家提出自己的issue。

ElasticView
-----------
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/1340691923/ElasticView)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/1340691923/ElasticView/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/1340691923/ElasticView.svg?label=Release)](https://gitee.com/cynthia520/elastic-view/releases)
> ElasticView 是一款用来监控elasticsearch状态和操作elasticsearch索引的web可视化工具。它由golang开发而成，具有部署方便，占用内存小等优点
 * es连接树管理（更方便的切换测试/生产环境）
 * 支持页面级别的权限管理（未来将完善）
 * 支持sql转换成dsl语法
 * 更方便的重建索引
 * 任务管理
 * 备份管理
 * 可将查询内容下载为excel文件
 * 可进行索引创建，映射创建，别名创建，索引删除等操作
 * 支持版本 6.x, 7.x,8.x
 * 支持Navicat界面 （需要前往权限-角色界面 重新分配界面权限后刷新页面就可以看到）
 * docker部署
 
 
## ElasticView 部分截图

![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/1.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/2.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/3.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/4.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/5.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/6.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/7.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/8.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/9.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/10.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/11.png)
![image](https://raw.githubusercontent.com/1340691923/ElasticView/main/show_img/12.png)


## 应用程序下载
[下载地址]( https://gitee.com/cynthia520/elastic-view/releases/) 

##支持操作系统：
1. windows
2. linux
3. mac

## 安装教程
 1. 下载release里面的对应压缩包（windows用户下载ElasticView_windows.zip，linux用户下载ElasticView_linux.zip，mac用户下载 ElasticView_mac.zip）
 2. 解压后，新建mysql数据库 es_view，导入es_view.sql 
 3. 修改config.json文件中的 数据库连接信息，日志存放目录和应用启动端口等配置信息
 4. windows：双击ElasticView.exe  linux：chmod +x ElasticView && nohup ./ElasticView > ElasticView.log &
 5. 浏览器访问对应ip:端口，初始用户名：admin，初始密码：admin
 
## docker部署
 1. docker pull 1340691923/elastic_view:latest
 2. docker run -d -p 8090:8090 1340691923/elastic_view:latest
 3. 浏览器访问对应ip:8090，初始用户名：admin，初始密码：admin

## 手动编译
 1. git clone git@github.com:1340691923/ElasticView.git
 2. cd vue && cnpm install (安装前端依赖)
 3. npm run build:prod （打前端正式包）
 4. windows: go build -o ElasticView.exe  linux : 双击build_linux.bat // 编译go可执行程序
 
## 其他开源项目
 * 铸龙用户行为分析系统     https://github.com/1340691923/xwl_bi
 * 软考成绩快查工具               https://github.com/1340691923/SoftTestMonitor
 
