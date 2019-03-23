const state = {
    groupList:[],
    device:[]
}
const mutations = {
   groupList(state,value){
       state.groupList = value
   },
   deviceList(state,value){
       state.device = value
   }
}
const action = {
    getgroupList (context){
        context.commit('groupList')
    }
 
}
const getters = {
    // getdevice:state =>{
    //     return state.device
    //    }
}
export default {
    state,
    mutations,
    action,
    getters

}