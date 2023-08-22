/**
 * 生产环境
 * ===========
 * 当打包构建应用时，Javascript 包会变得非常大，影响页面加载。
 * 如果我们能把不同路由对应的组件分割成不同的代码块，然后当路由被访问的时候才加载对应组件，这样就更加高效了。
 *
 */
/*
module.exports = (file)=>{
  if (file == 'layout'){
    return import('@/layout/' + file + '.vue')
  }
  return import('@/' + file + '.vue')
}
*/
module.exports = (file) => {
  if (file == 'layout') {
    return require('@/layout/' + file + '.vue').default
  }
  return require('@/' + file + '.vue').default
}

