import request from '@/utils/request'

export function getRoutes() {
  return request({
    url: '/api/gm_user/routes',
    method: 'get'
  })
}

export function getRoles() {
  return request({
    url: '/api/gm_user/roles',
    method: 'get'
  })
}

export function addRole(data) {
  return request({
    url: '/api/gm_user/role/add',
    method: 'post',
    data
  })
}

export function updateRole(data) {
  return request({
    url: `/api/gm_user/role/update`,
    method: 'post',
    data
  })
}

export function deleteRole(id) {
  return request({
    url: `/api/gm_user/role/delete`,
    method: 'get',
    params: id
  })
}
