<template>
<div class=" canvas_box">
    <vue-particle-line class="login-wrap">
        <div class="ms-login">
            <div class="ms-title">{{ $t("title.logo_name") }}</div>
            <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="0px" class="ms-content">
                <el-form-item prop="username">
                    <el-input v-model="ruleForm.username" :placeholder="$t('prompt_message.account')">
                        <i class="el-icon-tickets" slot="prepend"></i>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input type="password" :placeholder="$t('prompt_message.pwd')" v-model="ruleForm.password" @keyup.enter.native="submitForm('ruleForm')">
                        <i class="el-icon-edit" slot="prepend"></i>
                    </el-input>
                </el-form-item>
                <div class="login-btn">
                    <el-button type="primary" @click="submitForm('ruleForm')">{{ $t("button_message.launcher") }}</el-button>
                </div>
            </el-form>
            <span class="language"  @click="change_language">{{language}}</span>
        </div>
    </vue-particle-line>

</div>
</template>

<script>

    export default {
        data(){
            return {
                language:'切换语言',
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
                        this.language = '切换语言'
                    }else {
                        this.lang = 'zh-CN';
                        this.$i18n.locale = this.lang;//关键语句
                        // localStorage.setItem('language', '切换语言');
                        sessionStorage.setItem('language', '切换语言');
                        // localStorage.setItem('lang', 'zh-CN');
                            this.language = 'Change Language'
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
                            .catch( (error) => {
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
     .canvas_box{
         height: 100%;
     }
    .login-wrap{
        position: relative;
        width:100%;
        height:100%;
        /* background-image: url(../../assets/img/login-background.jpg); */
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
        z-index: 999;
        /* position: absolute;
        left:50%;
        top:50%; */
         position: absolute;
         margin:auto; 
         top: 0;
         left: 0;
         right: 0;
         bottom: 0;
        width:350px;
        height: 300px;
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