const state = {
    setSession_id :'',
    loginName: '',
    id:'',
}
const mutations={
   setSession_id(state,value){
       state.setSession_id=value
   },
   loginName(state,value){
    state.loginName=value
   },
   id(state,value){
    state.id=value
   },

}
export default {
    state,
    mutations

}