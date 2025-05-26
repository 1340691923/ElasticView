import request from '@/utils/request'

var api = '/api/plugins/'

export function PluginMarket(data) {
  return request({
    url: api + 'PluginMarket',
    method: 'post',
    data
  })
}

export function GetPluginInfo(data) {
  return request({
    url: api + 'GetPluginInfo',
    method: 'post',
    data
  })
}

export function InstallPlugin(data) {
  return request({
    url: api + 'InstallPlugin',
    method: 'post',
    data
  })
}

export function UploadPlugin(data) {
  return request({
    url: api + 'UploadPlugin',
    method: 'post',
    data
  })
}

export function StarPlugin(data) {
  return request({
    url: api + 'StarPlugin',
    method: 'post',
    data
  })
}

export function ImportEvKey(data) {
  return request({
    url: api + 'ImportEvKey',
    method: 'post',
    data
  })
}

export function UnInstallPlugin(data) {
  return request({
    url: api + 'UnInstallPlugin',
    method: 'post',
    data
  })
}

export function LikeComment(data) {
  return request({
    url: api + 'LikeComment',
    method: 'post', 
    data
  })
}

 
export function ListComments(data) {
  return request({
    url: api + 'ListComments',
    method: 'post', 
    data
  })
}

 
export function AddComment(data) {
  return request({
    url: api + 'AddComment',
    method: 'post', 
    data
  })
}

export function GetWxArticleList(data) {
  return request({
    url: api + 'GetWxArticleList',
    method: 'post',
    data
  })
}

export function GetLocalPluginList(data) {
  return request({
    url: api + 'GetLocalPluginList',
    method: 'post',
    data
  })
}
