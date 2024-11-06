import request from "@/utils/request";

export function getRoles() {
  return request({
    url: '/api/gm_user/roles',
    method: 'post'
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

export function deleteRole(data) {

  return request({
    url: `/api/gm_user/role/delete`,
    method: 'post',
    data
  })
}
