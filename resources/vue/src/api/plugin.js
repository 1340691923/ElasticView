import request from "@/utils/request";

export function CallPluginApi(req) {
  let pluginId = req.pluginAlias
  let path = req.url
  let method = req.method
  let header = req.header
  let data = req.data
  let responseType = req.responseType
  let transformResponse = req.transformResponse
  console.log("req",req)
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
  if (req.hasOwnProperty("transformResponse") ){
    cfg.transformResponse = transformResponse
  }

  if (req.hasOwnProperty("header") && Object.keys(headers).length > 0 ){
    cfg.header = header
  }

  return request(cfg)
}
