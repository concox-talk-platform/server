import Vue from 'vue'
import App from './App.vue'
import store from './VuexStore'
import router from './router'
import axios from 'axios';
import ElementUI from 'element-ui';
import VueI18n from 'vue-i18n'
import './assets/icon/iconfont.css';
import 'element-ui/lib/theme-chalk/index.css'; // 默认主题
import "babel-polyfill";
// import { push_uniq } from '../node_modules/terser';
// import { is } from 'css-select';
// import { isNull } from 'util';


Vue.config.productionTip = false
Vue.use(ElementUI, {
    size: 'small'
});
Vue.use(VueI18n)
const i18n =new VueI18n({
    locale:'en-US',
    messages: {
      'zh-CN': require('./components/common/lang/zh'),
      'en-US': require('./components/common/lang/en')
    },
    silentTranslationWarn: true,
  })
axios.defaults.baseURL = 'http://172.16.0.74:8080'
// axios.defaults.baseURL = 'http://113.105.153.240:8080'
Vue.prototype.$axios = axios;

//使用钩子函数对路由进行权限跳转
router.beforeEach((to, from, next) => {
    // let role=localStorage.getItem('setSession_id');
    let role=sessionStorage.getItem('setSession_id');
    if(role===null&&to.path=='/login'){
    next()
    }else if(role===null){
    router.push('/login')
    next()
    }
    else{
    if(to.path=='/login'){
    router.push('/')
    next()
    }else{
    next()
    }
    }
    })

new Vue({
    
    i18n,
    router,
    store,
    render: h => h(App),
}).$mount('#app')