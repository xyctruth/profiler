export default {
  namespaced: true,
  state: {
    userInfo: null,
  },
  mutations: {
    setUserInfo(state, userInfo) {
      state.userInfo = userInfo
    },

  },
  actions: {
    async LoginIn({commit}, loginInfo) {
      // 登录后返回信息
      const res = {
        data: {
          userName: "我是谁",
        }
      }
      commit('setUserInfo', res.data)
      return res
    },
  },
  getters: {
    userInfo(state) {
      return state.userInfo
    },
  },
}
