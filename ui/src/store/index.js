import VuexPersistence from 'vuex-persist'
import {createStore} from 'vuex'
// https://webpack.js.org/guides/dependency-management/#requirecontext
const modulesFiles = import.meta.globEager("./modules/*.js")
// you do not need `import app from './modules/app'`
// it will auto require all vuex module from modules file
// const modules = Object.values(modulesFiles).map((item) => {
//   return item.default
// })
const modules = Object.keys(modulesFiles).reduce((modules, modulePath) => {
  // set './app.js' => 'app'
  const moduleName = modulePath.replace(/^\.\/modules\/(.*)\.\w+$/, '$1')
  modules[moduleName] = modulesFiles[modulePath].default
  return modules
}, {})

const vuexLocal = new VuexPersistence({
  key:  'vuex',
  storage: window.localStorage,
  modules: ['user'],
})

export default createStore({
  modules: modules,
  plugins: [vuexLocal.plugin],
})
