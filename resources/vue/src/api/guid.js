import request from '@/utils/request'

const api = '/api/gm_guid'

export function Finish(data) {
  return request({
    url: api + '/Finish',
    method: 'post',
    data
  })
}

export function IsFinish(data) {
  return request({
    url: api + `/IsFinish`,
    method: 'post',
    data
  })
}
