<template>
<div style="height:100%">
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
                          <span class="group_on">  ({{online}} </span>
                          <span class="group_totalnum">  /{{totalnum}}) </span>
                     </div> 
                </div>
                 <div class="group_members" v-show='members_show' >
                       <div v-for="(item,index) in group_list" :key="item.group_info.id" class="members"
                        @click="group_detail(index)" 
                        @mouseover="enter_group(index)"
                        @mouseout="leave_group(index)"
                        :class="{active_on:active === index,active_hover:actived === index}"> 
                           <span class="interphonefamily groupicon">&#xe71b;</span>
                            <div class="members_num" >
                                <span class="members_name">{{item.group_info.group_name}}</span>
                                <span class="members_on">  ({{online}} </span>
                                <span class="members_totalnum">  /{{totalnum}}) </span>
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
                       <ul>
                            <li>
                                <div class='group_pic'>
                                   <img src="../../assets/img/computer.png" alt="">
                                </div>
                                <div class="member_name">
                                    2001
                                </div> 
                            </li>                                          
                                                                                     
                       </ul>
                </div>
                <el-button  class="addbotton" @click="members_div= true">{{$t('button_message.add')}}</el-button>
            </div>
        </div>
        <el-dialog
        width="31%" :title="$t('group.member_title')" :visible.sync="members_div" append-to-body>
        <!-- 左右列表移动 -->
            <el-transfer
                style="text-align: left; display: inline-block"
                v-model="value3"
                filterable
                :render-content="renderFunc"
                :titles="[Source, Target]"
                :button-texts="[toleft, toright]"
                :format="{
                    noChecked: '${total}',
                    hasChecked: '${checked}/${total}'
                }"
                @change="handleChange"
                :data="data">
                <el-button class="transfer-footer" slot="left-footer" size="small">操作</el-button>
                <el-button class="transfer-footer" slot="right-footer" size="small">操作</el-button>
            </el-transfer>
        </el-dialog>
        <div slot="footer" class="dialog-footer">
        <el-button @click="outerVisible = false">{{ $t("button_message.cancel") }}</el-button>
        <el-button type="primary" >{{ $t("button_message.ensure") }}</el-button>
        </div>
    </el-dialog>
    <!-- 组内详细设备 -->
    <div class="device_detail">
            <div class="device_detail_tittle">
                <!-- {{device_group_name}} -->
                Text
            </div>
            <div class="device_detail_num">
                <ul>
                        <li  @contextmenu.prevent='control(index)'  v-for="(item,index) in monise" :key="item.id">
                            <div class='device_pic'>
                                <img src="../../assets/img/computer.png" alt="">
                            </div>
                            <div class="device_name">
                                {{item.name}}
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



</div>
</template>

<script>
export default {
    data() {
      const generateData = _ => {
        const data = [];
        const cities = this.device_member;
         this.cities = this.device_member
        // const pinyin = ['shanghai', 'beijing', 'guangzhou', 'shenzhen', 'nanjing', 'xian', 'chengdu'];
        cities .forEach((city, index) => {
          data.push({
            label: city,
            key: index,
            // pinyin: pinyin[index]
          });
        });
        return data;
      };       
        return {
        data: generateData(),
        value3: [],
        value4: [],
        renderFunc(h, option) {
          return <span>{ option.key } - { option.label }</span>;
        },
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
         device_group_name:'',
         control_show:false,
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
         monizu:[{name:'a',id:1},{name:'b',id:2},{name:'class',id:3},{name:'d',id:4},{name:'e',id:5},{name:'f',id:6},{name:'g',id:7}],
         monise:[{name:'232',id:0},{name:'2132',id:1},{name:'2132',id:2},{name:'2132',id:3},{name:'2132',id:4},{name:'2132',id:5},{name:'2132',id:6},
         {name:'232',id:7},{name:'2132',id:8},{name:'2132',id:9},{name:'222',id:10},{name:'2132',id:11},{name:'2132',id:12},]

        }
    },
    methods: {
        // 鼠标右击
        control(index){
            alert(index)
        },
        unfold(){
            this.show=!this.show;
            window.console.log(this.device_member);

        },
        groupbody_show(){
            this.members_show=!this.members_show;
            this.down_show=!this.down_show;
            this.right_show=!this.right_show;
        },
        // 点击变色事件
        group_detail(index){
             this.active = index                  
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
        // 列表移动
      handleChange(value, direction, movedKeys) {
        window.console.log(value, direction, movedKeys);
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
/* .device_pic{
    width: 60px;
    height: 60px;
} */

</style>
