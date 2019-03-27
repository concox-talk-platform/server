<template>
<div class="client">
    <div class="client_left">
        <div class="client_left_tittle">
            <i class="el-icon-caret-bottom client_left_icon"></i>
            <span class="client_left_name">{{ $t("client_lang.client_list") }} </span>
            <div class="client_left_regiter" @click="register">{{ $t("client_lang.client_add") }}</div>            
        </div>
        <div class="client_left_body">
            <el-input
            :placeholder="$t('ztree.filter')"
            v-model="filterText">
            </el-input>

            <el-tree
            class="filter-tree"
            :data="ztree_data"
            :props="defaultProps"
            default-expand-all
            :filter-node-method="filterNode"
            ref="ztree">
            </el-tree>

        </div>
    </div>
    <div class="client_right">
        <div class="client_details">
            <div class="account_info"><span class="account_info_tittle">{{$t('account.account_information')}}</span></div> 
            <div class="account_detailed_info">
                <!-- <div class="account_detailed_tittle"><span class="account_detailed_name">小明</span></div> -->
               <table>
                   <tbody>
                       <tr>
                           <td colspan="3">
                               <label class="account_detailed_name">{{information_name}}</label> 
                               <i class="el-icon-edit" @click="edit"></i>
                            </td>
                       </tr>
                       <tr>
                           <td> 
                               <label >{{$t('information.login_name')}}:</label>
                                <span>{{information_login}}</span>
                           </td>
                           <td>
                               <label >{{$t('information.type')}}:</label>
                                <span>{{information_type}}</span>
                           </td>
                           <td>
                               <label >{{$t('information.number')}}:</label>
                               <span>{{information_number}}</span>
                           </td>
                       </tr>
                       <tr>
                           <td>
                                <label>{{$t('information.contact')}}:</label>
                                <span id="account_detailed_contact">{{information_contact}}</span>
                           </td>
                           <td>
                               <label>{{$t('information.phone')}}:</label>
                               <span id="account_detailed_phone">{{information_phone}}</span>
                           </td>
                           <td>
                               <label>{{$t('information.adress')}}:</label>
                               <span id="account_detailed_addres">{{information_adress}}</span>
                           </td>

                       </tr>
                       
                   </tbody>
                </table> 
            </div> 
        </div>
        <div class="equipment_form">
            <el-tabs v-model="activeName"  @tab-click="handleClick">
                <el-tab-pane :label="$t('information.equipment')" name="first">
                    <div class="equipment_table">
                                <el-table ref="multipleTable" :data="tableData.slice((currentPage-1)*pagesize,currentPage*pagesize)" tooltip-effect="dark" style="width: 100%"  @selection-change="handleSelectionChange">
                                    <el-table-column type="selection" width="55" ></el-table-column>
                                    <el-table-column  type="index" width="80" :label="$t('table.number')"></el-table-column>
                                    <el-table-column prop="imei" label="IMEI" width="240"> </el-table-column>
                                    <el-table-column prop="version" :label="$t('table.model')" width="240" > </el-table-column>
                                    <el-table-column prop="device_name" :label="$t('table.name')" width="240" > </el-table-column>
                                    <el-table-column prop="time" :label="$t('table.time')" width="240"> </el-table-column>
                                    <el-table-column :label="$t('table.operation')">
                                        <template slot-scope="scope">
                                                    <el-button size="mini" @click="andexporth(scope.$index, scope.row)">{{$t('table.export')}}</el-button>

                                            </template>
                                    </el-table-column>
                                </el-table>
                                <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="[10, 20, 30, 40]"
                                :page-size="10" layout="total, sizes, prev, pager, next, jumper" :total="400">
                                </el-pagination>                        
                   
                    </div>
                </el-tab-pane>
                <el-tab-pane :label="$t('information.data')" name="second">
                    <div class="subordinate_div">
                            <el-form ref="subordinate" :model="subordinate"  label-width="136px" @submit.native.prevent>
                            <el-form-item :label="$t('reg_message.name')">
                                <el-input  v-model="subordinate.name" ></el-input>
                            </el-form-item>                         
                            <el-form-item :label="$t('reg_message.account')">
                                <el-input  v-model="subordinate.account"  :disabled="ban"></el-input>
                            </el-form-item>
                            <el-form-item :label="$t('reg_message.account_type')" prop="register_type">
                                    <el-radio-group v-model="subordinate.type">
                                    <el-radio v-for="item in Account_typedata" :key="item.Account_type"  :label="item.Account_type" :value="item.value" ></el-radio>       
                                    </el-radio-group>
                                </el-form-item>
                            <el-form-item :label="$t('reg_message.contact')">
                            <el-input v-model="subordinate.contact" autocomplete="off" ></el-input>
                            </el-form-item>
                            <el-form-item :label="$t('reg_message.phone')" >
                            <el-input v-model="subordinate.phone" autocomplete="off" ></el-input>
                            </el-form-item>
                            <el-form-item :label="$t('reg_message.email')">
                            <el-input v-model="subordinate.email" autocomplete="off" ></el-input>
                            </el-form-item>                
                            <el-form-item :label="$t('reg_message.adress')" >
                            <el-input v-model="subordinate.adress" autocomplete="off" ></el-input>
                            </el-form-item>
                        </el-form>
                        <div slot="footer" class="dialog-footer subordinate_footer">
                            <!-- <el-button @click="register_Cancle">{{$t('button_message.cancel')}}</el-button> -->
                            <el-button type="primary" @click="subordinate_submit">{{$t('button_message.confirm')}}</el-button>
                        </div>
                    </div>
                </el-tab-pane>
          
            </el-tabs>
        </div>
    </div>
  <!-- 注册 -->
        <el-dialog :title="$t('reg_message.title')" :visible.sync="registerVisible">
            <el-form ref="registerForm" :model="registerForm"  :rules="register_rules" label-width="136px" @submit.native.prevent>
                <el-form-item :label="$t('reg_message.name')" prop="register_name">
                    <el-input ref="register_name" v-model="registerForm.register_name" :placeholder="$t('prompt_message.name')"></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.account')" prop="register_Account">
                    <el-input ref="register_Account" v-model="registerForm.register_Account" :placeholder="$t('prompt_message.account')"></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.pwd')" prop="register_Password">
                    <el-input ref="register_Password" v-model="registerForm.register_Password" :placeholder="$t('prompt_message.pwd')" type="password"></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.cfm_pwd')" prop="register_cfmPassword">
                    <el-input v-model="registerForm.register_cfmPassword" :placeholder="$t('prompt_message.again_pwd')" type="password"></el-input>
                </el-form-item>
                  <el-form-item :label="$t('reg_message.account_type')" prop="register_type">
                        <el-radio-group v-model="registerForm.register_type">
                        <el-radio v-for="item in Account_typedata" :key="item.Account_type"  :label="item.Account_type" :value="item.value" ></el-radio>       
                        </el-radio-group>
                    </el-form-item>
                <el-form-item :label="$t('reg_message.contact')">
                  <el-input v-model="registerForm.name" autocomplete="off" ></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.phone')" >
                  <el-input v-model="registerForm.phone" autocomplete="off" ></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.email')">
                  <el-input v-model="registerForm.email" autocomplete="off" ></el-input>
                </el-form-item>                
                 <el-form-item :label="$t('reg_message.adress')" >
                  <el-input v-model="registerForm.adress" autocomplete="off" ></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="register_Cancle">{{$t('button_message.cancel')}}</el-button>
                <el-button type="primary" @click="submit_register('registerForm')">{{$t('button_message.sign_up')}}</el-button>
            </div>
        </el-dialog>     
</div>
    
</template>
<script>
export default {
    data() {
        var Same_check = (rule, value, callback) => {
            if (value === '') {
            callback(new Error(this.$t('prompt_message.again_pwd')));
            } else if (value !== this.registerForm.register_Password) {
            callback(new Error(this.$t('prompt_message.pwd_err')));
            } else {
            callback();
            }
        };
        return {
            registerVisible:false,
            registerForm: {
                    register_Account: '',
                    register_Password: '',
                    register_cfmPassword: '',
                    register_type: ''
                }, 
            register_rules: {
                register_name: [
                    { required: true, message: this.$t('prompt_message.name'), trigger: 'blur' }
                ], 
                register_Account: [
                    { required: true, message: this.$t('prompt_message.account'), trigger: 'blur' }
                ],             
                register_Password: [
                    { required: true, message: this.$t('prompt_message.pwd'), trigger: 'blur' },
                    { min: 6, max: 15, message:this.$t('prompt_message.pwd_length'), trigger: 'blur' }
                ],            
                register_cfmPassword: [
                    { required: true, message: this.$t('prompt_message.again_pwd'), trigger: 'blur' },
                    { min: 6, max: 15, message: this.$t('prompt_message.pwd_length'), trigger: 'blur' },
                    { validator: Same_check, trigger: 'blur' },
                ],
                register_type: [
                    { required: true, message: this.$t('prompt_message.account_type'), trigger: 'blur' }
                ]
            },
                // 账户类型
            Account_typedata:
                [
                    {
                        Account_type:this.$t('reg_message.dealer'),
                        value: 1,
                    },
                    {
                        Account_type:this.$t('reg_message.company'),
                        value: 2,
                    },
                    {
                        Account_type:this.$t('reg_message.administrator'),
                        value: 3
                    },
                    {
                        Account_type:this.$t('reg_message.dispatcher'),
                        value: 4
                    }
                ], 
            // 树组件数据
            filterText: '',
                    ztree_data: [{
                        id: 1,
                        label: '一级 1',
                        children: [{
                            id: 4,
                            label: '二级 1-1',
                            children: [{
                            id: 9,
                            label: '三级 1-1-1'
                            }, {
                            id: 10,
                            label: '三级 1-1-2'
                            }]
                        }]
                        }, {
                        id: 2,
                        label: '一级 2',
                        children: [{
                            id: 5,
                            label: '二级 2-1'
                        }, {
                            id: 6,
                            label: '二级 2-2'
                        }]
                        }, {
                        id: 3,
                        label: '一级 3',
                        children: [{
                            id: 7,
                            label: '二级 3-1'
                        }, {
                            id: 8,
                            label: '二级 3-2'
                        }]
                    }],
                    defaultProps: {
                    children: 'children',
                    label: 'label'
                    },
            // 登录信息
            information_name:'程涛',
            information_login:'小小',
            information_type:'经销商',
            information_number:'21',
            information_contact: '库卡',
            information_phone:'13112344321',
            information_adress:'china shenzhen',
        
            // 选项卡
             activeName: 'first',
            //  上级修改下级信息
            ban:true,
            subordinate :{
                name:'',
                account:'11111',
                contact: '',
                phone:'',
                email:'',
                type:'',
                adress:'',
                
            },
            // 设备表格
         tableData: [
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},
                    {imei: '21424152312321312',version: 'v.2.0',device_name: '测试组一',time: '2015/04/12'},    
        ],
        multipleSelection: [],
        currentPage: 1,
        pagesize:10,
             
        
        }
    },
    methods: {
            register(){
               this.registerVisible=true; 
               this.register_Account=''
               this.register_Password=''
               this.register_cfmPasswor=''
            },
            register_Cancle(){
                this.registerVisible=false;
                this.$refs['registerForm'].clearValidate();
                this.$refs['registerForm'].resetFields();
            },
            submit_register(registerForm){
                this.$refs[registerForm].validate((valid) => {
                  if (valid) {
                        let register_info={}
                        register_info.nikename=this.registerForm.register_name.trim();
                        register_info.username=this.registerForm.register_Account.trim();
                        register_info.pwd=this.registerForm.register_Password.trim();
                        switch(this.registerForm.register_type){
                                case '经销商':
                                register_info.roletype=1;
                                break;
                                case 'Dealer':
                                register_info.roletype=1;
                                break;                             
                                case '公司':
                                register_info.roletype=2;
                                break;
                                case 'Company':
                                register_info.roletype=2;
                                break;
                                case '管理员':
                                register_info.roletype=3;
                                break;  
                                case 'Administrator':
                                register_info.roletype=3;
                                break;                             
                                case '调度员':
                                register_info.roletype=4;
                                break;
                                case 'Dispatcher':
                                register_info.roletype=4;
                                break;
                            }
                        if(this.registerForm.name !== undefined){
                            register_info.name=this.registerForm.name.trim();
                        }else{
                            register_info.name=''
                        }
                        if(this.registerForm.phone !== undefined){
                            register_info.phone=this.registerForm.phone.trim();  
                        }else{
                            register_info.phone=''
                        }                                    
                        if(this.registerForm.email !== undefined){
                            register_info.email=this.registerForm.email.trim(); 
                        }else{
                            register_info.email=''
                        }                                       
                        if(this.registerForm.adress !== undefined){
                            register_info.adress=this.registerForm.adress.trim();  
                        }else{
                            register_info.adress=''
                        }        
                            window.console.log(register_info)                                      
                        // this.$axios.post('/account',register_info)
                        // .then(function (response) {
                        // //  this.$router.push('/homePage');
                        //   window.console.log(response);
                        //   window.console.log(response.data.success);
                        //   if(response.data.success){
                        //         this.$message({
                                    
                        //         message: this.$t('prompt_message.register_succ'),
                        //         type: 'success'
                        //         });
                        //         this.registerVisible=false;
                        //         this.$refs['registerForm'].clearValidate();
                        //        this.$refs['registerForm'].resetFields();
                        //   }else{
                        //         this.$message({
                        //         message: '创建失败，请重新创建',
                        //         type: 'warning'
                        //         });
                        //   }
                        // }.bind(this))
                        // .catch( (error) => {
                        //  window.console.log(error);
                        //         this.$message({
                        //         message: '创建失败，请重新创建',
                        //         type: 'warning'
                        //         });
                    
                        // }); 
                  } else {
                    // console.log('error submit!!');
                    return false;
                  }
                });
            },
            filterNode(value, data) {
                if (!value) return true;
                return data.label.indexOf(value) !== -1;
            },
            edit(){
            this.activeName='second';
            },
            // 选项卡
            handleClick(tab, event) {
                window.console.log(tab, event);
            },
            // 修改下级信息
            subordinate_submit(){
                let subordinate_info = {};
                subordinate_info.nikename = this.subordinate.name;
                subordinate_info.username = this.subordinate.account;
                subordinate_info.contact = this.subordinate.contact;
                subordinate_info.phone = this.subordinate.phone;
                subordinate_info.email = this.subordinate.email;
                subordinate_info.adress = this.subordinate.adress;
                switch(this.subordinate.type){
                    case '经销商':
                    subordinate_info.roletype=1;
                    break;
                    case 'Dealer':
                    subordinate_info.roletype=1;
                    break;                             
                    case '公司':
                    subordinate_info.roletype=2;
                    break;
                    case 'Company':
                    subordinate_info.roletype=2;
                    break;
                    case '管理员':
                    subordinate_info.roletype=3;
                    break;  
                    case 'Administrator':
                    subordinate_info.roletype=3;
                    break;                             
                    case '调度员':
                    subordinate_info.roletype=4;
                    break;
                    case 'Dispatcher':
                    subordinate_info.roletype=4;
                    break;
                }
                window.console.log(subordinate_info)
            }, 
            // 设备表格
                  handleSelectionChange(val) {
                    this.multipleSelection = val;
                    },
                    andexport(index, row) {
                        window.console.log(index, row);
                        // window.console.log(this.multipleSelection);                        
                    },
                    handleSizeChange(val) {
                        window.console.log(`每页 ${val} 条`);
                        // this.currentPage = currentPage;
                        this.pagesize = val;
                    },
                    handleCurrentChange(currentPage) {
                    //    window.console.log(`当前页: ${val}`);
                    this.currentPage =currentPage
                    }
           
    },
    watch: {
      filterText(val) {
        this.$refs.ztree.filter(val);
      }
    },
}
</script>
<style>
.client{
    display: flex;
    height: 100%;
}
.client_left{
width: 364px;
height: 100%;
border-right: 2px solid #a0b0c7;

background-color: white;
}
.client_right{
 height: 100%;
 flex: 1;
 background-color:white;
}
.client_left_tittle{
    height: 30px;
    background-color: #eef3f7;
    position: relative;
}
.client_left_icon{
    font-size: 29px;
}
.client_left_name{
    position: absolute;
    display: inline-block;
    font-size: 16px;
    top: 5px;
    left: 65px;
}
.client_left_regiter{
    display: inline-block;
    font-size: 12px;
    border: 1px solid #aab7c9;
    padding: 0px 10px 0px 10px;
    position: absolute;
    right: 13px;
    top: 8px;
    cursor: pointer;
}
.client_left_regiter:hover{
    color: #0072bd
}
.el-dialog{
    width: 35%;
}
/* 树组件盒子 */
.client_left_body{
    height: 100%;
    background-color: white;
}
.client_details{
    height: 139px;
    background-color: white;
    border:3px solid #cfdae7;
}
.equipment_form{
    padding-left: 10px;
    padding-right: 10px;
}
.account_info{
    height: 36px;
    line-height: 36px;
    background-color: #e3eaf0;
    border-bottom: 2px solid #d6d6d6;
}
.account_info_tittle{
    padding-left: 10px;
}
/* .account_detailed_info{
   width: 650px;
   height: 100px;
   border: 1px solid #d6d6d6;
   background-color: bisque;
}
.account_detailed_tittle{

} */
.account_detailed_info td{
        height: 32px;
    padding-right: 35px;
    padding-left: 10px;
    border: solid 1px #d6d6d6;
}
.account_detailed_info label{
     font-weight: bold;
    padding-left: 4px;
    font-size: 12px;
}
.account_detailed_info span{
 
    font-size: 14px;
}
.account_detailed_info table {
    border-collapse: collapse;
    border-spacing: 0;
}
.account_detailed_name{
    margin-right: 45px;
}
.subordinate_div{
    /* background-color: aqua; */
    width: 640px;
}
.subordinate_footer{
    text-align: center;
}
.equipment_table{
 height: 550px;
 overflow: auto;
}
</style>

