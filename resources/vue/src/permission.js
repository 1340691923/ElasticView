import router from './router'
import store from './store'
import { message } from '@/utils/singleMsg.js'
// import { Message } from 'element-ui'
// import NProgress from 'nprogress' // progress bar

import { getToken } from '@/utils/auth' // get token from cookie
import getPageTitle from '@/utils/get-page-title'

// NProgress.configure({ showSpinner: false }) // NProgress Configuration

const whiteList = ['/login', '/auth-redirect'] // no redirect whitelist

// 权限认证
router.beforeEach(async(to, from, next) => {
  // start progress bar
  // NProgress.start()

  // set page title
  document.title = getPageTitle(to.meta.title)

  // determine whether the user has logged in
  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      // if is logged in, redirect to the home page
      next({ path: '/' })
      // NProgress.done()
    } else {
      const hasRoles = store.getters.roles && store.getters.roles.length > 0
      if (hasRoles) {
        next()
      } else {
        try {
          const { roles, list } = await store.dispatch('user/getInfo')

          const accessRoutes = await store.dispatch('permission/generateRoutes', list)

          router.addRoutes(list)

          next({ ...to, replace: true })
        } catch (error) {
          // 清除token 并且重新登录
          await store.dispatch('user/resetToken')
          message.error(error || 'Has Error')
          next(`/login?redirect=${to.path}`)
          // NProgress.done()
        }
      }
    }
  } else {
    // 没有登录
    if (whiteList.indexOf(to.path) !== -1) {
      // in the free login whitelist, go directly
      next()
    } else {
      message.error('请先登录，谢谢')
      // other pages that do not have permission to access are redirected to the login page.
      next(`/login?redirect=${to.path}`)
      // NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  // NProgress.done()
})
