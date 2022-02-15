// 自动生成 路由映射文件 脚本
var asyncRoutes = [

]

function filterAsyncRoutes(routes) {
  routes.forEach(route => {
    const tmp = { ...route }
    if (tmp.children) {
      tmp.children = filterAsyncRoutes(tmp.children)
    }
    if (tmp.component != 'layout') {
      console.log("'" + tmp.component + "':() => import('@/" + tmp.component + "'),")
    }
  })
}

filterAsyncRoutes(asyncRoutes)

