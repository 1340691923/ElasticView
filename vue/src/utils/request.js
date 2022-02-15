// axios 拦截器
import store from '@/store'

import { getToken } from '@/utils/auth'

import { message } from '@/utils/singleMsg.js'

const CancelMsg = '用户已经取消请求'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 60 * 60 * 60 * 60 // request timeout
})

// 声明一个数组用于存储每个请求的取消函数和axios标识

// request interceptor
service.interceptors.request.use(
  config => {
    if (config.data) {
      if (config.data.hasOwnProperty('cancelToken')) {
        config.cancelToken = new axios.CancelToken((c) => {
          store.dispatch('baseData/SET_ReqCancelMap', { token: config.data['cancelToken'], fn: c })
        })
      }
    }

    // do something before request is sent
    if (store.getters.token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['X-Token'] = getToken()
    }
    return config
  },
  error => {
    console.log(' request err', err)
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
   */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    console.log('response', response)
    if (response.config.responseType == 'arraybuffer') {
      return response.data
    }

    const res = response.data

    // if the custom code is not 20000, it is judged as an error.
    if (res.code !== 0) {
      if (res.code === 40001) {
        message.error('请重新登录!')
      }

      if (res.code === 40002 || res.code === 40003) {
        // to re-login
        ELEMENT.MessageBox.confirm('未登录或登录超时，您可以取消以停留在此页，或重新登录', '确认注销', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          store.dispatch('user/resetToken').then(() => {
            location.reload()
          })
        })
      }
      return res
      // return Promise.reject(new Error(res.msg || '未知错误'))
    } else {
      return res
    }
  },
  error => {
    if (error.message == CancelMsg) {
      message({
        message: '已经手动取消请求', // error.message,
        type: 'success',
        duration: 5 * 1000
      })
      return Promise.reject(error)
    }
    console.error('err', error) // for debug
    message({
      message: '网络异常', // error.message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

service.defaults.transformResponse = [
  data=>{
    return jsonlint.parse(data)
  }
]

export default service
