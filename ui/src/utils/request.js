import axios from 'axios' // 引入axios
import config from "@/config/config"
import { ElMessage } from 'element-plus'


const service = axios.create({
  baseURL: config.reqUrl,
  timeout: 30000,
  validateStatus: () => true, // 是否自动判断状态码 即400以下走then，400或400以上走catch
})



// http request 拦截器
service.interceptors.request.use(
  config => {
    return config
  },
  error => {
    ElMessage({
      showClose: true,
      message: error,
      type: 'error',
    })
    return error
  },
)

// http response 拦截器
service.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    ElMessage({
      showClose: true,
      message: error,
      type: 'error',
    })
    return error
  },
)

export default service
