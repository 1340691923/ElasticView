// 鉴权 token 相关函数

const TokenKey = 'Admin-Token'

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
