/**
 * 开发环境
 * ===================
 * 当你的项目页面越来越多之后，在开发环境之中使用 lazy-loading 会变得不太合适，每次更改代码触发热更新都会变得非常的慢。
 * 所以建议只在生产环境之中使用路由懒加载功能。
 * &&这里注意一下该写法只支持 vue-loader at least v13.0.0+
 */
module.exports = (file) => {
  if (file == 'layout') {
    return require('@/layout/' + file + '.vue').default
  }
  return require('@/' + file + '.vue').default
}
