<template>
<div>
   <el-header height="96px">
    <div class="header">
        <div class="collapse-btn" >
            <i class="el-icon-menu"></i>
        </div>
        <div class="logo" @click="refresh">{{ $t("title.logo_name") }}</div>

        <div class="header-right">
            <div class="header-user-con">
                <!-- 全屏显示 -->
                <div class="btn-fullscreen" @click="handleFullScreen">
                    <el-tooltip effect="dark" :content="fullscreen?`取消全屏`:`全屏`" placement="bottom">
                        <i class="el-icon-rank"></i>
                    </el-tooltip>
                </div>
                <!-- 用户头像 -->
                <div class="user-avator"><img src="../../assets/img/img.jpg"></div>
                <!-- 用户名下拉菜单 -->
                <!-- <el-dropdown class="user-name" trigger="click" @command="handleCommand"> -->
                <el-dropdown class="user-name">
                    <span class="el-dropdown-link">
                        <!-- {{$store.state.User.loginName}} <i class="el-icon-caret-bottom"></i> -->
                        {{nikename}} <i class="el-icon-caret-bottom"></i>
                    </span>
                    <el-dropdown-menu slot="dropdown">
                            <!-- <el-dropdown-item  @click.native="FormVisible">{{ $t("account.account_information") }}</el-dropdown-item> -->
                            <el-dropdown-item  @click.native="changePassword">{{ $t("account.change_passwords") }}</el-dropdown-item>
                        <!-- <el-dropdown-item divided  command="loginout">退出登录</el-dropdown-item> -->
                        <el-dropdown-item  @click.native="change_lang">{{$t("change.language")}}</el-dropdown-item>
                        <el-dropdown-item  @click.native="outlogin">{{ $t("account.log_out") }}</el-dropdown-item>
                    </el-dropdown-menu>
                </el-dropdown>
     
            </div>
        </div>
         <v-sidebar></v-sidebar>
             <!-- 账户信息 -->
            <!-- <el-dialog :title="$t('reg_message.account_info')" :visible.sync="dialogFormVisible">
                <el-form :model="accountInfo" ref="accountInfo">
                    <el-form-item :label="$t('reg_message.contact')" :label-width="formLabelWidth">
                    <el-input v-model="accountInfo.name" autocomplete="off" ></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('reg_message.phone')" :label-width="formLabelWidth">
                    <el-input v-model="accountInfo.phone" autocomplete="off" ></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('reg_message.email')" :label-width="formLabelWidth">
                    <el-input v-model="accountInfo.email" autocomplete="off" ></el-input>
                    </el-form-item>                
                    <el-form-item :label="$t('reg_message.adress')" :label-width="formLabelWidth">
                    <el-input v-model="accountInfo.adress" autocomplete="off" ></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('reg_message.remark')" :label-width="formLabelWidth" >
                    <el-input type="textarea" v-model="accountInfo.remark" autocomplete="off" ></el-input>
                    </el-form-item>        
                </el-form>
                <div slot="footer" class="dialog-footer">
                    <el-button @click="formCancle">{{$t('button_message.cancel')}}</el-button>
                    <el-button type="primary" @click="updateInfor">{{$t('button_message.ensure')}}</el-button>
                </div>
            </el-dialog> -->
            <!-- 修改密码 -->
            <el-dialog :title="$t('change_pwd.title')" :visible.sync="pwdVisible" status-icon>
                <el-form ref="form" :model="form" :rules="rule" label-width="136px" @submit.native.prevent>
                    <el-form-item :label="$t('change_pwd.old_pwd')" prop="oldPassword">
                        <el-input ref="oldPassword" v-model="form.oldPassword" :placeholder="$t('change_pwd.put_oldpwd')" type="password"></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('change_pwd.new_pwd')" prop="newPassword">
                        <el-input ref="newPassword" v-model="form.newPassword" :placeholder="$t('change_pwd.put_newpwd')" type="password"></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('change_pwd.cfm_pwd')" prop="cfmPassword">
                        <el-input v-model="form.cfmPassword" :placeholder="$t('change_pwd.put_cfmpwd')" type="password"></el-input>
                    </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer">
                    <!-- <el-button @click="pwdVisible = false">取 消</el-button> -->
                    <el-button @click="pwdCancle">{{$t('button_message.cancel')}}</el-button>
                    <el-button type="primary" @click="submitPwd">{{$t('button_message.ensure')}}</el-button>
                </div>
            </el-dialog>
    </div>
   </el-header>
 </div>
</template>
<script>
import vSidebar from './Sidebar.vue';
    // import bus from '../web-api/bus';
    export default {
        data() {
            var validatePass = (rule, value, callback) => {
                if (value === '') {
                callback(new Error(this.$t('change_pwd.put_newpwd')));
                } else if (value == this.form.oldPassword) {
                callback(new Error(this.$t('change_pwd.put_newagain')));
                } else {
                callback();
                }
            };
            var validatePass2 = (rule, value, callback) => {
                if (value === '') {
                callback(new Error(this.$t('prompt_message.again_pwd')));
                } else if (value !== this.form.newPassword) {
                callback(new Error(this.$t('prompt_message.pwd_err')));
                } else {
                callback();
                }
            };
            return {
                // nikename:localStorage.getItem('loginName'),
                nikename:sessionStorage.getItem('loginName'),
                collapse: false,
                fullscreen: false,
                // name: 'linxin',
                // dialogFormVisible: false,
                lang:sessionStorage.getItem('lang'),
                pwdVisible:false,
                accountInfo: {
                    name:'',
                    phone:'',
                    email:'',
                    adress:'',
                    remark:''
                    
                },
                formLabelWidth: '80px',
            // 修改密码
               requesting: false,
                    form: {
                       
                    oldPassword: undefined,
                    newPassword: undefined,
                    cfmPassword: undefined
                    },
                    rule: {
                        oldPassword: [
                            {required: true, message: this.$t('change_pwd.put_oldpwd'), trigger: 'blur'},
                            {min: 6, max:15, message: this.$t('prompt_message.pwd_length'), trigger: 'blur'}
                        ],
                        newPassword: [
                            { validator: validatePass, trigger: 'blur' },
                            {required: true, message: this.$t('change_pwd.put_newpwd'), trigger: 'blur'},
                            {min: 6, max: 15, message: this.$t('prompt_message.pwd_length'), trigger: 'blur'},
                        ],
                        cfmPassword: [
                            { validator: validatePass2, trigger: 'blur' },
                            {required: true, message:this.$t('change_pwd.put_cfmpwd'), trigger: 'blur'},
                            {min: 6, max: 15, message: this.$t('prompt_message.pwd_length'), trigger: 'blur'},
                        ]
                    }
        
                }
        },
         components:{
            vSidebar
        },
        computed:{

        },
        methods:{
    
            // 切换语言
            change_lang(){
                 if ( this.lang === 'zh-CN' ) {
                        this.lang = 'en-US';
                        this.$i18n.locale = this.lang;//关键语句
                        // localStorage.setItem('lang', 'en-US');
                        sessionStorage.setItem('lang', 'en-US');
                        // localStorage.setItem('language', 'Change Language');
                        sessionStorage.setItem('language', 'Change Language');
                        window.console.log(this.lang)
                        this.$router.go(0);    
                    }else {
                        this.lang = 'zh-CN';
                        this.$i18n.locale = this.lang;//关键语句
                        // localStorage.setItem('language', '切换语言');
                        sessionStorage.setItem('language', '切换语言');
                        // localStorage.setItem('lang', 'zh-CN');
                        sessionStorage.setItem('lang', 'zh-CN');
                        this.$router.go(0);
                    }
            },
            session_language(){
                this.$i18n.locale = sessionStorage.getItem('lang');
            },
            changePassword(){
                this.pwdVisible=true;
            },
            pwdCancle(){
                this.pwdVisible = false;
                
            },
            submitPwd(){
            let sendpwd={};
            // sendpwd.id= localStorage.getItem('id');
            sendpwd.id= sessionStorage.getItem('id');
            // sendpwd.session_id =  localStorage.getItem('setSession_id');
            sendpwd.session_id =  sessionStorage.getItem('setSession_id');
            sendpwd.old_pwd=this.form.oldPassword;
            sendpwd.new_pwd=this.form.newPassword;
            sendpwd.confirm_pwd=this.form.cfmPassword;
            this.$axios.post('/account/pwd/update',sendpwd, 
            { headers: 
            {"Authorization" : sessionStorage.getItem('setSession_id')}
            }
            )
            .then(function (response) {
            window.console.log(response);
            this.$message({
                message: this.$t('change_pwd.change_success'),
                type: 'success'
            });       
             this.$router.push('/login'); 
             localStorage.clear();       
             sessionStorage.clear();     
              
             
            }.bind(this))
            .catch( (error) => {
            window.console.log(error);
            
            }); 
    
            },

            outlogin(){
            //  this.$router.push('/homePage'); 
             this.$router.push('/login');
             this.clearinfor()
            },
            // 清除数据
            clearinfor(){
            localStorage.clear();       
             sessionStorage.clear();
            },
            // 全屏事件
            handleFullScreen(){
                let element = document.documentElement;
                if (this.fullscreen) {
                    if (document.exitFullscreen) {
                        document.exitFullscreen();
                    } else if (document.webkitCancelFullScreen) {
                        document.webkitCancelFullScreen();
                    } else if (document.mozCancelFullScreen) {
                        document.mozCancelFullScreen();
                    } else if (document.msExitFullscreen) {
                        document.msExitFullscreen();
                    }
                } else {
                    if (element.requestFullscreen) {
                        element.requestFullscreen();
                    } else if (element.webkitRequestFullScreen) {
                        element.webkitRequestFullScreen();
                    } else if (element.mozRequestFullScreen) {
                        element.mozRequestFullScreen();
                    } else if (element.msRequestFullscreen) {
                        // IE11
                        element.msRequestFullscreen();
                    }
                }
                this.fullscreen = !this.fullscreen;
            },
            // 点击logo刷新
            refresh(){
                this.$router.go(0);
            }
        },

        beforeMount(){
            this.session_language()
        },
        mounted(){
 
        }
    }
</script>
<style>
.el-header {
    padding: 0px;
}
    .header {
    position: relative;
    -webkit-box-sizing: border-box;
    box-sizing: border-box;
    width: 100%;
    height: 96px;
    font-size: 22px;
    color: #fff;
    background-color: #206ba2;
    }
    .collapse-btn{
        float: left;
        padding: 0 21px;
        cursor: pointer;
        line-height: 70px;
    }
    .header .logo{
        float: left;
        width: 265px;
        line-height: 70px;
        cursor: pointer;
    }
    .header-right{
        float: right;
        padding-right: 50px;
    }
    .header-user-con{
        display: flex;
        height: 70px;
        align-items: center;
    }
    .btn-fullscreen{
        transform: rotate(45deg);
        margin-right: 5px;
        font-size: 24px;
    }
    .btn-bell, .btn-fullscreen{
        position: relative;
        width: 30px;
        height: 30px;
        text-align: center;
        border-radius: 15px;
        cursor: pointer;
    }
    .btn-bell-badge{
        position: absolute;
        right: 0;
        top: -2px;
        width: 8px;
        height: 8px;
        border-radius: 4px;
        background: #f56c6c;
        color: #fff;
    }
    .btn-bell .el-icon-bell{
        color: #fff;
    }
    .user-name{
        margin-left: 10px;
    }
    .user-avator{
        margin-left: 20px;
    }
    .user-avator img{
        display: block;
        width:40px;
        height:40px;
        border-radius: 50%;
    }
    .el-dropdown-link{
        color: #fff;
        cursor: pointer;
    }
    .el-dropdown-menu__item{
        text-align: center;
    }
    .el-dialog{
        width: 30%;
    }
</style>
