import request from '@/utils/request'

const api = '/api/es_crud/'

export function GetList(data) {
  return request({
    url: api + 'GetList',
    method: 'post',
    transformResponse : [
      data => {
        return jsonlint.parse(data)
      }
    ],
    data,
  })
}

export function GetDSL(data) {
  return request({
    url: api + 'GetDSL',
    method: 'post',
    data
  })
}

export function Download(data) {
  return request({
    responseType: 'arraybuffer', // 必填
    url: api + 'Download',
    method: 'post',
    data
  })
}

