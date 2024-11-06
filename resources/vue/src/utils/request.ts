import axios, { InternalAxiosRequestConfig } from "axios";
import { useUserStoreHook } from "@/store/modules/user";
import { ResultEnum } from "@/enums/ResultEnum";
import { TOKEN_KEY } from "@/enums/CacheEnum";
import {getToken, removeToken, setToken} from "@/utils/auth";

// 创建 axios 实例
const service = axios.create({
  baseURL: getBaseURL(),
  timeout: 600000,
  headers: { "Content-Type": "application/json;charset=utf-8" },
});

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    config.headers["X-Token"] = getToken()

    config.headers["X-Version"] = window["appVersion"]

    if (localStorage.getItem('lang') != null) {
      config.headers['Current-Language'] = localStorage.getItem('lang')
    }

    return config;
  },
  (error: any) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response) => handleData(response),
  (error: any) => {

    const { response } = error
    if (response === undefined) {
      ElMessage({
        message: '网络异常', // error.message,
        type: 'error',
        duration: 5 * 1000
      })
      return Promise.reject(error)
    } else return handleData(response)
  }
);

const handleData = (response) => {
  // 检查配置的响应类型是否为二进制类型（'blob' 或 'arraybuffer'）, 如果是，直接返回响应对象
  if (response.config.responseType == 'arraybuffer') {
    return response.data
  }

  const res = response.data

  // if the custom code is not 20000, it is judged as an error.
  if (res.code !== 0) {
    if (res.code === 40001) {
      ElMessage.error('请重新登录!')
    }

    if (res.code === 40007){
      ElMessageBox.confirm(res.msg, "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        location.reload();
      });

    }

    if (res.code === 40003 ) {

      let config = response.config
      if (!config._retry){

        let newToken = res.newToken
        setToken(newToken)
        config._retry = true
        config.headers["X-Token"] = getToken()
        return axios(config);
      }

      ElMessageBox.confirm('登录超时，您可以取消以停留在此页，或重新登录', "警告", {
        confirmButtonText: "重新登录",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        removeToken()
        location.reload();
      });
    }

    if (res.code === 40008 ) {
      ElMessage.error(res.msg)
      removeToken()
      location.reload();
    }

    if (res.code === 40002 ) {

      ElMessageBox.confirm('登录验证失败，您可以取消以停留在此页，或重新登录', "警告", {
        confirmButtonText: "重新登录",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        removeToken()
        location.reload();
      });
    }
  }
  return res
}

// 导出 axios 实例
export default service;

function getBaseURL(){
  if(import.meta.env.PROD) {
    return window["appUrl"]
  }
  return import.meta.env.VITE_APP_API_URL
}
