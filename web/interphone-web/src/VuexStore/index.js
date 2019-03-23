import Vue from 'vue'
import Vuex from 'vuex'
import Header from './Store/Header'
import Account from './Store/Account'
import Area from './Store/Area'
import Control from './Store/Control'
import Equipment from './Store/Equipment'
import Statistics from './Store/Statistics'
import User from './Store/User'
Vue.use(Vuex)
export default new Vuex.Store({
    strict: process.env.NODE_ENV !== 'production',
    modules : {
        Header,
        Account,
        Area,
        Control,
        Equipment,
        Statistics,
        User
    }

})