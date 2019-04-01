<template>
<!-- import { debug } from 'util'; -->
    <div class="login-wrap">
        <div class="ms-login">
            <div class="ms-title">{{ $t("title.logo_name") }}</div>
            <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="0px" class="ms-content">
                <el-form-item prop="username">
                    <el-input v-model="ruleForm.username" :placeholder="$t('prompt_message.account')">
                        <el-button slot="prepend" icon="el-icon-document"></el-button>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input type="password" :placeholder="$t('prompt_message.pwd')" v-model="ruleForm.password" @keyup.enter.native="submitForm('ruleForm')">
                        <el-button slot="prepend" icon="el-icon-edit"></el-button>
                    </el-input>
                </el-form-item>
                 <!-- <div class="register" @click="register">{{ $t("reg_message.enroll") }}</div> -->
                <div class="login-btn">
                    <el-button type="primary" @click="submitForm('ruleForm')">{{ $t("button_message.launcher") }}</el-button>
                </div>
            </el-form>
            <span class="language"  @click="change_language">{{language}}</span>
        </div>
  <!-- 注册 -->
        <!-- <el-dialog :title="$t('reg_message.title')" :visible.sync="registerVisible">
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
        </el-dialog>               -->
    </div>
</template>

<script>

    export default {
        data(){
            return {
                language:'Change Language',
                ruleForm: {
                    username: '',
                    password: '',
                },
                lang: 'en-US',
                // 规则
                rules: {
                    username: [
                        { required: true, message: this.$t('prompt_message.account'), trigger: 'blur' }
                    ],
                    password: [
                        { required: true, message:this.$t('prompt_message.pwd'), trigger: 'blur' }
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
                    ]
                    
            }
        },
        computed:{

        },
        methods: {
            // 切换语言
            change_language(){
                    if ( this.lang === 'zh-CN' ) {
                        this.lang = 'en-US';
                        this.$i18n.locale = this.lang;//关键语句
                        // localStorage.setItem('lang', 'en-US');
                        sessionStorage.setItem('lang', 'en-US');
                        // localStorage.setItem('language', 'Change Language');
                        sessionStorage.setItem('language', 'Change Language');
                        this.language = 'Change Language'
                    }else {
                        this.lang = 'zh-CN';
                        this.$i18n.locale = this.lang;//关键语句
                        // localStorage.setItem('language', '切换语言');
                        sessionStorage.setItem('language', '切换语言');
                        // localStorage.setItem('lang', 'zh-CN');
                            this.language = '切换语言'
                        sessionStorage.setItem('lang', 'zh-CN');
                    }
            },
            // 登录
            submitForm(ruleForm) {
            // submitForm() {
            //     this.$router.push('/homePage');
                this.$refs[ruleForm].validate((valid) => {
                    if (valid) {
                       
                        let login_key={};
                        login_key.username=this.ruleForm.username.trim();
                        login_key.pwd=this.ruleForm.password.trim();
                        window.console.log(login_key)
                        this.$axios.post('/account/login.do/'+this.ruleForm.username.trim(),login_key)
                        .then(function (response) {
                        // localStorage.setItem('setSession_id', response.data.session_id);
                        sessionStorage.setItem('setSession_id', response.data.session_id);
                        window.console.log( sessionStorage.getItem('setSession_id'));
                        sessionStorage.setItem('loginName', this.ruleForm.username.trim());  
                        this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
                            { headers: 
                            {
                            "Authorization" : sessionStorage.getItem('setSession_id')
                            }
                            })
                            .then((response) =>{
                             window.console.log(response);
                            sessionStorage.setItem('id', response.data.account_info.id);
                            sessionStorage.setItem('lang', this.lang);
                            sessionStorage.setItem('account_info', JSON.stringify(response.data.account_info));
                            localStorage.setItem('device_list', JSON.stringify(response.data.device_list));
                            localStorage.setItem('group_list', JSON.stringify(response.data.group_list));
                            localStorage.setItem('children_list', JSON.stringify(response.data.tree_data));
                            this.$router.push('/homePage');
                            })
                            .catch(function (error) {
                            window.console.log(error);
                            });          
                        //  this.$router.push('/homePage');
                        }.bind(this))
                        .catch( (error) => {
                        window.console.log(error.response.data);
                        if(error.response.data.error_code == '0021'){
                                this.$message({
                                message: this.$t('prompt_message.account_error'),
                                type: 'warning'
                                });
                        }else if(error.response.data.error_code == '0022'){
                                this.$message({
                                message: this.$t('prompt_message.login_error'),
                                type: 'warning'
                                });
                         }
                        }); 
                                           
                    } else {
                        window.console.log('error submit!!');
                        return false;
                    }
                });
            },
         
        }
    }
</script>
<style>
     .language{
            font-size: 14px;
            color: white;
            display: inline-block;
            float: right;
            margin-bottom: 14px;
            margin-right: 8px;
            cursor: pointer;
     }
    .login-wrap{
        position: relative;
        width:100%;
        height:100%;
        background-image: url(../../assets/img/login-background.jpg);
        background-size: cover;
        background-repeat: no-repeat;
    }
    .ms-title{
        width:100%;
        line-height: 50px;
        text-align: center;
        font-size:20px;
        color: #fff;
        border-bottom: 1px solid #ddd;
    }
    .ms-login{
        position: absolute;
        right:0px;
        top:50%;
        width:350px;
        margin:-190px 200px 0 0px;
        border-radius: 5px;
        background: rgba(255,255,255, 0.3);
        overflow: hidden;
    }
    .ms-content{
        padding: 30px 30px;
    }
    .login-btn{
        text-align: center;
    }
    .login-btn button{
        width:100%;
        height:36px;
        margin-bottom: 10px;
    }
    .login-tips{
        font-size:12px;
        line-height:30px;
        color:#fff;
    }
    .el-alert .el-alert--error{
      width: 50%
    }
   .el-select--small{
        width: 100%;
    }

</style>