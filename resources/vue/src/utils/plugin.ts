import router from '@/router'

import {
  registerMicroApps,
  addGlobalUncaughtErrorHandler
} from 'qiankun'
import {CallPluginApi} from "@/api/plugin";
import {GetEsConnect, GetEsConnectVer} from "@/utils/es_link";
import {OptAction} from "@/api/es-link";
import { useAppStoreHook, useSettingsStoreHook} from "@/store";
import {ElLoading} from "element-plus";

const useSettingsStore = useSettingsStoreHook()
const useAppStore = useAppStoreHook();

const onChangeSettingsStore = (fn:any) => {
  return watch(
    () => useSettingsStore.$state, // 监听状态的函数
    (state) => {
      fn(useSettingsStore);//fn为子应用传过来的函数，即回调函数
    },{
      deep:true
    }
  );
};

const onChangeAppStore = (fn:any) => {
  return watch(
    () => useAppStore.$state, // 监听状态的函数
    (state) => {
      fn(useAppStore);//fn为子应用传过来的函数，即回调函数
    },{
      deep:true
    }
  );
};

let loadingInstance = null; // 保存Loading实例

export function RegisterMicroApps(pluginList){
  let props = {
    CallPluginApi:CallPluginApi,
    GetSelectEsConnID:()=>{
      return GetEsConnect()
    },
    GetSelectEsVersion:()=>{
      return GetEsConnectVer()
    },
    LinkOptAction:async ()=>{
      return OptAction({
        getByLocal:1
      })
    },
    GetI18nMessage: ()=>{
      return window["ev_i18n_message"]
    },
    store:{
      useSettingsStore,
      useAppStore,
      onChangeSettingsStore,
      onChangeAppStore
    },
    router,
  }
  for(let i in pluginList){

    pluginList[i].container =  "#Appmicro"
    pluginList[i].props =  props
  }
  console.log(pluginList)

  registerMicroApps(pluginList, {
    // qiankun 生命周期钩子 - 加载前
    beforeLoad: (app) => {
      loadingInstance = ElLoading.service({
        lock: true,
        text: '插件正在加载中...',
        background: 'rgba(0, 0, 0, 0.7)'
      }); // 开启全局Loading
      console.log('加载插件前，加载进度条', app.name)
      return Promise.resolve()
    },
    beforeMount:(app) =>{
      console.log('挂载插件前，加载进度条', app.name)
      return Promise.resolve()
    },
    // qiankun 生命周期钩子 - 挂载后
    afterMount: (app) => {
      if (loadingInstance) {
        loadingInstance.close(); // 关闭全局Loading
      }
      console.log('挂载插件后，进度条加载完成', app.name)
      return Promise.resolve()
    }
  })
}

/**
 * 添加全局的未捕获异常处理器
 */
addGlobalUncaughtErrorHandler((event) => {
  const { message: msg } = event
  console.log("event",event)
  // 加载失败时提示
  if (msg && msg.includes('died in status LOADING_SOURCE_CODE')) {
    console.log('请检查插件是否可运行，插件加载未成功',msg)
  }
})
