const state = {
  reqCancelMap: {},
  EsConnectID: 0
}

const mutations = {
  SET_EsConnect: (state, EsConnect) => {
    state.EsConnectID = EsConnect
  },
  SET_ReqCancelMap: (state, obj) => {
    state.reqCancelMap[obj['token']] = obj['fn']
  },
  DElETE_ReqCancelMap: (state, token) => {
    delete (state.reqCancelMap[token])
  }
}

const actions = {
  SetEsConnect({ commit }, p) {
    commit('SET_EsConnect', p)
  },
  SET_ReqCancelMap({ commit }, p) {
    commit('SET_ReqCancelMap', p)
  },
  DElETE_ReqCancelMap({ commit }, p) {
    commit('DElETE_ReqCancelMap', p)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
