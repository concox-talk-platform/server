module.exports = {
    title: {
        logo_name: 'PTT Management System',
        Account_center: 'My Clients',
        control_center: 'Console',
        equipment: 'Device',
        consumer: 'User',
        area: 'Area',
        monitor: 'Monitor platform'
    },
    account: {
        account_information: 'Account Information',
        change_passwords: 'Change Passwords',
        log_out: 'Log Out'
    },
    button_message: {
        launcher: 'Log In',
       
        cancel: 'Cancel',
        sign_up: 'Sign Up',
        ensure: 'Submit',
        add: 'Add',
        confirm :'Confirm',
        reset: 'Reset'
        
    },
    reg_message: {
        enroll: 'Sign Up',
        title: 'Register Account',
        account: 'Account',
        pwd: 'Password',
        name:'Client Name',
        cfm_pwd: 'Confirm Password',
        account_type: 'Account Type',
        dealer: 'Dealer',
        company: 'Company',
        administrator: 'Administrator',
        dispatcher: 'Dispatcher',
        contact: 'Contact',
        phone: 'Phone',
        email: 'Email',
        adress: 'Adress',
        remark: 'Remarks',
        account_info: 'Account Information'
    },
    prompt_message: {
         account: 'Please enter an account',
         pwd: 'Please enter your password',
         name: 'Please type user name',
         again_pwd: 'Please enter your password again',
         pwd_length: 'The length should be between 6 and 15 characters',
         account_type: 'Please select account type',
         pwd_err: 'Entered passwords differ',
         account_error: 'Account does not exist',
         login_serve:'No connection was established to the server',
         login_error: ' Wrong password',
         login_length:'Password length incorrect',
         device_num: 'Please add device',
         subordinate: 'Please select a subordinate account'
        },
        change_pwd: {
            title: 'change Password',
            old_pwd:' Old Password',
            new_pwd: 'New Password',
            cfm_pwd: 'Confirm Password',
            put_oldpwd: 'Please importation of old passwords',
            put_newpwd: 'Please import new password',
            put_cfmpwd: ' Please enter your passwords again',
            put_newagain: 'The new password cannot be the same as the old password',
            change_success: 'The password has been modified successfully. Please log in again',
        },
        change: {
            language: 'Change Language'
        },
        group: {
            add_group: 'Create a temporary group',
            add_name: 'Group name',
            message: 'Please enter a group name',
            members: 'Temporary group member',
            member_title: 'Member List',
            all_member: 'All Members',
            select_member: 'Selected Member',
            add:"Add",
            remove:"Remove", 
            dbremove:"Double-click the delete",
            dissolve:"Dissolve ",
            amend: 'Modify',
            modified_member: 'Member editor',
            dele_success: 'Deleted successfully',
            modify_suc: 'Changed successfully',
            modify_failed: 'Changed failed',
            name:'The group name already exists',
            voice:'Voice',
            video:'Video',
            messaging:'Messaging',
     
        },
        // 我的用户
        client_lang: {
            client_list: 'Client list',
            client_add: 'Add',
            import: 'Import new device'
        },
        // 操作组
         control: {
          hint: 'Hint',
          delete_this: 'Are you sure to dissolve this group?',
          Modify_group:'Modify the group',
          group_num: 'Group id',
          group_name: 'Group name'
         },
    // homepage页面
        ztree: {
            filter:"Please enter a search phrase"
        },
        information:{
            login_name: 'Login name ',
            type: 'Type', 
            number:'Equipment number',
            contact: 'Contacts',
            phone:'Phone',
            adress:'Adress',
            equipment: 'Equipment',
            data:  'Personal data',
        },
        table:{
           number:'Number',
           model: 'Model',
           name: 'Device name',
           time: 'Import time',
           operation :'Operation',
           export: 'Transfer',
           mass:'Batch transfer',
           device: 'Please select a device',
           info:'Transinformation',
           message:'Are you sure to transfer out?',
           select: 'Please select',
           new_device:'Thse import equipment',
           lead:'Are you sure to import?',
           no_data:"No data",
        },
        failed:{
            transfer:'The operation failed. Please refresh the page and try again',
            imei:'IMEI is a number with a length of 15, please enter again',
            export_success:'Successfully added device',
            transfer_success:'Successful transfer device',
        },
        login:{
            expired:'Login time expired, please login again',
            
        },
        establish:{
            failed:'Failed to create. Please create again',
            success:'Created successfully',
        },
        registration:{
            name:'Customer name requirements: 1~20 characters, can only contain English letters, numbers, Chinese characters, underline',
            client:'Account requirements: 5~20 characters, can only contain English letters, numbers, underscores',
            same_account:'The account name already exists',
            nick_name: 'Please enter a customer name',
            success:'Data updated successfully',
        },
        video:{
            message:"Do you initiate a video call?",
            audio:"Do you initiate a voice call?",
            audio_text:'Voice Communication',
            video_text:'Video Call',
            loding:'Connecting......',
            call_message:'Incoming call from'
        },
        im:{
            group_list:'Group List',
            member_list:'Member List',
            chat_list:'Chat List',
            chat_history:'Chat Record',
            text:'Text',
            file:'File',
            send:'Send',
            close:'Close',
            media:'Send pictures, videos, files',
            tip:"Notification",
            from:"You have received a message from",
            news:" ",
            answer:"Reply",

        },
        apply:{
            request:'requests to initiate a',
            connection:'connection',
            accept:'Accept',
            refuse:'Refuse',
            agree:'The other party agreed to your request',
            call:'Initiate a call',
            reject:'The other party rejected your request',
            offline:' The other party is not online',
            exist:'The other party does not exist',
         }
    

        
}