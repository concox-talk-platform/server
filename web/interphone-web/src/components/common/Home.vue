<template>
    <div class="wrapper">
        <v-head></v-head>
        <div class="content-box" :class="{'content-collapse':collapse}">
            <div class="content">
                <transition name="move" mode="out-in">
                    <keep-alive :include="tagsList">
                        <router-view></router-view>
                    </keep-alive>
                </transition>
            </div>
        </div>
    </div>
</template>

<script>
    import vHead from './Header.vue';
    
    export default {
        data(){
            return {
                tagsList: [],
                collapse: false
            }
        },
        components:{
            vHead
        },
        // beforeCreate(){
        //     // this.$axios.get('/account/'+localStorage.getItem('loginName'),
        //     this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
        //     { headers: 
        //     {
        //     "Authorization" : sessionStorage.getItem('setSession_id')
        //      }
        //      })
        //     .then((response) =>{
        //         window.console.log(response);
        //      sessionStorage.setItem('id', response.data.account_info.id);
        //      localStorage.setItem('account_info', response.data.account_info);
        //      localStorage.setItem('device_info', response.data.account_info);
        //      this.$store.commit("groupList",response.data.group_list);
             
        //      this.$store.commit("deviceList",response.data.device_list);
        //      this.$store.commit("information",response.data.account_info);
        //         this.$axios.get('/account_class/'+response.data.account_info.id,
        //         { headers: 
        //         {
        //         // "Authorization" : localStorage.getItem('setSession_id')
        //         "Authorization" : sessionStorage.getItem('setSession_id')
        //         }
        //         })
        //         .then((response) =>{
        //         //  localStorage.setItem('id', response.data.account_info.id);    
        //         window.console.log(response);
        //         this.$store.commit("subordinate",response.data.tree_data);
        
        //         })
        //         .catch(function (error) {
        //         window.console.log(error);
        //         });

        //     })
        //     .catch(function (error) {
        //     window.console.log(error);
        //     });
        // },
        // beforeMount(){
        // // this.lang =localStorage.getItem('lang');
        // this.lang =sessionStorage.getItem('lang');
        // this.$i18n.locale = this.lang;

        // // this.accountInfo.name=
        // },
    }
</script>
