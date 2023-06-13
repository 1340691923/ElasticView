import request from '@/utils/request'

export function setIndexCfg(data) {
  return request({
    url: '/api/search/setIndexCfg',
    method: 'post',
    data
  })
}

export function getIndexCfg(data) {
  return request({
    url: '/api/search/getIndexCfg',
    method: 'post',
    data
  })
}

export function SearchLog(data) {
  return request({
    url: '/api/search/SearchLog',
    method: 'post',
    data
  })
}

