<template>
<div style="height:100%"  @click="hide">
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
                                <span class="members_on">  ({{online}} </span>
                                <span class="members_totalnum">  /{{totalnum}}) </span>
                            </div>  
                            <div class="editorial_group" v-show="editorial_show === index">
                                <div class="ungroup" @click.stop="dissolve">{{$t('group.dissolve')}}</div>
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
    <!-- 创建组 -->
    <el-dialog :title="$t('group.add_group')" :visible.sync="group_div">
        <div class="dialog_div">
            <div class="group_top">
                <el-form ref="addgroup_form" :model="addgroup_form" :rules="group_rules" label-width="113px">
                    <el-form-item :label="$t('group.add_name')" prop="name">
                        <el-input v-model="addgroup_form.name"></el-input>
                    </el-form-item>
                 </el-form>   
            </div>
            <div class="group_addmember">
                <div class="group_membertitle">
                    {{$t('group.members')}}
                </div>
                <div class="group_memberdiv">
                       <ul v-if="group_memberdiv_show">
                            <li   v-for="(item,index) in confirm_device_List" :key="item.id" v-on:dblclick="dele_device(index)" :title="$t('group.dbremove')">
                                <div class='group_pic' >
                                   <img src="../../assets/img/inter.png" alt="">
                                </div>
                                <div class="member_name">
                                    {{item.user_name.String}}
                                </div> 
                            </li>                                          
                                                                                     
                       </ul>
                </div>
                <el-button  class="addbotton" @click="add_select_div">{{$t('button_message.add')}}</el-button>
            </div>
        </div>
        <el-dialog
                width="31%" :title="$t('group.member_title')" :visible.sync="members_div" append-to-body>
            <!-- 左右列表移动 -->
                <el-transfer
                    style="text-align: left; display: inline-block"
                    v-model="yesData"
                    :titles="[Source, Target]"
                    :button-texts="[toleft, toright]"
                    :format="{
                        noChecked: '${total}',
                        hasChecked: '${checked}/${total}'
                    }"
                    @change="handleChange"
                    :props="{key: 'id',label:'name'}"
                    :data="noData">
                    <!-- <el-button class="transfer-footer" slot="left-footer" size="small">操作</el-button>
                    <el-button class="transfer-footer" slot="right-footer" size="small">操作</el-button> -->
                </el-transfer>
            <div class="select_add" @click="selected_add">{{$t("button_message.confirm")}}</div>
        </el-dialog>
        <div slot="footer" class="dialog-footer">
        <el-button @click="group_add_cancle">{{ $t("button_message.cancel") }}</el-button>
        <el-button type="primary" @click="group_add_submit('addgroup_form')" >{{ $t("button_message.ensure") }}</el-button>
        </div>
    </el-dialog>
    <!-- 组内详细设备 -->
    <div class="device_detail" v-if="group_div_show">
            <div class="device_detail_tittle">
                <span class="device_detail_name">{{device_group_name}}</span>
                <span class="device_detail_close" @click="close_device_detail">x</span>
                <!-- Text -->
            </div>
            <div class="device_detail_num">
                <ul>
                        <li  @contextmenu.prevent='control(index)'  v-for="(item,index) in select_group_num" :key="item.id">
                            <div class='device_pic'>
                                <img src="../../assets/img/inter.png" alt="">
                            </div>
                            <div class="device_name">
                                {{item.user_name.String}}
                            </div> 
                            <div class="control_menum" v-show="control_show">
                                <div class="control_voice">语音通话</div>
                                <div class="control_vedio">视频通话</div>
                                <div class="control_look">视频查看</div>
                                <div class="control_text">即时消息</div>
                            </div>
                        </li>                             
                </ul>
            </div>
    </div>
     <!-- 删除组弹出框 -->
        <el-dialog :title="$t('control.hint')" :visible.sync="dialogVisible" width="30%">
            <span>{{$t('control.delete_this')}}</span>
            <span slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">{{$t('button_message.cancel')}}</el-button>
                <el-button type="primary" @click="dele_submit">{{$t('button_message.confirm')}}</el-button>
            </span>
        </el-dialog>

     <!-- 修改组弹出框 -->


</div>
</template>

<script>
export default {
    data() {   
        return {
         yesData: [],
         show:false,
         online:6,
         totalnum:20,
         active:null,
         actived:null,
         active_on:'',
         active_hover:'',
         members_show:true,
         down_show:true,
         right_show:false,
        //  group_list:[],
         group_div:false,
         members_div:false, 
         control_show:false,
         editorial_show:false,
         dialogVisible: false,
         select_group_nummber:'a' ,
         group_div_show: false,
         dele_num:null,
         group_memberdiv_show: false,
         confirm_device:[],
         confirm_device_List:[],
         titles_left:this.$t('group.all_member'), 
         titles_right:this.$t('group.select_member'), 
         button_left: this.$t('group.remove'),
         button_right: this.$t('group.add'),
         Source: this.$t('group.all_member'), 
         Target: this.$t('group.select_member'),
         toleft: this.$t('group.remove'),
         toright: this.$t('group.add'),
         addgroup_form: {
         name:'',
         },

         group_rules : {
            name: [
            { required: true, message: this.$t('group.message'), trigger: 'blur' },
          ],
         },

        }
    },
    methods: {
        hide(){
          this.editorial_show=false.stop
        },
        // 鼠标右击
        control(index){
            alert(index)
        },
        unfold(){
            this.show=!this.show;

        },
        groupbody_show(){
            this.members_show=!this.members_show;
            this.down_show=!this.down_show;
            this.right_show=!this.right_show;
        },
        // 点击变色事件
        group_detail(index){
             this.active = index ;
             this.select_group_nummber=index;
             this.group_div_show = true;
           
        },
        close_device_detail(){
             this.group_div_show = false;
        },
        // 移入移出
        enter_group(index){
             this.actived = index     
        },
        leave_group(index){
            index=null;
            this.actived = index;
        },
        // 添加新的组
        group_add(){
         this.group_div = true;
        
        },
        add_select_div(){
         this.members_div=true;

        },
        // 列表移动
      handleChange(value, direction, movedKeys,) {
        window.console.log(value, direction, movedKeys,);
      },
    //   选择确认
       selected_add(){
        this.confirm_device=this.yesData;
        this.confirm_device_List=[];
        for( var i =0;i<this.$store.state.Equipment.device.length;i++){
            for(var j=0;j<this.confirm_device.length;j++){
                if(this.$store.state.Equipment.device[i].id ==this.confirm_device[j]){
                  this.confirm_device_List.push(this.$store.state.Equipment.device[i])  
                }
            }
        }
        this.members_div = false;
        this.group_memberdiv_show = true;

       },
       group_add_cancle(){
            this.group_div = false;
            this.group_memberdiv_show = false;
            // 清空表单
            this.$refs['addgroup_form'].clearValidate();
            this.$refs['addgroup_form'].resetFields();
       },
    //    添加组提交
       group_add_submit(addgroup_form){
           let submit_form={};
       
         this.$refs[addgroup_form].validate((valid) => {
                  if (valid) {
                   
                    submit_form.group_name=this.addgroup_form.name;
                    submit_form.group_device=this.confirm_device_List;
                    if(this.confirm_device_List.length !== 0){
                      window.console.log(submit_form)
                    }else{
                        this.$message({
                        message: this.$t('prompt_message.device_num'),
                        type: 'warning'
                        });
                    }
                           
                  }else{
                      return false;
                  }

             })
       },
    //    双击删除
        dele_device(qua){
            // window.console.log(this.confirm_device_List)
            this.confirm_device_List.splice(qua,1)
            // this.yesData=this.confirm_device_List
            this.yesData.splice(qua,1)
        },
        // 编辑组成员
        editorial(index){
            // alert(index);
            this.editorial_show=index
            this.dele_num=index
        },
        // 解散组
        dissolve(){
            this.dialogVisible = true
             window.console.log(this.dele_num);
             let dele_mumber=this.dele_num
             window.console.log(this.$store.state.Equipment.groupList);
             let dele_groupList= this.$store.state.Equipment.groupList.slice(0,this.$store.state.Equipment.groupList.length);
             window.console.log(dele_mumber);
             window.console.log(dele_groupList);
             window.console.log(11111);
             dele_groupList.splice(dele_mumber,1);
             window.console.log(dele_groupList);
             this.$store.commit("groupList",dele_groupList);

        },
        dele_submit(){

         this.dialogVisible = false
        },
        // 修改组
        amend(){
          window.console.log(111)
        },
        // 成员编辑
        modified_member(){
            window.console.log(111)
        }
    },
    computed:{
   
            group_list(){
                return this.$store.state.Equipment.groupList;
            },
            device_member(){
                 let device_id=this.$store.state.Equipment.device.map(e =>{
                     if(e.hasOwnProperty('id')){
                         return e.id
                     }
                 })
                return device_id
            },
            noData(){
                let  transfer_name=this.$store.state.Equipment.device;
                let transfer_newData = [];
                transfer_name.forEach((obj) => {
                    var  transfer_obj ={};
                    transfer_obj.id = obj.id;
                    transfer_obj.name =obj.user_name.String;
                    transfer_newData.push(transfer_obj)
                });
               
                return  transfer_newData
            },
            // 选中组的设备信息

            select_group_num(){
                let  group_select_device = this.$store.state.Equipment.groupList;
                let  group_selected_num = this.select_group_nummber
                if( group_selected_num !== 'a'){
                    group_select_device = group_select_device[group_selected_num].device_infos
                }
                return   group_select_device
            },
            // 组名
            device_group_name(){
                let group_selected_member =this.$store.state.Equipment.groupList;
                let  group_selected_name = this.select_group_nummber
                if( group_selected_name !== 'a'){
                    group_selected_member = group_selected_member[group_selected_name].group_info.group_name
                }
                return   group_selected_member
            }

            
    },

    mounted(){
        // this.group_list=this.$store.state.Equipment.groupList;

        // this.group_list=this.$store.state.Equipment.groupList;
    }
}
</script>

<style scoped>
.active_on{
    background-color: #206ba2
}
.active_hover{
    background-color: #60A9E1
}
.homePage_sideGroup{
    width:267px;
    margin-top: 10px;
    margin-left: 10px;
    /* margin-bottom: 20px; */
    /* background-color: aqua; */
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
.group_total{
    height: 33px;
    position: relative;
    background-color: white;
    cursor: pointer;
}
.interphonefamily {
  font-family: "interphonefamily" !important;
  font-size: 25px;
  font-style: normal;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
.el-icon-caret-bottom,.el-icon-caret-right{
    font-size: 25px
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
}
.members{
    height:35px;
    position: relative;
}
.groupicon{
    font-size: 26px;
    display: inline-block;
    position: absolute;
    top: 3px;
    left: 4px;
}
.members_num{
    position: absolute;
    top: 8px;
    left: 36px;
    font-size: 14px;

}
.device_detail_close{
    float: right;
    padding-right: 8px;
    cursor: pointer;
}
/* .group_members:hover{
    background-color: 
} */
/* .sideGroup_body{
    width: 100%;
    background-color: brown;
    animation:groupshow .3s linear;
    animation: groupshow .3s linear;
    -moz-animation: groupshow .3s linear;
    -webkit-animation: groupshow .3s linear;
    -o-animation: groupshow .3s linear;
} */
.sideGroup_body{
    width: 100%;
    background-color: white;
    height: 750px;
    position: relative;
}
.group_foot{
    height: 40px;
    width: 100%;
    position: relative;
    position: absolute;
    bottom: 0px;
    background-color: #206ba2;
    cursor: pointer;
}
.group_foot:hover{
  background-color: #60A9E1  
}
.foot_icon{
    font-size: 38px;
    position: absolute;
    left: 12px;
    top: 1px;
}
.foot_add{
    position: absolute;
    top: 8px;
    left: 41px;
}
.foot_text{
    position: absolute;
    top: 10px;
    left: 54px;
}
.dialog_div{
    width: 100%;
    height: 400px;
}
.group_top{
   width: 100%;
   height: 50px;
   
}
.group_addmember{
    /* background-color: burlywood; */
    height: 330px;
    width: 100%;
    border: 1px solid #ccc;
}
.group_membertitle{
    height: 30px;
    text-align: center;
}
.group_memberdiv{
    width: 100%;
    height: 300px;
    /* background-color: aquamarine; */
    overflow: auto ;
}
.group_memberdiv li{
    float: left;
    margin: 7px;
    list-style:none;
}
.group_pic {
    width: 60px;
    height: 60px;
}
.group_pic img,.device_pic img{
    width: 60px;
    height: 60px;
}
.member_name,.device_name{
    text-align: center
}
.addbotton{
    margin-top: 10px
}
.device_detail{
    width: 500px;
    height: 360px;
    background-color: white;
    border: 1px solid #ccc;
    position: absolute;
    top: 129px;
    left: 282px;
    
}
.device_detail_num li{
    float: left;
    list-style:none;
    margin: 7px;
    border: 1px solid black;
    width: 80px;
    height: 80px;
    position: relative;

}
.device_detail_tittle{
    height: 20px;
    background-color: white;
}
.device_detail_num{
    /* background-color: bisque; */
        height: 340px;
        overflow: auto ;
}
.device_pic{
    width: 80px;
    height: 60px;
    text-align: center;
}
.control_menum{
    background-color: white;
    position: absolute;
    height: 121px;
    width: 93px;
    top: 24px;
    left: 30px;
    
}
.control_voice,.control_vedio,.control_look,.control_text{
    font-size: 14px;
    text-align: center;
    height: 30px;
    line-height: 30px;
    cursor: pointer;
}
.select_add{
    width: 59px;
    height: 26px;
    line-height: 26px;
    text-align: center;
    cursor: pointer;
    margin-top: 25px;
    background-color: #409eff;
    color: white;
    border-radius: 5px 5px 5px 5px;
}
.editorial_group{
    height: 90px;
    width: 115px;
    position: absolute;
    text-align: center;
    top: 27px;
    background-color: white;
    left: 74px;
    z-index: 22;
    font-size: 12px;
    border: 1px solid black;
}
.ungroup,.modification,.editor_member{
    line-height: 30px;
    height: 30px
}
.ungroup:hover,.modification:hover,.editor_member:hover{
     background-color: #60a9e1;
}


</style>
