<template>
<div>
<el-aside>
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
</div>
</el-aside>

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
                            <li  v-for="(item,index) in confirm_device_List" :key="item.id" v-on:dblclick="dele_device(index)" :title="$t('group.dbremove')">
                                <div class='group_pic' >
                                   <img src="../../assets/img/inter.png" alt="">
                                </div>
                                <div class="member_name">
                                    {{item.user_name}}
                                </div> 
                            </li>                                          
                       </ul>
                </div>
                <el-button  class="addbotton" @click="add_select_div">{{$t('button_message.add')}}</el-button>
            </div>
        </div>
        <el-dialog    style="text-align: center;"   width="42%" :title="$t('group.member_title')" :visible.sync="members_div" append-to-body>
            <!-- 添加新组左右列表移动 -->
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
            <div>
            <div class="transfer_cancel" @click="transfer_cancel">{{$t("button_message.cancel")}}</div>
            <div class="select_add" @click="selected_add">{{$t("button_message.confirm")}}</div>
            </div>
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
        </div>
        <div class="device_detail_num">
            <ul>
                <li  @contextmenu.prevent='control(index)'  v-for="(item,index) in select_group_num" :key="item.id">
                    <div class='device_pic'>
                        <img src="../../assets/img/inter.png" alt="">
                    </div>
                    <div class="device_name">
                        {{item.user_name}}
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
    <!-- 左侧删除组弹出框 -->
    <el-dialog :title="$t('control.hint')" :visible.sync="dialogVisible" width="30%">
        <span>{{$t('control.delete_this')}}</span>
        <span slot="footer" class="dialog-footer">
            <el-button @click="dialogVisible = false">{{$t('button_message.cancel')}}</el-button>
            <el-button type="primary" @click="dele_submit">{{$t('button_message.confirm')}}</el-button>
        </span>
    </el-dialog>
    <!-- 左侧修改组弹出框 -->
    <el-dialog :title="$t('control.Modify_group')" :visible.sync="amend_show" >
        <el-form :model="Modify_group_form">
            <el-form-item :label="$t('control.group_num')" label-width="110px">
            <el-input v-model="Modify_group_form.num" autocomplete="off" :disabled="dis_control"></el-input>
            </el-form-item>
            <el-form-item :label="$t('control.group_name')" label-width="110px">
            <el-input v-model="Modify_group_form.name" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="amend_show = false">{{$t('button_message.cancel')}}</el-button>
            <el-button type="primary" @click="amend_group_submit">{{$t('button_message.confirm')}}</el-button>
        </div>
    </el-dialog>
    <!-- 修改组成员弹出框 -->
    <el-dialog width="42%" style="text-align: center;" :title="$t('group.modified_member')" :visible.sync="modified_member_show" append-to-body >
        <!-- 左侧左右列表移动 -->
            <el-transfer
                style="text-align: left; display: inline-block"
                v-model="select_Data"
                :titles="[Source, Target]"
                :button-texts="[toleft, toright]"
                :format="{
                    noChecked: '${total}',
                    hasChecked: '${checked}/${total}'
                }"
                @change="handleChange"
                :props="{key: 'id',label:'name'}"
                :data="member_data">
                <!-- <el-button class="transfer-footer" slot="left-footer" size="small">操作</el-button>
                <el-button class="transfer-footer" slot="right-footer" size="small">操作</el-button> -->
            </el-transfer>
            <div>
             <div class="select_cancel" @click="modified_cancel">{{$t("button_message.cancel")}}</div>
             <div class="select_add" @click="modified_add">{{$t("button_message.confirm")}}</div>
            </div>
       
       
    </el-dialog>


</div>
</template>

<script>
export default {
    data() {   
        return {
                // local_device_list:JSON.parse(localStorage.getItem('device_list')),
                local_group_list:JSON.parse(localStorage.getItem('group_list')),
                yesData: [],
                select_Data:[],
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
                modified_member_show: false,
                dele_num:null,
                group_memberdiv_show: false,
                amend_show:false,
                dis_control:true,
                confirm_device:[],
                confirm_device_List:[],
                modified_add_member : [],
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
                Modify_group_form: {
                    num:'',
                    name: ''
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
          this.editorial_show = false.stop
        },
        // 鼠标右击
        control(index){
            window.console.log(index)
        },
        unfold(){
            this.show = !this.show;
        },
        groupbody_show(){
            this.members_show = !this.members_show;
            this.down_show = !this.down_show;
            this.right_show = !this.right_show;
        },
        // 点击变色事件
        group_detail(index){
            this.active = index ;
            this.select_group_nummber = index;
            this.group_div_show = true;
        },
        close_device_detail(){
            this.group_div_show = false;
        },
        // 移入移出
        enter_group(index){
            this.actived = index;     
        },
        leave_group(index){
            index = null;
            this.actived = index;
        },
        // 添加新的组
        group_add(){
            this.group_div = true;
        },
        add_select_div(){
            this.members_div = true;
        },
        // 列表移动
        handleChange(value, direction, movedKeys,) {
            window.console.log(value, direction, movedKeys,);
        },
        // 选择确认
        selected_add(){
            this.confirm_device = this.yesData;
            this.confirm_device_List = [];
            // let device_list = JSON.parse(localStorage.getItem('device_list'));
            // for( var i =0;i<device_list.length;i++){
            if(this.local_device_list == null){
                return 
            }else{
                for( var i =0;i<this.local_device_list.length;i++){
                    for(var j=0;j<this.confirm_device.length;j++){
                        // if(device_list[i].id == this.confirm_device[j]){
                        if(this.local_device_list[i].id == this.confirm_device[j]){
                        // this.confirm_device_List.push(device_list[i])  
                        this.confirm_device_List.push(this.local_device_list[i])  
                        }
                    }
                }
            }

            window.console.log(this.confirm_device_List)
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
       //添加新组提交
        group_add_submit(addgroup_form){
            let submit_form={};
            this.$refs[addgroup_form].validate((valid) => {
                if (valid) {
     
                    // submit_form.group_name = this.addgroup_form.name;
                    // submit_form.group_device = this.confirm_device_List;
                if(this.confirm_device_List.length !== 0){
                    let group_info={};
                    group_info.group_name = this.addgroup_form.name;
                    group_info.account_id = sessionStorage.getItem('id');
                    submit_form.group_info = group_info;
                    submit_form.device_infos = this.confirm_device_List;
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
    // 双击删除
        dele_device(qua){
            this.confirm_device_List.splice(qua,1);
            this.yesData.splice(qua,1);
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
        // 左侧删除组提交
        dele_submit(){
            let dele_mumber = this.dele_num;
            window.console.log(dele_mumber);
            // let dele_group_list = JSON.parse(localStorage.getItem('group_list'));
            // let dele_groupList = dele_group_list.slice(0,dele_group_list.length);
            let dele_groupList = this.local_group_list.slice(0,this.local_group_list.length);
            dele_groupList.splice(dele_mumber,1);
            window.console.log(dele_groupList);
            // this.$store.commit("groupList",dele_groupList);
            this.dialogVisible = false;
        },
        // 左侧修改组
        amend(){
            this.amend_show = true;
            let  amend_num = this.dele_num;
            // let amend_group_list = JSON.parse(localStorage.getItem('group_list'));
            // this.Modify_group_form.num = amend_group_list[amend_num].group_info.id
            this.Modify_group_form.num = this.local_group_list[amend_num].group_info.id
        },
        // 左侧修改组提交
        amend_group_submit(){
            if(this.Modify_group_form.name !== ''){
            let  new_group_info = {};
            new_group_info.id = this.Modify_group_form.num;
            new_group_info.name = this.Modify_group_form.name;
            this.amend_show = false;
            }else{
                this.$message({
                    message: this.$t('group.login_error'),
                    type: 'warning'
                });
            }
        },
        // 左侧成员编辑
        modified_member(){
            this.modified_member_show=true;
            let  modified_num = this.dele_num;
            window.console.log(modified_num);
            let modified_member_data =[];
            // let modified_group_list = JSON.parse(localStorage.getItem('group_list'));
            // modified_member_data = modified_group_list[modified_num].device_infos;
            modified_member_data = this.local_group_list[modified_num].device_infos;
            window.console.log(modified_member_data);
            let selected_member_data=modified_member_data.map(e =>{
                if(e.hasOwnProperty('id')){
                    return e.id
                }
            })
            this.select_Data=selected_member_data;
        },
        modified_cancel(){
            this.modified_member_show=false;
        },
        // 左侧成员编辑提交
        modified_add(){
            // let modified_divice = JSON.parse(localStorage.getItem('device_list'));
            // for( var i=0;i<modified_divice.length;i++){
                
            for( var i=0;i<this.local_device_list.length;i++){
                for(var j=0;j< this.select_Data.length;j++){
                    // if(modified_divice[i].id == this.select_Data[j]){
                    if(this.local_device_list[i].id == this.select_Data[j]){
                        // this.modified_add_member.push(modified_divice[i])
                        this.modified_add_member.push(this.local_device_list[i])
                    }
                }
            }

            window.console.log(this.modified_add_member);
            this.modified_member_show=false;
            //  this.group_div_show=false;
        },
        transfer_cancel(){

            this.confirm_device = this.yesData;
            this.members_div = false;
        }
    },
    computed:{
            group_list(){
                // return JSON.parse(localStorage.getItem('group_list'));
                return this.local_group_list;
            },
            device_member(){
                //  let device_device_list = JSON.parse(localStorage.getItem('device_list'));
                //  let device_id = device_device_list.map(e =>{
                    if(this.local_device_list == null){
                        return
                    }else{
                        let device_id = this.local_device_list.map(e =>{
                            if(e.hasOwnProperty('id')){
                                return e.id
                            }
                        })
                        return device_id;                   
                    }

            },
            noData(){
                // let  transfer_name = JSON.parse(localStorage.getItem('device_list'));
                let transfer_newData = [];
                // transfer_name.forEach((obj) => {
                    if (this.local_device_list == null){
                        return
                    }else{
                        this.local_device_list.forEach((obj) => {
                                    var  transfer_obj = {};
                                    transfer_obj.id = obj.id;
                                    transfer_obj.name =obj.user_name;
                                    transfer_newData.push(transfer_obj)
                                });
                        return  transfer_newData
                    }
      
            },
            // 选中组的设备信息
            member_data(){
                // let modified_name = JSON.parse(localStorage.getItem('device_list'));
                let modified_newData = [];
                // modified_name.forEach((obj) => {
                if( this.local_device_list == null){
                      return
                }else{
                    this.local_device_list.forEach((obj) => {
                        var  modified_obj ={};
                        modified_obj.id = obj.id;
                        modified_obj.name =obj.user_name;
                        modified_newData.push(modified_obj)
                    });
                    return  modified_newData
                }

            },
            select_group_num(){
                // let  group_select_device = JSON.parse(localStorage.getItem('group_list'));
                let  group_select_device = this.local_group_list;
                let  group_selected_num = this.select_group_nummber
                if( group_selected_num !== 'a'){
                    group_select_device = group_select_device[group_selected_num].device_infos
                        if(this.modified_add_member.length !== 0){
                            group_select_device =this.modified_add_member
                        }
                }
                return   group_select_device
            },
            // 组名
            device_group_name(){
                // let group_selected_member = JSON.parse(localStorage.getItem('group_list'));
                let group_selected_member = this.local_group_list;
                let  group_selected_name = this.select_group_nummber
                if( group_selected_name !== 'a'){
                    group_selected_member = group_selected_member[group_selected_name].group_info.group_name
                }
                return   group_selected_member
            },
            local_device_list(){
                return JSON.parse(localStorage.getItem('device_list'))
            }
    },
    mounted(){
    }
}
</script>

<style scoped>
.el-aside{
    overflow: visible
}
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
    height: 700px;
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
    margin-top: 64px;
    background-color: #409eff;
    color: white;
    border-radius: 5px 5px 5px 5px;
    display: inline-block;
    margin-left: 43px;
}
.select_cancel{
    width: 59px;
    height: 26px;
    line-height: 26px;
    text-align: center;
    cursor: pointer;
    margin-top: 25px;
    background-color: #ccc;
    color: white;
    border-radius: 5px 5px 5px 5px;
    display: inline-block;
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
.transfer_cancel{
      width: 59px;
    height: 26px;
    line-height: 26px;
    text-align: center;
    cursor: pointer;
    margin-top: 64px;
    background-color: #ccc;
    color: white;
    border-radius: 5px 5px 5px 5px;
    display: inline-block;
    margin-left: 7px;  
}


</style>
