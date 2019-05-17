<template>
    <div   id="group_sider" style="height:100%"  @click="hide">
        <!-- 侧边栏 -->
        <div class="homePage_sideGroup">
            <div class="sideGroup_head" @click="unfold" >
                <span class="sideGroup_headname">Text群组</span>
            </div>
            <div class="sideGroup_body" v-show="show" >
                    <div class="group_total">
                        <i class="el-icon-caret-bottom"  @click="groupbody_show" v-show="down_show"></i> 
                        <i class="el-icon-caret-right"  @click="groupbody_show" v-show="right_show"></i> 
                        <span class="interphonefamily ">&#xe71b;</span>
                        <div class="group_num" >
                            <span class="group_name">Text</span>
                            <!-- <span class="group_on">  ({{online}} </span>
                            <span class="group_totalnum">  /{{totalnum}}) </span> -->
                        </div> 
                    </div>
                    <div class="group_members" v-show='members_show' >
                        <div v-for="(item,index) in group_list" :key="item.group_info.id" class="members"
                        @click="group_detail(index)" 
                        @mouseover="enter_group(index)"
                        @mouseout="leave_group(index)"
                        @contextmenu.prevent='editorial(index)'
                        :class="{active_on:active === index,active_hover:actived === index}"> 
                            <span class="interphonefamily groupicon">&#xe71b;</span>
                            <div class="members_num" >
                                <span class="members_name">{{item.group_info.group_name}}</span>
                                <span>(</span>
                                <span class="members_on"> {{item.group_info.online_num + 1}}</span>
                                <span >  / </span>
                                <span class="members_totalnum">{{item.device_infos.length}}</span>
                                <span>)</span>
                            </div>  
                            <div class="editorial_group" v-show="editorial_show === index">
                                <div class="ungroup" @click.stop="dissolve" >{{$t('group.dissolve')}}</div>
                                <div class="modification" @click.stop="amend">{{$t('group.amend')}}</div>
                                <div class="editor_member" @click.stop="modified_member">{{$t('group.modified_member')}}</div>
                            </div>
                        </div> 
                    </div>
                    <div class="group_foot"  @click="group_add" >
                            <span class="interphonefamily foot_icon">&#xe71b;</span>
                            <span class="foot_add">+</span>
                            <span class="foot_text">{{ $t("group.add_group") }}</span>                
                    </div>
            </div>     
        </div>
    </div>
</template>
<script>
    export default{
     data() {
         return {
             
         }
     },
     methods: {
        hide(){
                this.editorial_show = false.stop;
                this.media_show = false.stop
        }, 
        unfold(){
                    this.show = !this.show;
                },
        groupbody_show(){
            this.members_show = !this.members_show;
            this.down_show = !this.down_show;
            this.right_show = !this.right_show;
        },  
        group_detail(index){
                    this.active = index ;
                    this.select_group_nummber = index;
                    this.group_div_show = true;
        },    
        // 移入移出
        enter_group(index){
            this.actived = index;     
        }, 
        leave_group(index){
                    index = null;
                    this.actived = index;
        },
        // 编辑组成员
        editorial(index){
            this.editorial_show = index;
            this.dele_num = index;
        },
        // 左侧解散组
        dissolve(){
            this.dialogVisible = true;
        },
        // 左侧修改组
        amend(){
            this.amend_show = true;
            let  amend_num = this.dele_num;
            // let amend_group_list = JSON.parse(localStorage.getItem('group_list'));
            // this.Modify_group_form.num = amend_group_list[amend_num].group_info.id
            this.Modify_group_form.num = this.local_group_list[amend_num].group_info.id;
            this.Modify_group_info = this.local_group_list[amend_num].group_info
        },
        // 左侧成员编辑
        modified_member(){
            this.modified_member_show=true;
            let  modified_num = this.dele_num;
            window.console.log(modified_num);
            let modified_member_data =[];
            modified_member_data = this.local_group_list[modified_num].device_infos;
            this.updata_group_info = this.local_group_list[modified_num].group_info;
            window.console.log(this.local_group_list[modified_num]);
            window.console.log(this.local_group_list[modified_num].group_info);
            let selected_member_data=modified_member_data.map(e =>{
                if(e.hasOwnProperty('id')){
                    return e.id
                }
            })
            this.select_Data=selected_member_data;
                window.console.log(selected_member_data)
            window.console.log(this.select_Data)
        },
        // 添加新的组
        group_add(){
            this.group_div = true;
        },
     },


    }


</script>
<style>
    .homePage_sideGroup{
        width:267px;
        margin-top: 10px;
        margin-left: 10px;
        position: relative;
    }
    .sideGroup_head{
        width: 100%;
        height: 40px;
        margin-bottom: 10px;
        color: white;
        position: relative;
        background-color: #206ba2;
        cursor: pointer;
    }
    .sideGroup_headname{
      position: absolute;
      top: 10px;
      left: 85px;
    }
    .sideGroup_body{
        width: 100%;
        background-color: white;
        height: 700px;
        position: relative;
    }
    .group_total{
        height: 33px;
        position: relative;
        background-color: white;
        cursor: pointer;
    }
    .el-icon-caret-bottom,.el-icon-caret-right{
        font-size: 25px
    }
    .interphonefamily {
      font-family: "interphonefamily" !important;
      font-size: 25px;
      font-style: normal;
      -webkit-font-smoothing: antialiased;
      -moz-osx-font-smoothing: grayscale;
    }
    .group_num{
        display: inline-block;
        position: absolute;
        top: 8px;
        left: 58px;
    }
    .group_members{
        cursor: pointer;
        margin-left:52px;
        overflow: auto;
        height: 600px;
    }
    .members{
        height:35px;
        position: relative;
    }
</style>