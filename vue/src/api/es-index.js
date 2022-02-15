import request from '@/utils/request'

const api = '/api/es_index/'

export function CreateAction(data) {
  return request({
    url: api + 'CreateAction',
    method: 'post',
    data
  })
}

export function GetSettingsAction(data) {
  return request({
    url: api + 'GetSettingsAction',
    method: 'post',
    data
  })
}

export function IndexNamesAction(data) {
  return request({
    url: api + 'IndexNamesAction',
    method: 'post',
    data
  })
}

export function ReindexAction(data) {
  return request({
    url: api + 'ReindexAction',
    method: 'post',
    data
  })
}

export function GetAliasAction(data) {
  return request({
    url: api + 'GetAliasAction',
    method: 'post',
    data
  })
}

export function DeleteAction(data) {
  return request({
    url: api + 'DeleteAction',
    method: 'post',
    data
  })
}

export function OperateAliasAction(data) {
  return request({
    url: api + 'OperateAliasAction',
    method: 'post',
    data
  })
}

export function GetSettingsInfoAction(data) {
  return request({
    url: api + 'GetSettingsInfoAction',
    method: 'post',
    data
  })
}

export function StatsAction(data) {
  return request({
    url: api + 'StatsAction',
    method: 'post',
    data
  })
}

export function CatStatusAction(data) {
  return request({
    url: api + 'CatStatusAction',
    method: 'post',
    data
  })
}

