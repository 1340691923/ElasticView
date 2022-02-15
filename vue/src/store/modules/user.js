import { getInfo, login, logout } from '@/api/user'
import { getToken, removeToken, setToken } from '@/utils/auth'
import router, { resetRouter } from '@/router'
import { RoutesComponentmaps } from '@/utils/router'
import { message } from '@/utils/singleMsg.js'

const state = {
  token: getToken(),
  name: '',
  avatar: '',
  introduction: '',
  roles: [],
  channos: []
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_INTRODUCTION: (state, introduction) => {
    state.introduction = introduction
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },

  SET_ROLEID: (state, roleid) => {
    state.roleid = roleid
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password }).then(response => {
        const { data, code, msg } = response

        if (code != 0) {
          message({
            message: msg || '未知错误',
            type: 'error',
            duration: 5 * 1000
          })
          reject(msg)
        }

        commit('SET_TOKEN', data.token)
        setToken(data.token)

        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo().then(response => {
        var data = response.data
        // const { data } = response
        if (!data) {
          reject(data.msg)
        }

        const { roles, name, avatar, introduction, list, channos } = data

        var routerList = JSON.parse(list)

        convertTree(routerList)

        routerList.push({ path: '*', redirect: '/404', hidden: true })
        // roles must be a non-empty array
        if (!roles || roles.length <= 0) {
          reject('getInfo: roles must be a non-null array!')
        }
        data.list = routerList
        commit('SET_ROLES', data.roles)
        commit('SET_NAME', data.name)
        commit('SET_AVATAR', data.avatar)
        commit('SET_INTRODUCTION', data.introduction)

        commit('SET_ROLEID', data.roleId)
        resolve(data)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state, dispatch }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        removeToken()
        resetRouter()
        dispatch('tagsView/delAllViews', null, { root: true })

        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      removeToken()
      resolve()
    })
  },

  // dynamically modify permissions
  changeRoles({ commit, dispatch }, role) {
    return new Promise(async resolve => {
      const token = role + '-token'

      commit('SET_TOKEN', token)
      setToken(token)

      const { roles } = await dispatch('getInfo')

      resetRouter()

      // generate accessible routes map based on roles
      const accessRoutes = await dispatch('permission/generateRoutes', roles, { root: true })

      // dynamically add accessible routes
      router.addRoutes(accessRoutes)

      // reset visited views and cached views
      dispatch('tagsView/delAllViews', null, { root: true })

      resolve()
    })
  }
}

function convertTree(routers) {
  const menuRouters = []

  routers.forEach(r => {
    r.component = RoutesComponentmaps[r.component]
    if (r.children) convertTree(r.children)
  })
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
