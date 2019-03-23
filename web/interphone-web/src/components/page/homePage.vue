<template>
<div class="client">
    <div class="client_left">
        <div class="client_left_tittle">
            <i class="el-icon-caret-bottom client_left_icon"></i>
            <span class="client_left_name">{{ $t("client_lang.client_list") }} </span>
            <div class="client_left_regiter" @click="register">{{ $t("client_lang.client_add") }}</div>            
        </div>
    </div>
    <div class="client_right"></div>
  <!-- 注册 -->
        <el-dialog :title="$t('reg_message.title')" :visible.sync="registerVisible">
            <el-form ref="registerForm" :model="registerForm"  :rules="register_rules" label-width="136px" @submit.native.prevent>
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
                        this.$axios.post('/account',register_info)
                        .then(function (response) {
                        //  this.$router.push('/homePage');
                          window.console.log(response);
                          window.console.log(response.data.success);
                          if(response.data.success){
                                this.$message({
                                    
                                message: this.$t('prompt_message.register_succ'),
                                type: 'success'
                                });
                                this.registerVisible=false;
                                this.$refs['registerForm'].clearValidate();
                               this.$refs['registerForm'].resetFields();
                          }else{
                                this.$message({
                                message: '创建失败，请重新创建',
                                type: 'warning'
                                });
                          }
                        }.bind(this))
                        .catch( (error) => {
                         window.console.log(error);
                                this.$message({
                                message: '创建失败，请重新创建',
                                type: 'warning'
                                });
                    
                        }); 
          

                  } else {
                    // console.log('error submit!!');
                    return false;
                  }
                });
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

background-color: white;
}
.client_right{
 height: 100%;
 flex: 1;
 background-color: aqua
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
</style>

