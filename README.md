
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

[简体中文](./README-cn.md) | English

# Project Documentation
[Official Website](http://www.elastic-view.cn) 

[Video Tutorial](https://www.bilibili.com/video/BV12tDDYWEP2/?vd_source=d03eb2249d8310afce3f5b90c6081bb3)

[Communication community](https://txc.qq.com/products/666253)


# Important Tips

1. This project has documents and detailed video tutorials from start-up to development to deployment

2. This project requires you to have a certain foundation in golang and vue3

3. You can complete all operations through our tutorials and documents, so we no longer provide free technical services. If you need services, please [add the author's paid support](https://raw.githubusercontent.com/1340691923/ElasticView/main/resources/show_img/weixin.jpg)



## 1. Basic Introduction

### 1.1 Project Introduction

> ElasticView is a full-stack front-end and back-end separated data source management plugin platform developed based on [vue](https://vuejs.org) and [gin](https://gin-gonic.com), integrating jwt authentication, dynamic routing, dynamic menu, casbin authentication, data source management, plugin market and other functions.


## 2. Main functions

- Permission management: Permission management based on `jwt` and `casbin`.

- User management: System administrators assign user roles and role permissions.

- Role management: Create the main object of permission control, and assign different api permissions and menu permissions to roles.

- Data source management: You can set the data source to be managed, and have integrated elasticsearch (6, 7, 8), mysql, redis, clickhouse, postgres, mongodb data sources

- Plugin market: You can install various plugins for operating data sources.

## 3. Secondary Development

```
- node version >= v20.14.0
- golang version >= v1.23
- IDE：Goland
```

### 3.1 run golang


```bash

# git clone
git clone https://github.com:1340691923/ElasticView.git

# install gowatch
go install github.com/silenceper/gowatch@latest

# run
gowatch

```

### 3.2 run vue

```bash

cd resources\vue

# Install Dependencies
pnpm install

# run
npm run dev
```

### 3.3. Technology selection

- Front-end: Use [Element](https://github.com/ElemeFE/element) based on [Vue](https://vuejs.org) to build basic pages.
- Back-end: Use [Gin](https://gin-gonic.com/) to quickly build basic APIs. [Gin](https://gin-gonic.com/) is a web framework written in go language.
- Database: Use [gorm](http://gorm.cn) to implement basic database operations.
- API documentation: Use `Swagger` to build automated documentation.
- Configuration file: Use [viper](https://github.com/spf13/viper) to implement configuration files in `yaml` format.
- Log: Use [zap](https://github.com/uber-go/zap) to implement logging.

## 4. Plugin related

### 4.1 Official plugin
- [ev tools](https://github.com/1340691923/ev-tools): plugin for managing elasticsearch6, 7, 8 version index

### 4.2 Community plugin
- To be improved

### 4.3 Release plugin
- [Developer Backend](http://dev.elastic-view.cn): used to publish your own plugin for ElasticView users


## 5. Contact information

### QQ communication group: 685549060

### WeChat public account: gh_7247127deece

### WeChat communication group
| WeChat |
| :---: |
| <img width="150" src="https://raw.githubusercontent.com/1340691923/ElasticView/main/resources/show_img/weixin.jpg">



## 6. Donation

If you think this project is helpful to you, you can buy the author a drink :tropical_drink: [click me](http://www.elastic-view.cn/suporrt.html)

## 7. Commercial precautions

If you use this project for commercial purposes, please comply with the Apache2.0 agreement and keep the author's technical support statement.
