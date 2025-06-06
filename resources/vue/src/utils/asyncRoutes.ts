export const asyncRoutes = [
  {
    "path": "/permission",
    "component": "layout",
    "redirect": "/permission/role",
    "alwaysShow": true,
    "meta": {
      "title": "权限",
      "icon": "system"
    },
    "children": [
      {
        "path": "role",
        "name": "role",
        "component": "views/permission/role",
        "meta": {
          "title": "权限组管理",
          "icon": "role",
          "hidden": false
        },
      },
      {
        "path": "user",
        "name": "user",
        "component": "views/permission/user",
        "meta": {
          "title": "用户管理",
          "icon": "el-icon-user",
          "hidden": false
        },

      },
      {
        "path": "oauth",
        "name": "oauth",
        "component": "views/permission/oauth",
        "meta": {
          "title": "第三方登录",
          "icon": "el-icon-promotion",
          "hidden": false
        },

      },
      {
        "path": "operater_log",
        "name": "operater_log",
        "component": "views/permission/operater_log",
        "meta": {
          "title": "操作日志列表",
          "icon": "dict",
          "hidden": false
        },

      }
    ]
  },
  {
    "path": "/connect-tree",
    "component": "layout",
    "redirect": "/connect-tree/link",
    "alwaysShow": false,
    "meta": {
      "title": "数据源",
      "icon": "el-icon-link"
    },
    "children": [
      {
        "path": "link",
        "name": "link",
        "component": "views/connect-tree/link",
        "meta": {
          "title": "数据源管理",
          "icon": "el-icon-link",
          "hidden": false
        },

      },
      {
        "path": "auth",
        "name": "auth",
        "component": "views/connect-tree/auth",
        "meta": {
          "title": "鉴权管理",
          "icon": "el-icon-user",
          "hidden": false
        },
      }
    ]
  },
  {
    "path": "/plugins",
    "component": "layout",
    "redirect": "/plugins/market",
    "alwaysShow": false,
    "meta": {
      "title": "插件",
      "icon": "el-icon-goods-filled"
    },
    "children": [
      {
        "path": "market",
        "name": "market",
        "component": "views/plugins/market",
        "meta": {
          "title": "插件市场",
          "icon": "el-icon-goods-filled",
          "hidden": false
        },

      },
      {
        "path": "manager",
        "name": "manager",
        "component": "views/plugins/manager",
        "meta": {
          "title": "插件列表",
          "icon": "el-icon-list",
          "hidden": false
        },

      }
    ]
  }
]
