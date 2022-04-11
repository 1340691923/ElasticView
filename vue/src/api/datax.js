import request from '@/utils/request'

let api = "/api/datax/"

export function LinkInfoList(data) {
  return request({
    url: api + 'LinkInfoList',
    method: 'post',
    data
  })
}

export function InsertLink(data) {
  return request({
    url: api + 'InsertLink',
    method: 'post',
    data
  })
}


export function DelLinkById(data) {
  return request({
    url: api + 'DelLinkById',
    method: 'post',
    data
  })
}

export function TestLink(data) {
  return request({
    url: api + 'TestLink',
    method: 'post',
    data
  })
}





