
<div align=center>
<img src="https://raw.githubusercontent.com/1340691923/ElasticView/c01b67cf1f97fb543d4513d1b6a4a7eac20a8387/resources/vue/src/assets/logo.png" width="300" height="300" />
</div>
<div align=center>
<img src="https://img.shields.io/badge/golang-1.23-blue"/>
<img src="https://img.shields.io/badge/gin-1.10-lightBlue"/>
<img src="https://img.shields.io/badge/vue-3.4.31-brightgreen"/>
<img src="https://img.shields.io/badge/element--plus-2.7.6-green"/>
<img src="https://img.shields.io/badge/gorm-1.25.7-red"/>
</div>

[英文](./README.md) | 简体中文

# 项目文档
[官网](http://www.elastic-view.cn)

[视频教程](https://www.bilibili.com/video/BV12tDDYWEP2/?vd_source=d03eb2249d8310afce3f5b90c6081bb3)

[交流社区](https://txc.qq.com/products/666253)


# 重要提示

0.由于码云最大只能上传100m，所以从0.0.20后的最新版本只能从[官网](http://www.elastic-view.cn/deploy/)下载

1.本项目从起步到开发到部署均有文档和详细视频教程

2.二开本项目需要您有一定的golang和vue3基础

3.您完全可以通过我们的教程和文档完成一切操作，因此我们不再提供免费的技术服务，如需服务请进行[加作者付费支持](https://raw.githubusercontent.com/1340691923/ElasticView/main/resources/show_img/weixin.jpg)

4.您可通过[交流社区](https://txc.qq.com/products/666253)来进行反馈需求和bug，谢谢

## 1. 基本介绍

### 1.1 项目介绍

> ElasticView是一个基于 [vue](https://vuejs.org) 和 [gin](https://gin-gonic.com) 开发的全栈前后端分离的数据源管理插件平台，集成jwt鉴权，动态路由，动态菜单，casbin鉴权，数据源管理，插件市场等功能。


## 2. 主要功能

- 权限管理：基于`jwt`和`casbin`实现的权限管理。
- 用户管理：系统管理员分配用户权限组和权限组权限。
- 权限组管理：创建权限控制的主要对象，可以给权限组分配不同api权限和菜单权限。
- 数据源管理：可进行设置需要管理的数据源，已集成 elasticsearch(6,7,8),mysql,redis,clickhouse,postgres,mongodb数据源
- 插件市场：可安装操作数据源的各种插件。

## 3. 二开使用说明

```
- node版本 >= v20.14.0
- golang版本 >= v1.23
- IDE推荐：Goland
```

### 3.1 后端工程


```bash

# 克隆项目
git clone https://github.com:1340691923/ElasticView.git

# 安装gowatch
go install github.com/silenceper/gowatch@latest

# 运行
gowatch

# 后台默认端口
0.0.0.0:8090

```

### 3.2 前端工程

```bash
# 进入web文件夹
cd resources\vue

# 安装依赖
pnpm install

# 启动web项目
npm run dev
```
### 3.3 进行打包

```bash

# 安装打包工具

go install github.com/1340691923/ElasticView/cmd/ev_builder@v0.0.22

# 执行打包命令

./ev_builder

# 生成于 resource/dist 文件夹内

```

### 3.4. 技术选型

- 前端：用基于 [Vue](https://vuejs.org) 的 [Element](https://github.com/ElemeFE/element) 构建基础页面。
- 后端：用 [Gin](https://gin-gonic.com/) 快速搭建基础API，[Gin](https://gin-gonic.com/) 是一个go语言编写的Web框架。
- 数据库：使用 [gorm](http://gorm.cn) 实现对数据库的基本操作。
- API文档：使用`Swagger`构建自动化文档。
- 配置文件：使用 [viper](https://github.com/spf13/viper) 实现`yaml`格式的配置文件。
- 日志：使用 [zap](https://github.com/uber-go/zap) 实现日志记录。

## 4.插件相关

### 4.1 官方插件
-  [ev工具箱](https://github.com/1340691923/ev-tools)：用于管理elasticsearch6，7，8版本索引的插件
-  [插件开发模板](https://github.com/1340691923/eve-plugin-vue3-template)：这是用于快速开发插件的模板工程
-
### 4.2 社区插件
- 待完善

### 4.3 发布插件
-  [ev插件开发者后台](http://dev.elastic-view.cn)：用于发布自己的插件提供给ElasticView用户使用


## 5. 联系方式



### QQ交流群：685549060

### 微信公众号：gh_7247127deece

### 微信交流群
| 微信 |
|  :---:  | 
| <img width="150" src="https://raw.githubusercontent.com/1340691923/ElasticView/main/resources/show_img/weixin.jpg"> 

### 知识星球
| 知识星球 |
|  :---:  | 
| <img width="150" src="https://raw.githubusercontent.com/1340691923/ElasticView/refs/heads/master/resources/vue/src/assets/zsxq.jpg"> 


## 6. 捐赠

如果你觉得这个项目对你有帮助，你可以请作者喝饮料 :tropical_drink: [点我](http://www.elastic-view.cn/suporrt.html)

## 7. 商用注意事项

如果您将此项目用于商业用途，请遵守Apache2.0协议并保留作者技术支持声明。
