import request from '@/utils/request'

export function getList(data) {
  return request({
    url: '/api/operater_log/ListAction',
    method: 'post',
    data
  })
}
