import router from '@/router'

import {
  registerMicroApps,
  addGlobalUncaughtErrorHandler
} from 'qiankun'
import {CallPluginApi} from "@/api/plugin";
import {GetEsConnect, GetEsConnectVer} from "@/utils/es_link";
import {OptAction} from "@/api/es-link";
import { useAppStoreHook, useSettingsStoreHook,useUserStoreHook } from "@/store";
import {ElLoading} from "element-plus";
import {SubscribeToChannel, publish,unsubscribeFromChannel} from "@/utils/centrifuge";

const useSettingsStore = useSettingsStoreHook()
const useAppStore = useAppStoreHook();
const userStore = useUserStoreHook()

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

import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

export function RegisterMicroApps(pluginList){

  for(let i in pluginList){
    let pluginData = pluginList[i]

    let props = {
      CallToChannel:(channel,msg)=>{
        let pluginId = pluginData["name"]
        publish(`${pluginId}$v$${channel}`,msg,(res)=>{

        },(err)=>{
          console.error(`${pluginId}$v$${channel}`,msg,err)
        })
      },
      SubToChannel:(channel,msgCb)=>{
        let pluginId = pluginData["name"]
        SubscribeToChannel(`${pluginId}$v$${channel}`,(res)=>{
          msgCb(res)
        })

      },
      UnSubscribeToChannel:(channel)=>{
        let pluginId = pluginData["name"]
        unsubscribeFromChannel(`${pluginId}$v$${channel}`)
      },
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
      getUserId:()=>{
        return  userStore.user.userId
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

    pluginList[i].container =  "#Appmicro"
    pluginList[i].props =  props
    if(!import.meta.env.PROD){
      if(pluginList[i].entry.indexOf("http") ===-1){
        pluginList[i].entry = import.meta.env.VITE_APP_API_URL+pluginList[i].entry
      }
    }
  }
  console.log(pluginList)

  registerMicroApps(pluginList, {
    // qiankun 生命周期钩子 - 加载前
    beforeLoad: (app) => {
      NProgress.start()
      loadingInstance = ElLoading.service({
        lock: true,
        text: app.name+'插件正在加载中...',
        background: 'rgba(0, 0, 0, 0.7)'
      }); // 开启全局Loading
      console.log('加载插件前，加载进度条',app.name)
      return Promise.resolve()
    },
    beforeMount:(app) =>{
      console.log('挂载插件前，加载进度条', app.name)
      return Promise.resolve()
    },
    // qiankun 生命周期钩子 - 挂载后
    afterMount: (app) => {
      NProgress.done()
      if (loadingInstance) {
        loadingInstance.close(); // 关闭全局Loading
      }
      console.log('挂载插件后，进度条加载完成', app.name)
      return Promise.resolve()
    },
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
