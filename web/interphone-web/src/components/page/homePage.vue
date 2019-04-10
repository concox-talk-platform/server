<template>
<div class="client">
    <el-container>
        <el-aside width="345px">
    <div class="client_left">
        <div class="client_left_tittle">
            <i class="el-icon-caret-bottom client_left_icon"></i>
            <span class="client_left_name">{{ $t("client_lang.client_list") }} </span>
            <div class="client_left_import" @click="device_import" v-if="rootShow">{{ $t("client_lang.import") }}</div>            
            <div class="client_left_regiter" @click="register">{{ $t("client_lang.client_add") }}</div>            
        </div>
        <div class="client_left_body">
            <!-- <el-input :placeholder="$t('ztree.filter')" v-model="filterText">
            </el-input> -->
            <!-- <el-tree @node-click="handleNodeClick"  class="filter-tree" :data="ztree_data" :props="defaultProps" default-expand-all :filter-node-method="filterNode" ref="ztree" :empty-text="$t('table.no_data')">
            </el-tree> -->
            <el-tree @node-click="handleNodeClick"  class="filter-tree"  :props="defaultProps" 
            lazy :load="get_newtree"   ref="ztree" :empty-text="$t('table.no_data')" v-if="tree_show">
            </el-tree>
        </div>
    </div>
    </el-aside>
    <el-main>
    <div class="client_right">
        <div class="client_details">
            <div class="account_info"><span class="account_info_tittle">{{$t('account.account_information')}}</span></div> 
            <div class="account_detailed_info">
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
                        <div class="transfer_tittle">
                            <span class="mass_transfer" @click="transfer">{{$t('table.mass')}}</span>
                        </div>
                        <!-- 完整分页 -->
                        <el-table ref="multipleTable" :data="table_page" tooltip-effect="dark"
                         :empty-text="$t('table.no_data')" style="width: 100%"  @selection-change="handleSelectionChange">
                               <!-- <el-table ref="multipleTable" :data="tableData" tooltip-effect="dark"
                         :empty-text="$t('table.no_data')" style="width: 100%"  @selection-change="handleSelectionChange"> -->
                            <el-table-column type="selection" width="55" ></el-table-column>
                            <el-table-column  type="index" width="80" :label="$t('table.number')"></el-table-column>
                            <el-table-column prop="imei" label="IMEI" width="240"> </el-table-column>
                            <el-table-column prop="bind_status.String" :label="$t('table.model')" width="240" > </el-table-column>
                            <el-table-column prop="user_name" :label="$t('table.name')" width="240" > </el-table-column>
                            <el-table-column prop="crate_time.String" :label="$t('table.time')" width="240"> </el-table-column>
                            <el-table-column :label="$t('table.operation')">
                                <template slot-scope="scope">
                                            <el-button size="mini" @click="device_export(scope.$index, scope.row)">{{$t('table.export')}}</el-button>
                                    </template>
                            </el-table-column>
                        </el-table>
                        <!-- 完整分页 -->
                        <!-- <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="[10, 20, 30, 40]"
                        :page-size="10" layout="total, sizes, prev, pager, next, jumper" :total="400" >
                        </el-pagination>       -->
                        <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="[10, 20, 30, 40]"
                        :page-size="10" layout="prev, pager, next" :total="page_mumber" >
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
                            <!-- <el-form-item :label="$t('reg_message.account_type')" prop="register_type">
                                    <el-radio-group v-model="subordinate.type">
                                    <el-radio v-for="item in Account_typedata" :key="item.Account_type"  :label="item.Account_type" :value="item.value" ></el-radio>       
                                    </el-radio-group>
                            </el-form-item> -->
                            <el-form-item :label="$t('reg_message.account_type')" prop="register_type">
                                    <el-input  v-model="subordinate.type"  :disabled="ban"></el-input>
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
                            <el-form-item :label="$t('reg_message.remark')"  >
                                <el-input type="textarea" v-model="subordinate.remark" autocomplete="off" ></el-input>
                            </el-form-item> 
                        </el-form>
                        <div slot="footer" class="dialog-footer subordinate_footer">
                            <el-button @click="reset ">{{$t('button_message.reset')}}</el-button>
                            <el-button type="primary" @click="subordinate_submit">{{$t('button_message.confirm')}}</el-button>
                        </div>
                    </div>
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
  <!-- 注册 -->
        <el-dialog :title="$t('reg_message.title')" :visible.sync="registerVisible" :show-close="false" width="33%">
            <el-form ref="registerForm" :model="registerForm"  :rules="register_rules" label-width="136px" @submit.native.prevent>
                <el-form-item :label="$t('reg_message.name')" prop="register_name">
                    <el-input ref="register_name" v-model="registerForm.register_name" :placeholder="$t('prompt_message.name')"></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.account')" prop="register_Account">
                    <el-input ref="register_Account" v-model="registerForm.register_Account" :placeholder="$t('prompt_message.account')"></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.pwd')" prop="register_Password">
                    <el-input ref="register_Password" v-model="registerForm.register_Password" :placeholder="$t('prompt_message.pwd')" ></el-input>
                </el-form-item>
                <el-form-item :label="$t('reg_message.cfm_pwd')" prop="register_cfmPassword">
                    <el-input v-model="registerForm.register_cfmPassword" :placeholder="$t('prompt_message.again_pwd')" ></el-input>
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
                <el-form-item :label="$t('reg_message.remark')"  >
                    <el-input type="textarea" v-model="registerForm.remark" autocomplete="off" ></el-input>
                </el-form-item> 
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="register_Cancle">{{$t('button_message.cancel')}}</el-button>
                <el-button type="primary" @click="submit_register('registerForm')">{{$t('button_message.sign_up')}}</el-button>
            </div>
        </el-dialog> 
  <!-- 转移 -->
        <el-dialog :title="$t('table.info')" :visible.sync="transfer_dialog" width="30%" :show-close="false">
            <el-select v-model="customer" filterable :placeholder="$t('table.select')">
                <!-- <el-option  v-for="item in customer_List" :key="item.id" :label="item.account_name" :value="item.id"> -->
                <el-option  v-for="item in customer_List" :key="item.id" :label="item.account_name+item.account_nickname" :value="item.id">
                    <span style="float:left">{{item.account_name}}</span>
                    <span style="float:right">{{item.account_nickname}}</span>
                </el-option>
            </el-select>
            <el-table :data="gridData" :empty-text="$t('table.no_data')">
                <el-table-column property="imei" label="IMEI"  width="200"></el-table-column>
                <el-table-column property="bind_status.String" :label="$t('table.model')"  width="150"></el-table-column>
                <el-table-column property="user_name" :label="$t('table.name')"  width="150"></el-table-column>
            </el-table>
        <span slot="footer" class="dialog-footer">
            <el-button @click="dialig_hidden">{{$t('button_message.cancel')}}</el-button>
            <el-button type="primary" @click="submit_transfer">{{$t('button_message.confirm')}}</el-button>
        </span>
        </el-dialog>
      <!-- 超级管理员导入设备 -->
      <el-dialog :title="$t('table.new_device')" :visible.sync="import_show" width="30%"  :show-close="false">
        <div class="imei_div">
            <el-button @click="add_imei">{{$t('group.add')}}</el-button>
            <el-table :data="imei_data" :empty-text="$t('table.no_data')">
                   
                    <el-table-column prop="iemi" label="IMEI">
                    <template  slot-scope="scope">
                    <el-input v-model.number="imei_data[scope.$index].iemi"   ref="iemi_rule"></el-input>
                    </template>
                    </el-table-column>
                    <el-table-column prop="command" :label="$t('table.operation')">
                    <template  slot-scope="scope">
                    <el-button @click="delete_imei(scope.$index)">{{$t('group.remove')}}</el-button>
                    </template>
                    </el-table-column>
               
            </el-table>

        </div>
        <span slot="footer" class="dialog-footer">
            <el-button @click="import_cancle">{{$t('button_message.cancel')}}</el-button>
            <el-button type="primary" @click="import_submit">{{$t('button_message.confirm')}}</el-button>
        </span>
        </el-dialog>
        
</el-main>


 </el-container>
</div>
    
</template>
<script>
export default {
    inject:['reload'],
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
            // personal_info:JSON.parse(localStorage.getItem('account_info')),
            personal_info:{},
            device_info:[] ,
            total_mumber:'',
            treeNode:1,
            nodeId:'',
            rootShow:false,
            ztreeId:sessionStorage.getItem('id'),
            // device_info:JSON.parse(localStorage.getItem('device_list')) ,
            // childre_infor:JSON.parse(localStorage.getItem('children_list')),
            
            // lang:sessionStorage.getItem('lang'/),
            registerVisible:false,
            registerForm: {
                    register_name:'',
                    register_Account: '',
                    register_Password: '',
                    register_cfmPassword: '',
                    register_type: '',
                    name:'',
                    phone:'',
                    email:'',
                    adress:'',
                    remark:''
                    
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
            // export_rules:{
            //   iemi_rule:[
            //        { min: 6, max: 15, message: this.$t('prompt_message.pwd_length'), trigger: 'blur' },
            //   ]
            // },
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
            ztree_data:[],
            defaultProps: {
                children: 'children',
                label: 'account_nickname'
            },
            tree_show :true,
            second_tree:'',
            // 登录信息
            information_name:'程涛',
            information_login:'小小',
            information_type:'',
            information_number:'',
            information_contact: '',
            information_phone:'',
            information_adress:'',
        
            // 选项卡
             activeName: 'first',
            //  上级修改下级信息
            ban:true,
            subordinate :{
                name:'',
                account:'',
                contact: '',
                phone:'',
                email:'',
                type:'',
                adress:'',
                remark:'',
                
            },
            // 设备表格
        // tableData:[],
        multipleSelection: [],
        currentPage: 1,
        pagesize:10,
        // 转移设备
        transfer_dialog: false,
        device_data:[],
        gridData:[],
        // 选择客户
        customer: '',
        // customer_List:[],
        device_imei:[],
        // 管理员导入设备
        import_show:false,
        // 导入imei号
        imei_data:[],  
        updata_list:[]
        }
    },
    methods: {
            register(){
               this.registerVisible=true; 
               this.register_Account=''
               this.registerForm.register_Password='123456'
               this.registerForm.register_cfmPassword='123456'
            },
            register_Cancle(){
                this.registerVisible=false;
                this.$refs['registerForm'].clearValidate();
                this.$refs['registerForm'].resetFields();
            },
            submit_register(registerForm){
                this.$refs[registerForm].validate((valid) => {
                  if (valid) {
                        let register_info={};
                        register_info.nick_name = this.registerForm.register_name.trim();
                        register_info.username = this.registerForm.register_Account.trim();
                        register_info.pwd = this.registerForm.register_Password.trim();
                        register_info.confirm_pwd = this.registerForm.register_cfmPassword.trim();
                        register_info.remark = this.registerForm.remark;
                        register_info.contact = this.registerForm.name;
                        register_info.pid = parseInt(sessionStorage.getItem('id'));
                        switch(this.registerForm.register_type){
                                case '管理员':
                                register_info.role_id=1;
                                break;  
                                case 'Administrator':
                                register_info.role_id=1;
                                break;                             
                                case '调度员':
                                register_info.role_id=2;
                                break;
                                case 'Dispatcher':
                                register_info.role_id=2;
                                break;                            
                                case '经销商':
                                register_info.role_id=3;
                                break;
                                case 'Dealer':
                                register_info.role_id=3;
                                break;                             
                                case '公司':
                                register_info.role_id=4;
                                break;
                                case 'Company':
                                register_info.role_id=4;
                                break;
                                case '超级管理员':
                                register_info.role_id=5;
                                break;
                                case 'Superadmin':
                                register_info.role_id=5;
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
                            register_info.address=this.registerForm.adress.trim();  
                        }else{
                            register_info.address=''
                        }
                          
                        window.console.log(register_info)   
                        this.$axios.post('/account',register_info,{ headers: 
                            {"Authorization" : sessionStorage.getItem('setSession_id')}
                                })
                        .then(function (response) {
                        //  this.$router.push('/homePage');
                          window.console.log(response);
                          window.console.log(response.data.success);
                          if(response.data.success){
                                this.$message({                                   
                                message: this.$t('establish.success'),
                                type: 'success'
                                });
                                this.registerVisible=false;
                                this.$refs['registerForm'].clearValidate();
                                 this.$refs['registerForm'].resetFields();
                                this.ztree_updata();

                           }else{
                                this.$message({
                                message: this.$t('establish.failed'),
                                type: 'warning'
                                });
                          }
                        }.bind(this))
                        .catch( (error) => {
                        switch(error.response.data.error_code){
                            case '0002':
                            this.$message({
                            message: this.$t('registration.name'),
                            type: 'warning'
                            });
                            break; 
                            case '0003':
                            this.$message({
                            message: this.$t('registration.client'),
                            type: 'warning'
                            });
                            break;
                            case '0005':
                            this.$message({
                            message: this.$t('registration.same_account'),
                            type: 'warning'
                            });
                            break;
                                
                        }
                    
                        }); 
                  } else {
                    // console.log('error submit!!');
                    return false;
                  }
                });
            },
            filterNode(value, data) {
                window.console.log(value)
                window.console.log(data)
                // if (!value) return true;
                // return data.label.indexOf(value) !== -1;
            },
            edit(){
            this.activeName='second';
            },
            // 选项卡
            handleClick(tab, event) {
                window.console.log(tab, event);
            },
            // 渲染个人信息
            apply_info(){
              this.personal_info = JSON.parse(sessionStorage.getItem('account_info'));             
              this.device_info = JSON.parse(localStorage.getItem('device_list'))
              this.renders();
            },
            renders(){     
                this.information_name = this.personal_info.nick_name;
                this.information_login = this.personal_info.username;
                this.information_phone = this.personal_info.phone.String;
                this.information_adress = this.personal_info.address.String;
                this.information_contact = this.personal_info.contact.String;
                if( this.device_info !== null){
                    this.information_number = this.device_info.length;
                }else{
                  this.information_number =0;
                 }
                
                if(sessionStorage.getItem('lang') == 'en-US'){
                     switch(this.personal_info.role_id){
                         case 1:
                         this.information_type = "Administrator";
                         this.subordinate.type = "Administrator";
                         break;
                         case 2:
                         this.information_type = "Dispatcher";
                         this.subordinate.type = "Dispatcher";
                         break;       
                         case 3:
                         this.information_type = "Dealer";
                         this.subordinate.type = "Dealer";
                         break;
                         case 4:
                         this.information_type = "Company";
                         this.subordinate.type = "Company";
                         break;
                         case 5:
                         this.information_type = "Superadmin";
                         this.subordinate.type = "Superadmin";
                         break;                                                                 
                     }
                }else{
                     switch(this.personal_info.role_id){
                         case 1:
                         this.information_type = "管理员";
                         this.subordinate.type = "管理员";
                         break;
                         case 2:
                         this.information_type = "调度员";
                         this.subordinate.type = "调度员";
                         break;       
                         case 3:
                         this.information_type = "经销商";
                         this.subordinate.type = "经销商";
                         break;
                         case 4:
                         this.information_type = "公司";
                         this.subordinate.type = "公司";
                         break;
                         case 5:
                         this.information_type = "超级管理员";
                         this.subordinate.type = "超级管理员";
                         break;                                                                   
                     }
                }
                window.console.log(this.personal_info);
                this.subordinate.adress = this.personal_info.address.String;
                this.subordinate.email = this.personal_info.email.String;
                this.subordinate.phone = this.personal_info.phone.String;
                this.subordinate.account = this.personal_info.username;
                this.subordinate.name = this.personal_info.nick_name;
                this.subordinate.remark = this.personal_info.remark.String;
                this.subordinate.contact = this.personal_info.contact.String;
                

                
                
            },
            // 重置
            reset(){
               this.subordinate.name = '';
               this.subordinate.contact = '';
               this.subordinate.phone = '';
               this.subordinate.email = '';
               this.subordinate.adress = '';
               this.subordinate.remark = '';
            },
            // 修改下级信息
            subordinate_submit(){
                let subordinate_info = {};
                subordinate_info.login_id = sessionStorage.getItem('id');
                subordinate_info.id = this.ztreeId.toString();
                subordinate_info.nick_name = this.subordinate.name;
                subordinate_info.username = this.subordinate.account;
                subordinate_info.type_id = this.personal_info.role_id.toString();
                subordinate_info.phone = this.subordinate.phone;
                subordinate_info.email = this.subordinate.email;
                subordinate_info.address = this.subordinate.adress;
                subordinate_info.remark = this.subordinate.remark;
                subordinate_info.contact = this.subordinate.contact;
                window.console.log(subordinate_info);
                let  same_name  = true;
                if(this.personal_info.nick_name == this.subordinate.name){
                    same_name  = true;
                }else{
                   same_name  = false; 
                }
                if(subordinate_info.nick_name == ''){
                  this.$message({
                        message: this.$t('registration.nick_name'),
                        type: 'warning'
                    });

                }else{
                    this.$axios.post('/account/info/update',subordinate_info,
                    { headers: 
                    {"Authorization" : sessionStorage.getItem('setSession_id')}
                    })
                    .then(() =>{
                    this.$message({
                        message: this.$t('registration.success'),
                        type: 'success'
                    });
                    if(same_name  == true){ 
                        this.update_data();
                    }else{
                        this.tree_show =false
                        this.ztree_name();
                        
                        this.update_data();
                     
                    }                   
                    })
                    .catch((error)=>{
                     window.console.log(error)
                    });
                }

        
            }, 
            // 设备表格
            handleSelectionChange(val) {
            this.multipleSelection = val;
            this.device_data=this.multipleSelection
            },
            device_export(index, row) {
                this.transfer_dialog=true;
                // this.customer_List = this.children_infor.children;
                this.gridData=[];
                this.gridData.push(row)
            },
            // 完整分页
            handleSizeChange(val) {
                // this.currentPage = currentPage;
                this.pagesize = val;
            },
            handleCurrentChange(currentPage) {
            this.currentPage =currentPage
            },
            transfer(){ 
                window.console.log(this.children_infor)
                window.console.log(this.updata_list.children)
                // this.update_data()
                if(this.multipleSelection.length !== 0){
                    this.transfer_dialog = true;
                    // this.customer_List = this.children_infor.children;
                    this.gridData=this.device_data;
                }else{
                    this.$message({
                        message: this.$t('table.device'),
                        type: 'warning'
                    });
                }
            },
            dialig_hidden(){
                this.transfer_dialog=false; 
                this.customer='';
                
            },
            // 转移设备
            submit_transfer(){
                this.$confirm(this.$t('table.message'),{
                   confirmButtonText: this.$t("button_message.confirm"),
                   cancelButtonText: this.$t("button_message.cancel")
                })
                .then(_ => {
                    window.console.log(_)
                    let transinformation={};
                    transinformation.devices=this.gridData;
                    let receiver ={}
                    this.customer_List.forEach(element => {
                        if(element.id ==this.customer){
                            receiver =element
                        }
                    });
                    window.console.log(receiver)
                    let facility={};
                    facility.account_id =receiver.id;
                    facility.account_name=receiver.account_name;
                    transinformation.receiver=facility;
                    window.console.log(transinformation)
                    if(transinformation.account_id ==''){
                        this.$message({
                        message: this.$t('prompt_message.subordinate'),
                        type: 'warning'
                        });
                        return false
                    }else{
                        this.$axios.post('/account_device/'+sessionStorage.getItem('id'),transinformation,
                            { headers: 
                            {
                            "Authorization" : sessionStorage.getItem('setSession_id')
                            }
                            })
                            .then(() =>{
                              this.transfer_dialog=false; 
                              this.customer='';
                              this.$message({
                                message: this.$t('failed.transfer_success'),
                                type: 'success'
                                });
                              this.update_data();
                              this.reload();
                            })
                            .catch(() => {
                                this.$message({
                                  message: this.$t('failed.transfer'),
                                  type: 'warning'
                                });
                            }); 
                    }
                   
                })
                .catch(_ => {
                    window.console.log(_)
                });
            },
            // 超级管理员导入设备
            device_import(){
                this.import_show=true;

            },
            import_cancle(){
                 this.import_show=false;
                 this.imei_data=[];
            },
            import_submit(){
                this.$confirm(this.$t('table.lead'),{
                   confirmButtonText: this.$t("button_message.confirm"),
                    cancelButtonText: this.$t("button_message.cancel")
                })
                .then(_ => {
                    window.console.log(_);
                    this.device_imei = this.imei_data.map(e =>{
                        if(e.hasOwnProperty('iemi')){
                            return e.iemi.toString()
                        }
                    })
                    let imei_length=[];
                    for(var i=0;i<this.device_imei.length;i++){
                           if(this.device_imei[i].length !== 15){
                                imei_length.push(this.device_imei[i]); 
                           }
                    }
                    if(imei_length.length == 0){
                       let decive_data = {};
                            decive_data.device_imei =this.device_imei;
                            this.imei_data=[];
                            this.import_show=false;
                            this.$axios.post('/device/import/SuperRoot',decive_data,
                            { headers: 
                            {
                            "Authorization" : sessionStorage.getItem('setSession_id')
                            }
                            })
                            .then(() =>{
                                this.transfer_dialog=false; 
                                this.customer='';
                                this.update_data();
                                this.reload();
                                this.$message({
                                message: this.$t('failed.export_success'),
                                type: 'success'
                                });
                            })
                            .catch(() => {
                                this.$message({
                                message: this.$t('failed.transfer'),
                                type: 'warning'
                                });
                            }); 
                    }else{
                        this.$message({
                        message: this.$t('failed.imei'),
                        type: 'warning'
                        });
                    }
                 
                    //  done();
                })
                .catch(_ => {
                    window.console.log(_)
                });
            },
            // 添加imei
            add_imei(){
              this.imei_data.push({})
            },
            delete_imei(index){
              this.imei_data.splice(index,1)
            },
           // ztree操作
             handleNodeClick(data) {
                 this.treeNode = data.$treeNodeId; 
                 this.nodeId = data.id;
                 if(data.$treeNodeId ==1){
                    this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
                    { headers: 
                    {
                    "Authorization" : sessionStorage.getItem('setSession_id')
                    }
                    })
                    .then((response) =>{
                    this.personal_info = response.data.account_info;
                    this.device_info =  response.data.device_list;
                    this.total_mumber = response.data.device_list.length;
                    this.ztreeId =  response.data.account_info.id;
                    this.renders(); 
                    })
                    .catch((error) => {
                    window.console.log(error);
                    });
                 }else{
                    this.$axios.get('/account_device/'+sessionStorage.getItem('id')+'/'+data.id,
                    { headers: 
                    {
                    // "Authorization" : localStorage.getItem('setSession_id')
                    "Authorization" : sessionStorage.getItem('setSession_id')
                    }
                    })
                    .then((response) =>{
                    //  localStorage.setItem('id', response.data.account_info.id);    
                    this.personal_info = response.data.account_info;
                    this.ztreeId =  response.data.account_info.id;
                    if(  response.data.devices == null){
                        this.total_mumber = 0;
                        this.device_info = [];
                    }else{
                       this.total_mumber = response.data.devices.length;
                       this.device_info =  response.data.devices;
                    }
                    
                    
                    this.renders();
                    // this.device_info=response.data.devices
                    })
                    .catch( (error) =>{
                    window.console.log(error);
                    });
                 } 
            } ,
            get_ztreeData(){
                this.ztree_data.push(this.children_infor);
                this.updata_list = this.children_infor;
            },
           // 获取下层树
            get_newtree(node, resolve) {
               if(node.level == 0){
                   resolve(this.tree_info)
                   window.console.log(node);
               }
               if(node.level == 1){
                   this.second_tree = this.tree_info[0].children
                   resolve(this.second_tree)
                      window.console.log(this.tree_info[0]);
               }
               if(node.level > 1){
                   window.console.log(node.data)
                    this.$axios.get('/account_class/'+sessionStorage.getItem('id')+'/'+node.data.id,
                    { headers: 
                    {
                    "Authorization" : sessionStorage.getItem('setSession_id')
                    }
                    })
                    .then((response) =>{
                     window.console.log(response.data.tree_data.children)
                     resolve(response.data.tree_data.children)
                    })
                    .catch( (error) => {
                    window.console.log(error);
                    resolve([])
                    });                    
                }
            },
            // 设备赋值
            get_device_data(){
                if( JSON.parse(localStorage.getItem('device_list')) == null){
                   return
                }else{
                   this.device_info = JSON.parse(localStorage.getItem('device_list'))
                }
                
            },
            get_total_mumber(){
                if(JSON.parse(localStorage.getItem('device_list'))==null){
                    this.total_mumber = 0;
                }else{
                 this.total_mumber = JSON.parse(localStorage.getItem('device_list')).length
                }
              
            },
            // 超级管理员导入
            root_Show(){
                if(JSON.parse(sessionStorage.getItem('account_info')).role_id == 5){
                    this.rootShow=true
                }else{
                    this.rootShow=false
                }
            },
            // 更新左侧树
            ztree_updata(){
                this.ztree_data = [];
                this.$axios.get('/account_class/'+sessionStorage.getItem('id')+'/'+sessionStorage.getItem('id'),
                { headers: 
                {"Authorization" : sessionStorage.getItem('setSession_id')}
                })
                .then((response) =>{
                this.ztree_data.push( response.data.tree_data);  
                this.updata_list = response.data.tree_data;      
                window.console.log(this.tree_info[0].children);
                let tree_length = this.tree_info[0].children.length-1
                this.second_tree.push(this.tree_info[0].children[tree_length]);
                localStorage.setItem('children_list', JSON.stringify(response.data.tree_data));
                })
                .catch(function (error) {
                window.console.log(error);
                });  
            },
            //更新左侧树名字
            ztree_name(){
                this.ztree_data = [];
                this.$axios.get('/account_class/'+sessionStorage.getItem('id')+'/'+sessionStorage.getItem('id'),
                { headers: 
                {"Authorization" : sessionStorage.getItem('setSession_id')}
                })
                .then((response) =>{
                this.second_tree=[]
                this.ztree_data.push( response.data.tree_data);  
                this.updata_list = response.data.tree_data;      
                localStorage.setItem('children_list', JSON.stringify(response.data.tree_data));
                // window.console.log(this.tree_info[0].children);
                // this.get_newtree(node, resolve)

                
                this.tree_show =true

                })
                .catch(function (error) {
                window.console.log(error);
                }); 
            },
            // addNode(node,data){
            //     this.node.childNodes = [];
            //     this.loadNode(this.node,this.resolve)
            // },
            // 页面赋值公共方法
            update_data(){
              if(this.treeNode == 1){
                this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
                    { headers: 
                    {
                    "Authorization" : sessionStorage.getItem('setSession_id')
                    }
                    })
                    .then((response) =>{
                    this.personal_info = response.data.account_info;
                    this.total_mumber = response.data.device_list.length;
                    this.device_info =  response.data.device_list;
                    // this.children_infor = response.data.tree_data;
                    this.renders();
                    sessionStorage.setItem('account_info', JSON.stringify(response.data.account_info));
                    localStorage.setItem('device_list', JSON.stringify(response.data.device_list));
                    localStorage.setItem('group_list', JSON.stringify(response.data.group_list)); 
                    localStorage.setItem('children_list', JSON.stringify(response.data.tree_data));         
                    })
                    .catch( (error) => {
                    window.console.log(error);
                    });
              }else{
                    this.$axios.get('/account_device/'+sessionStorage.getItem('id')+'/'+this.nodeId,
                    { headers: 
                    {
                    "Authorization" : sessionStorage.getItem('setSession_id')
                    }
                    })
                    .then((response) =>{
                    //  sessionStorage.setItem('id', response.data.account_info.id);    
                    this.personal_info = response.data.account_info;
                    if(response.data.devices == null){
                           this.total_mumber = 0;
                           this.device_info =  [];
                    }else{
                      this.total_mumber = response.data.devices.length;
                      this.device_info =  response.data.devices;
                    }
                    this.renders();
                    // this.device_info=response.data.devices
                    })
                    .catch((error) => {
                    window.console.log(error);
                    });
              }
            }

    },
    watch: {
            filterText(val) {
                this.$refs.ztree.filter(val);
            },
    },
    computed:{    
        tableData(){
            if(JSON.parse(localStorage.getItem('device_list')) == null){
              return
            }else{
              return this.device_info;
            }
           
        },
        page_mumber(){
            return this.total_mumber;
        },
        children_infor(){
            return JSON.parse(localStorage.getItem('children_list'))
        },
        table_page(){
            if(this.tableData == null){
               return
            }else{
               return  this.tableData.slice((this.currentPage-1)*10,this.currentPage*10)
            }
        },
        // 转移成员列表
        customer_List(){
            return  this.updata_list.children
        },
        tree_info(){
            return this.ztree_data
        }

    },
    beforeCreate() {
    },
    created(){
        this.apply_info();
        this.get_total_mumber();

    },
    beforeMount(){
       this. get_ztreeData();
       this.get_device_data();
       this.root_Show();
       
    //    this.get_total_mumber()
       
    }
 
}
</script>
<style>
.client{
    display: flex;
    height: 100%;
}
.client_left{
/* width: 364px; */
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
    left: 33px;
}
.client_left_regiter{
    display: inline-block;
    font-size: 12px;
    border: 1px solid #aab7c9;
    padding: 0px 10px 0px 10px;
    border-radius: 3px;
    position: absolute;
    right: 13px;
    top: 8px;
    cursor: pointer;
}
.client_left_regiter:hover{
    color: #0072bd
}
.client_left_import{
    display: inline-block;
    font-size: 12px;
    border: 1px solid #aab7c9;
    padding: 0px 10px 0px 10px;
    border-radius: 3px;
    position: absolute;
    right: 70px;
    top: 8px;
    cursor: pointer; 
}
.client_left_import:hover{
    color: #0072bd
}
.el-dialog{
    width: 35%;
}
/* 树组件盒子 */
.client_left_body{
    background-color: white;
    height: 743px;
    overflow: auto;
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
 height: 560px;
 overflow: auto;
}
.transfer_tittle{
    height: 27px;
    line-height: 18px;
    border-bottom: 1px solid #e4e7ed
}
.mass_transfer{
    font-size: 14px;
    padding-left: 5px;
    padding-right: 5px;
    border-radius: 3px;
    cursor: pointer;
    border: 1px solid #e4e7ed;
}
.mass_transfer:hover{
    color: #409eff;
}
.imei_div{
    max-height: 400px;
    overflow: auto;
}
</style>

