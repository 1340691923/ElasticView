import Vue from 'vue'

import Cookies from 'js-cookie'

import 'normalize.css/normalize.css' // a modern alternative to CSS resets

import '@/styles/index.scss' // global css

import VueIntro from 'vue-introjs'
Vue.use(VueIntro)
import 'intro.js/introjs.css'
// 复制粘贴
import VueClipboard from 'vue-clipboard2'
Vue.use(VueClipboard)

Vue.config.productionTip = false

import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'

import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'

Vue.use(Antd)
Vue.config.devtools = true;
import { i18n } from '@/utils/lang'

// excel导出
import JsonExcel from 'vue-json-excel'
Vue.component('downloadExcel', JsonExcel)

import splitPane from 'vue-splitpane'
Vue.component('split-pane', splitPane)

import App from './App'
import store from './store'
import router from './router'

import './permission' // permission control
import './utils/error-log' // error log
import './icons' // error log

// 全局注册组件
import Editor from 'vue2-ace-editor'
Vue.component('editor', Editor)

import * as filters from './filters' // global filters

Vue.use(ElementUI, {
  size: Cookies.get('size') || 'medium' // set element-ui default size
})

// 注册全局过滤器
Object.keys(filters).forEach(key => {
  Vue.filter(key, filters[key])
})

Vue.config.productionTip = false

import { message } from './utils/singleMsg.js'

Vue.prototype.$message = message

new Vue({
  el: '#app',
  router,
  i18n,
  store,
  render: h => h(App)
})

