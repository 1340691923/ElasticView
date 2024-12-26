import request from "@/utils/request";

export function CallPluginApi(
  pluginId,path,method = 'post',
  data = {},
  transformResponse,
  responseType = '') {

  if (path[0] !== '/'){
    path = `/${path}`
  }

  method = method.toLowerCase()

  let cfg = {
    url: `/api/call_plugin/${pluginId}${path}`,
    method: method,
    responseType:responseType,
  }

  if (method == "post"){
    cfg.data = data
  }else{
    cfg.params = data
  }
  if (transformResponse){
    cfg.transformResponse = transformResponse
  }

  return request(cfg)
}
