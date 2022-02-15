import request from '@/utils/request'

const api = '/api/es/'

export function PingAction(data) {
  return request({
    url: api + `PingAction`,
    method: 'post',
    data
  })
}

export function CatAction(data) {
  return request({
    url: api + `CatAction`,
    method: 'post',
    data
  })
}

export function RunDslAction(data) {
  return request({
    url: api + `RunDslAction`,
    method: 'post',
    data
  })
}

export function SqlToDslAction(data) {
  return request({
    url: api + `SqlToDslAction`,
    method: 'get',
    params: data
  })
}

export function OptimizeAction(data) {
  return request({
    url: api + `OptimizeAction`,
    method: 'post',
    data
  })
}

export function RecoverCanWrite(data) {
  return request({
    url: api + `RecoverCanWrite`,
    method: 'post',
    data
  })
}

