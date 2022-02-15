import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/api/gm_user/login',
    method: 'post',
     data
  })
}

export function getInfo() {
  return request({
    url: '/api/gm_user/info',
    method: 'get',
    params: {}
  })
}

export function logout() {
  return request({
    url: '/api/gm_user/logout',
    method: 'post'
  })
}

export function userList() {
  return request({
    url: '/api/gm_user/userlist',
    method: 'post'
  })
}

export function roleOption() {
  return request({
    url: '/api/gm_user/roleOption',
    method: 'post'
  })
}

export function getUserById(data) {
  return request({
    url: '/api/gm_user/getUserById',
    method: 'get',
    params: data
  })
}

export function UpdateUser(data) {
  return request({
    url: '/api/gm_user/UpdateUser',
    method: 'get',
    params: data
  })
}

export function InsertUser(data) {
  return request({
    url: '/api/gm_user/InsertUser',
    method: 'get',
    params: data
  })
}

export function DelUser(id) {
  return request({
    url: '/api/gm_user/DelUser',
    method: 'get',
    params: id
  })
}

export function GetAddUserCount() {
  return request({
    url: '/api/GetAddUserCount',
    method: 'post'
  })
}
