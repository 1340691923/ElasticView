import { createApp } from "vue";
import App from "./App.vue";
import setupPlugins from "@/plugins";

// 本地SVG图标
import "virtual:svg-icons-register";

// 样式
import "element-plus/theme-chalk/dark/css-vars.css";
import "@/styles/index.scss";
import "uno.css";
import "animate.css";

//json view
import JsonViewer from "vue3-json-viewer";
import "vue3-json-viewer/dist/index.css";
import {GetI18nCfg} from "@/api/i18n";
import i18n from "@/lang";
import {createI18n} from "vue-i18n";
import zhCnLocale from "@/lang/package/zh-cn";
import enLocale from "@/lang/package/en"; // 引入样式


const app = createApp(App);
app.use(setupPlugins);
app.use(JsonViewer)

const loadMessages = async () => {
  try {

    const res = await GetI18nCfg({});
    if(res.code == 0){
      if(res.data !=null) {

        const i18n = createI18n({
          legacy: false,
          locale: window["lang"],
          messages: res.data,
          globalInjection: true,
        })

        app.use(i18n)
      }
    }

    // 挂载应用
    app.mount('#app');
  } catch (error) {
    console.error('Failed to load messages:', error);
    app.mount('#app'); // 即使加载失败，也需要挂载应用
  }
};


loadMessages()
