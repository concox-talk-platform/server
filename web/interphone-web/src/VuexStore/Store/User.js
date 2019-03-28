const state = {
    information: '',
    subordinate: '',
    
}
const mutations={
    information(state,value){
       state.information=value
    },
    subordinate(state,value){
        state.subordinate=value
     },
 
}
export default {
    state,
    mutations

}