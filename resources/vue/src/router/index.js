import Vue from 'vue'
import Router from 'vue-router'
const _import = require('@/utils/import.' + process.env.NODE_ENV)

Vue.use(Router)

// 基础路由
export const constantRoutes = [
  {
    path: '/redirect',
    component: _import('layout'),
    hidden: true,
    children: [
      {
        path: '/redirect/:path*',
        component: _import('views/redirect/index')
      }
    ]
  },
  {
    path: '/login',
    component: _import('views/login/index'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: _import('views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/404',
    component: _import('views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: _import('views/error-page/401'),
    hidden: true
  },
  {
    path: '/',
    component: _import('layout/layout'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: _import('views/dashboard/index'),
        name: 'Dashboard',
        meta: { title: '首页', icon: 'el-icon-s-home', affix: true }
      }
    ]
  }
  /* { path: '*', redirect: '/404', hidden: true }*/
]

export const asyncRoutes = [
  {
    path: '/permission',
    component: 'layout',
    redirect: '/permission/role',
    alwaysShow: true,
    meta: {
      title: '权限',
      icon: 'el-icon-user-solid'
    },
    children: [
      {
        path: 'role',
        component: 'views/permission/role',
        name: 'role',
        meta: {
          title: '角色管理',
          icon: 'el-icon-s-check'
        }
      },
      {
        path: 'user',
        component: 'views/permission/user',
        name: 'user',
        meta: {
          title: '用户管理',
          icon: 'el-icon-user'
        }
      },
      {
        path: 'operater_log',
        component: 'views/permission/operater_log',
        name: 'operater_log',
        meta: {
          title: '操作日志列表',
          icon: 'el-icon-s-order'
        }
      }
    ]
  },
  {
    path: '/connect-tree',
    component: 'layout',
    redirect: '/connect-tree/index',
    alwaysShow: false,
    meta: {
      title: '连接树管理',
      icon: 'el-icon-link'
    },
    children: [
      {
        path: 'index',
        component: 'views/connect-tree/index',
        name: 'index',
        meta: {
          title: '连接树管理',
          icon: 'el-icon-link'
        }
      }
    ]
  },
  {
    path: '/cat',
    component: 'layout',
    redirect: '/cat/index',
    alwaysShow: false,
    meta: {
      title: 'ES状态',
      icon: 'el-icon-pie-chart'
    },
    children: [
      {
        path: 'index',
        component: 'views/cat/index',
        name: 'index',
        meta: {
          title: 'ES状态',
          icon: 'el-icon-pie-chart'
        }
      }
    ]
  },

  {
    path: '/rest',
    component: 'layout',
    redirect: '/rest/index',
    alwaysShow: false,
    meta: {
      title: '开发工具',
      icon: 'el-icon-edit'
    },
    children: [
      {
        path: 'index',
        component: 'views/rest/index',
        name: 'index',
        meta: {
          title: '开发工具',
          icon: 'el-icon-search'
        }
      }
    ]
  },
  {
    path: '/indices',
    component: 'layout',
    redirect: '/indices/index',
    alwaysShow: true,
    meta: {
      title: '索引管理',
      icon: 'el-icon-coin'
    },
    children: [
      {
        path: 'index',
        component: 'views/indices/index',
        name: 'index',
        meta: {
          title: '索引管理',
          icon: 'el-icon-coin'
        }
      },
      {
        path: 'reindex',
        component: 'views/indices/reindex',
        name: 'reindex',
        meta: {
          title: '重建索引',
          icon: 'el-icon-document-copy'
        }
      }
    ]
  },
  {
    path: '/task',
    component: 'layout',
    redirect: '/task/index',
    alwaysShow: false,
    meta: {
      title: '任务',
      icon: 'el-icon-notebook-2'
    },
    children: [
      {
        path: 'index',
        component: 'views/task/index',
        name: 'index',
        meta: {
          title: '任务',
          icon: 'el-icon-notebook-2'
        }
      }
    ]
  },
  {
    path: '/back-up',
    component: 'layout',
    redirect: '/back-up/index',
    alwaysShow: true,
    meta: {
      title: '备份',
      icon: 'el-icon-copy-document'
    },
    children: [
      {
        path: 'index',
        component: 'views/back-up/index',
        name: 'index',
        meta: {
          title: '快照存储库',
          icon: 'el-icon-first-aid-kit'
        }
      },
      {
        path: 'snapshot',
        component: 'views/back-up/snapshot',
        name: 'index',
        meta: {
          title: '快照管理',
          icon: 'el-icon-shopping-bag-2'
        }
      }
    ]
  },
  {
    path: '/navicat',
    component: 'layout',
    redirect: '/navicat/index',
    alwaysShow: false,
    meta: {
      title: 'Navicat',
      icon: 'el-icon-copy-document'
    },
    children: [
      {
        path: 'index',
        component: 'views/navicat/index',
        name: 'index',
        meta: {
          title: 'Navicat',
          icon: 'el-icon-first-aid-kit'
        }
      }
    ]
  },
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
