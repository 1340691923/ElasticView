import request from '@/utils/request'

const api = '/api/es_link/'

export function DeleteAction(data) {
  return request({
    url: api + 'DeleteAction',
    method: 'post',
    data
  })
}

export function ListAction(data) {
  return request({
    url: api + 'ListAction',
    method: 'post',
     data
  })
}

export function OptAction(data) {
  return request({
    url: api + 'OptAction',
    method: 'post',
     data
  })
}

export function UpdateAction(data) {
  return request({
    url: api + 'UpdateAction',
    method: 'post',
    data
  })
}

export function InsertAction(data) {
  return request({
    url: api + `InsertAction`,
    method: 'post',
    data
  })
}

export function InsertEsCfgAction(data) {
  return request({
    url: api + 'InsertEsCfgAction',
    method: 'post',
    data
  })
}

export function UpdateEsCfgAction(data) {
  return request({
    url: api + 'UpdateEsCfgAction',
    method: 'post',
    data
  })
}

export function DeleteEsCfgAction(data) {
  return request({
    url: api + 'DeleteEsCfgAction',
    method: 'post',
    data
  })
}
export function TreeAction() {
  return request({
    url: api + 'TreeAction',
    method: 'post'
  })
}
export function GetEsCfgByRoleId(data) {
  return request({
    url: api + 'GetEsCfgByRoleId',
    method: 'post',
    data
  })
}

export function SetEsCfgByRoleID(data) {
  return request({
    url: api + 'SetEsCfgByRoleID',
    method: 'post',
    data
  })
}

export function GetEsCfgList(data) {
  return request({
    url: api + 'GetEsCfgList',
    method: 'post',
    data
  })
}

export function GetEsCfgOpt(data) {
  return request({
    url: api + 'GetEsCfgOpt',
    method: 'post',
    data
  })
}

export function DeleteEsCfgRelation(data) {
  return request({
    url: api + 'DeleteEsCfgRelation',
    method: 'post',
    data
  })
}
