// i18n-setup.js
import Vue from 'vue'
import VueI18n from 'vue-i18n'
import en from './en'
import zh from './zn'
import axios from 'axios'
import store from '@/store'

Vue.use(VueI18n)

export const i18n = new VueI18n({
  locale: localStorage.getItem('lang') || 'en', // 设置语言环境
  fallbackLocale: localStorage.getItem('lang') || 'en',
  messages: {
    en,
    zh
  } // 设置语言环境信息
})

const loadedLanguages = ['en', 'zh'] // 我们的预装默认语言

export function setI18nLanguage(lang) {
  i18n.locale = lang
  localStorage.setItem('lang', lang)

  const langMap = {
    zh: 'cn',
    en: 'en'
  }
  axios.defaults.headers.common['Accept-Language'] = langMap[lang]
  document.querySelector('html').setAttribute('lang', lang)
  return lang
}

export function i18nText(val) { // 在其他的js文件中引入
  return i18n.t(val)
}
