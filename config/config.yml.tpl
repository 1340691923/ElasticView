log:
  storageDays: 4          # 日志保留天数
  logDir: "logs"          # 日志保留文件夹
port: 8090                # 启动端口
pluginRpcPort: 8091            #插件内网访问端口
dbType: "sqlite3"       # 数据保留类型 分为 sqlite3 和 mysql
enableLogEs: false    #是否记录es请求记录
enableLogEsRes: false #是否记录es请求记录中返回的响应体
sqlite: # dbType为sqlite3时填 dbPath为数据保存文件地址
  dbName: "es_view.db"
mysql: # dbType为mysql时填
  username: "root"
  pwd: ""
  ip: "localhost"
  port: "3306"
  dbName: "test"
  maxOpenConns: 10
  maxIdleConns: 10
appSecret: "1340691923@qq.com" # jwt 加密密钥
esPwdSecret: "concat_mail!!->1340691923@qq.com" # es密码加密密钥 加密方式为 AES
version: "0.0.3"  # EV 版本号
deBug: false      # 是否为测试模式 如果为 false则打开默认浏览器直接访问地址
checkForevUpdates: true #是否自动检测ev更新
checkForPluginUpdates: true  #是否自动检测ev插件更新
evKey: #evKey 需要到插件者后台注册获取
storeFileDir: #临时文件存放目录 例如下载的excel
plugin:
  loadPath: plugins   #插件存放目录
  storePath: plugins_store #插件临时文件存放目录
watermarkContent: ElasticView #水印
translation:
  lang: zh-cn # zh-cn or en
  cfgDir: config/i18n  #i18n文件存放目录
Ai:
  bigModeKey: "" #阿里百炼大模型appkey
