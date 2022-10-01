// 鉴权 token 相关函数

const TokenKey = 'Admin-Token'
const UserName = 'UserName'

export function getToken() {
  // return Cookies.get(TokenKey)
  return localStorage.getItem(TokenKey)
}

export function setToken(token) {
  return localStorage.setItem(TokenKey, token)
}


export function removeToken() {
  return localStorage.removeItem(TokenKey)
}


export function getName() {
  // return Cookies.get(TokenKey)
  return localStorage.getItem(UserName)
}

export function setName(name) {
  return localStorage.setItem(UserName, name)
}
