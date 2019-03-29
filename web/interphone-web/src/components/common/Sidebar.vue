<template>
    <div class="navigationbar">
         <div class="ul">

         <!-- <router-link class="li"  v-for="(item, index) in navigationerArr" :key="index" :to="item.address">{{item.name}}</router-link> -->
         <router-link class="li" to="/homePage">{{$t('title.Account_center')}}</router-link>
         <router-link class="li" to="/monitor">{{$t('title.monitor')}}</router-link>
         </div>

        
    </div>
</template>

<script>
    export default {
        data() {
            return {
                // navigationerArr:[
                //     {name:this.$t('title.Account_center'),address:"/homePage"},
                //     // {name:this.$t('title.control_center'),address:'/mapControl'},
                //     // {name:this.$t('title.equipment'),address:'/equipment'},
                //     // {name:this.$t('title.consumer'),address:'/usergroup'},
                //     // {name:this.$t('title.area'),address:'/area'},
                //     {name:this.$t('title.monitor'),address:'/monitor'},]
            }
        },
        methods:{
        },
        computed:{
            onRoutes(){
                return this.$route.path.replace('/','');
            },
            
        },
        created(){
            // this.language()
        },
    }
</script>

<style scoped>
    .navigationbar{
        display: block;
        position: absolute;
        top: 58px;
        font-size: 18px;
        width: 100%;
        height: 38px;
        line-height: 38px;
        background-color: #e4eaf2;
        border-bottom: 4px solid #409eff;
        box-sizing: border-box;
    }
    .li{
    display: inline-block;
    width: 117px;
    text-align: center;
    height: 30px;
    line-height: 30px;
    margin-left: 63px;
    color: black;
    border-radius: 6px 7px 0px 0px;
    font-size: 16px;
    }
    .li:hover{
        background-color: #6b9eff;
        color: white ;
        /* font-weight: 600; */
        
    }

    .router-link-active{
        background-color:#6b9eff;
        color: white  ;
        /* font-weight: 600; */
    }
    /* .sidebar::-webkit-scrollbar{
        width: 0;
    }
    .sidebar-el-menu:not(.el-menu--collapse){
        width: 250px;
    }
    .sidebar > ul {
        height:100%;
    } */

</style>
