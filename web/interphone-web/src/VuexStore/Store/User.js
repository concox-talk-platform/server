const state = {
    information: '',
    subordinate: '',
    Language: ''
    
}
const mutations={
    information(state,value){
       state.information=value
    },
    subordinate(state,value){
        state.subordinate=value
     },
     Language(state,value){
         state.Language=value;
     }
 
}
export default {
    state,
    mutations

}