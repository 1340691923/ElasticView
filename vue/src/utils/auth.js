// 鉴权 token 相关函数
import Cookies from 'js-cookie'

const TokenKey = 'Admin-Token'
const UserName = 'UserName'

export function getToken() {
  // return Cookies.get(TokenKey)
  return sessionStorage.getItem(TokenKey)
}

export function setToken(token) {
  return sessionStorage.setItem(TokenKey, token)
}

export function getName() {
  // return Cookies.get(TokenKey)
  return sessionStorage.getItem(UserName)
}

export function setName(name) {
  return sessionStorage.setItem(UserName, name)
}

export function removeToken() {
  return sessionStorage.removeItem(TokenKey)
}
