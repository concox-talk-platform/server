<template>
    <div id="app">
        <router-view isRouterAlive></router-view>
    </div>
</template>
<script>
export default {
    provide(){
      return{
         reload: this.reload
      }  
    },
    data(){
        return{
            isRouterAlive: true
        }
    },
    methods:{
        reload(){
            this.isRouterAlive=false
            this.$nextTick(function(){
                this.isRouterAlive=true
            })
            },
        
        },
        mounted(){
            // // 关闭浏览器窗口的时候清空浏览器缓存在localStorage的数据
            // window.onbeforeunload = function () {
            //     var storage = window.localStorage;
            //     storage.clear()
            // }
        }

    
    }

</script>
<style>
    @import "./assets/css/main.css";
    @import "./assets/css/color-dark.css";     /*深色主题*/
    /*@import "./assets/css/theme-green/color-green.css";   浅绿色主题*/
</style>