import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import {join} from 'path'
const isProductionEnv = process.env.NODE_ENV === "production"
function resolve(dir) {
  return join(__dirname, dir)
}

import {URL_SETTING_DATA} from "./src/config/getVariable.js"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: isProductionEnv?URL_SETTING_DATA.publicPath:'/',
  server: {
    host:"0.0.0.0",
    port: "80"
  },
  resolve: {
    alias: {
      "@": resolve('src')
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/assets/css/mixin.scss";`
      }
    }
  },
  define: {
    URL_SETTING_DATA: JSON.stringify(URL_SETTING_DATA)
  },
})
