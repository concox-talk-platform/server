module.exports = {

    title: {
        logo_name: 'PTT后台管理系统',
        Account_center: '我的账户',
        control_center: '控制台',
        equipment: '设备',
        consumer: '用户',
        area: '区域',
        monitor: '监控平台'
    },
    account: {
        account_information: '账户信息',
        change_passwords: '修改密码',
        log_out: '退出登录'
    },
    button_message: {
        launcher: '登录',
        cancel: '取消',
        sign_up: '注册',
        ensure: '提交',
        add: '添加',
        confirm :'确定',
        reset: '重置'

        
    },
    reg_message: {
        enroll: '注册',
        title: '注册账号',
        account: '账号',
        name:'客户名称',
        pwd: '密码',
        cfm_pwd: '确认密码',
        account_type: '账号类型',
        dealer: '经销商',
        company: '公司',
        administrator: '管理员',
        dispatcher: '调度员',
        contact: '联系人',
        phone: '电话',
        email: '邮箱',
        adress: '地址',
        remark: '备注',
        account_info: '账号信息'
    },
    prompt_message: {
        account: '请输入账号',
        name: '请输入用户名称',
        pwd: '请输入密码',
        again_pwd: '请再次输入密码',
        pwd_length: '长度在 6 到 15个字符',
        account_type: '请选择账号类型',
        pwd_err: '两次输入密码不一致',
        account_error: '账号不存在',
        login_error: '密码错误',
        login_serve:'未与服务器建立连接',
        login_length:'密码长度不对',
        device_num: '请添加设备',
        subordinate: '请选择转移的用户'
   },
   change_pwd: {
       title: '修改密码',
       old_pwd:' 旧密码',
       new_pwd: '新密码',
       cfm_pwd: '确认密码',
       put_oldpwd: '请输入旧密码',
       put_newpwd: '请输入新密码',
       put_cfmpwd: ' 请再次输入新密码',
       put_newagain: '新密码不能与旧密码相同',
       change_success: '密码修改成功,请重新登录'
   },
   change: {
       language: '切换语言'
   },
   group: {
       add_group: '创建临时分组',
       add_name: '组名称',
       message: '请输入组名称',
       members: '临时组成员',
       member_title: '成员列表',
       all_member: '全体成员',
       select_member: '选中成员',
       add:"添加",
       remove:"移除",
       dbremove:"双击删除",
       dissolve:"解散本组",
       amend: '修改本组',
       modified_member: '成员编辑',
       dele_success: '删除成功',
       modify_success: '修改成功',
       modify_failed: '修改失败',
       name:'群名已存在',
       voice:'语音',
       video:'视频',
       messaging:'消息',


       
   },
           // 我的用户
    client_lang: {
        client_list: '客户列表',
        client_add: '添加',
        import: '导入新设备'

    },
    // 操作组
    control: {
        hint: '提示',
        delete_this: '确定解散此组？',
        Modify_group:'修改小组',
        group_num: '组号码'
        },
    // homepage页面
    ztree: {
        filter:"输入关键字进行搜索"
    },
    information:{
        login_name: '登录名 ', 
        type: '类型',
        number: '设备数量',
        contact: '联系人',
        phone:'电话',
        adress:'地址',
        equipment:"设备",
        data:  '资料'
    },
    table:{
        number:'序号',
        model: '型号',
        name: '设备名',
        time: '导入时间',
        operation :'操作',
        export: '转移',
        mass:'批量转移',
        device: '请选择转移的设备',
        info:'转移信息',
        message:'确定转出?',
        select: '请选择',
        new_device:'导入设备',
        lead:'确定导入?',
        no_data:"暂无数据",
     },
     failed:{
         transfer:'操作失败，请刷新页面重新操作',
         imei:'IMEI号的长度为15位数字，请确定后提交',
         export_success:'设备添加成功',
         transfer_success:'设备转移成功',

     },
     login:{
        expired:'登录时间过期，请重新登录',
    },
    establish:{
        failed:'创建失败，请重新创建',
        success:'创建成功',
    },
    registration:{
        name:'客户名要求：不能大于20位字符，只能包含英文字母、数字、汉字、下划线',
        client:'账号要求：5~20位字符，只能包含英文字母、数字、下划线',
        same_account:'账户名已存在',
        nick_name:'请填写客户名称',
        success:'资料更新成功',
    },
    video:{
        message:"是否发起视频通话?",
        audio:"是否发起语音通话?",
        audio_text:'语音通话',
        video_text:'视频通话',
        loding:'连接中......',
        call_message:'有视频呼叫来自'
    },
    im:{
        group_list:'群组列表',
        member_list:'成员列表',
        chat_list:'聊天列表',
        chat_history:'聊天记录',
        text:'文本',
        file:'文件',
        send:'发送',
        close:'关闭',
        media:'发送图片、视频、文件',
        tip:"消息提示",
        from:"收到来自",
        news:"的消息",
        answer:"回复",
    },
    apply:{
        request:'请求发起',
        connection:'连接',
        accept:'接受',
        refuse:'拒绝',
        agree:'对方同意了您的请求',
        call:'发起通话',
        reject:'对方拒绝了您的请求',
        offline:'对方不在线',
        exist:'对方不存在',
     }

}