const childProcess = require('child_process')
try{
  childProcess.execSync('git rev-parse --abbrev-ref HEAD').toString().replace(/\s+/, '')
}
catch (e) {}

const URL_SETTING_DATA = function (){
  console.log(process.env.npm_config_base_api_url)
  return  {
    namespace: 'profiler-ui',
    reqUrl: process.env.npm_config_base_api_url || "http://localhost:8080",
    publicPath: '/',
  }
}()
module.exports.URL_SETTING_DATA = URL_SETTING_DATA

