<template>
    <div @click="hide" id="big_box">
            <div  class="other_div"  v-if="other_show"  >
                <div style="text-align: center">{{$t('im.tip')}}</div>
                <span style="margin-left: 10px">{{$t('im.from')}}</span><span>{{other_name}}</span><span>{{$t('im.news')}}</span>
                <span style="display: none" id="other_name">{{other_name}}</span>
                <div class="other_answer" @click="look_other">{{$t('im.answer')}}</div>
            </div>
            <!-- SOS呼叫 -->
            <div class="sos_div" v-if="sos_show"  @click="look_sos">
                <span class="sos_title">SOS</span>
                <div class="sos_body">
                    <span>{{sos_name}}</span>&nbsp;
                    <span>{{$t('sos.alert')}}</span>
                </div>
                <div class="sos_place">{{$t('sos.place')}}</div>
                
            </div>
            <!-- 下线提醒 -->
            <div  class="off_line_div"  v-show="off_line_show">
                <div class="off_title">{{$t('im.off_line')}}</div>
                <div class="off_body">
                    <span>{{off_line_name}}</span> &nbsp; &nbsp;
                    <span>{{$t('im.off_span')}}</span>
                </div>
               <div class="off_foot">
                    <el-button type="primary" @click="off_line_button">{{$t('button_message.confirm')}}</el-button>
               </div>
            </div>
            <!-- 上线提醒 -->
            <div  class="off_line_div"  v-show="online_show">
                    <div class="off_title">{{$t('im.on_line')}}</div>
                    <div class="off_body">
                        <span>{{online_name}}</span> &nbsp; &nbsp;
                        <span>{{$t('im.on_span')}}</span>
                    </div>
                   <div class="off_foot">
                        <el-button type="primary" @click="online_button">{{$t('button_message.confirm')}}</el-button>
                   </div>
            </div>
            <!-- websocket断开提醒 -->
            <div  class="chat_div"  v-if="refresh_show">
                    <div class="chat_title">{{$t('im.chat')}}</div>
                   <div class="off_foot">
                        <el-button type="primary" @click="chat_refresh">{{$t('button_message.refresh')}}</el-button>
                   </div>
            </div>
            <!-- 地图切换 -->
            <div class="map_select">
                <div class="baidumap" :class="{map_active:mapshow}" @click="baidu_selected">{{$t('map.baidu')}}</div>
                <div class="googlemap" :class="{map_active:!mapshow}" @click="google_selected">{{$t('map.google')}}</div>
            </div>
    
    
        <!-- <intermap :postSOS='this.sosGps' ref="soschild" v-if="mapshow"></intermap> -->
        <intermap  ref="soschild" v-if="mapshow"></intermap>
        <intergooglemap v-if="!mapshow" ref="googlechild"></intergooglemap>
            <div id="monotor_content">
                    <el-aside>
                    <div   id="group_sider" style="height:100%"  @click="hide">
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
                                                <span>(</span>
                                                <span class="members_on"> {{item.group_info.online_num + 1}}</span>
                                                <span >  / </span>
                                                <span class="members_totalnum">{{item.device_infos.length}}</span>
                                                <span>)</span>
                                            </div>  
                                            <div class="editorial_group" v-show="editorial_show === index">
                                                <div class="ungroup" @click.stop="dissolve" >{{$t('group.dissolve')}}</div>
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
                    <!-- 视频 -->
                    <div  class="video_room" v-show="video_show" @mousedown.self="video_move"> 
                        <span class="video_tittle">{{$t('video.video_text')}}</span>
                            <div id="videocall">
                            <div id="videos" >
                                <div id="videoleft"></div>
                                <div id="videoright"></div>
                                </div>
                            </div>
                            <el-button type="danger" round @click="video_hangup" class="video_close">hangup</el-button>
                    </div> 
                    <div  class="audio_room" v-show="aduio_show" @mousedown.self="video_move"> 
                        <div class="audio_tittle" >{{$t('video.audio_text')}}</div>
                        <div class="audio_loding"  v-show="audio_logo"> {{$t('video.loding')}}</div>
                        <div id="audio_box" v-show="!audio_logo"> <img class="audio_img" src="../../assets/img/inter.png" alt=""> </div>
                            <div id="audiocall">
                            <div id="audios" class="hide">
                                <div id="audioleft"></div>
                                <div id="audioright"></div>
                                </div>
                            </div>
                            <el-button type="danger" round @click="audio_hangup" class="audio_close">hangup</el-button>
                    </div>
                    <!-- 即时聊天 -->
                    <div class="im_box" v-show="im_show" @mousedown.self="video_move">
                        <div class="im_aside">
                            <div class="im_aside_logo"><img class="im_aside_img" src="../../assets/img/computer.png" alt=""></div>
                            <div class="im_recent_icon" @click="select_imrecent" ><div class="interphonefamily im_size " v-bind:class="{aside_active:im_recent_show}">&#xe60d;</div></div>
                            <div class="im_device_icon" @click="select_imdevice"><div class="interphonefamily im_size"  v-bind:class="{aside_active:im_device_show}">&#xe66c;</div></div>
                            <div class="im_group_icon" @click="select_imgroup" ><div class="interphonefamily im_size" v-bind:class="{aside_active:im_group_show}">&#xe6d6;</div></div>
                        </div>
                        <div class="im_group">
                            <div v-show="im_recent_show">
                                <div class="im_group_list">{{ $t("im.chat_list") }}</div>   
                                <div class="im_overheigh"> 
                                        <div v-for="(item,index) in im_off_line" :key="item.id" class="im_members"
                                        @click="select_im_outline(index,item.id,item.Name,item.MsgReceiverType,item.imMsgData,item)" 
                                        @mouseover="enter_im_outline(index)"
                                        @mouseout="leave_im_outline(index)"
                                        :class="{active_on:im_line_active === index,active_hover:im_line_actived === index}"> 
                                        <span class="interphonefamily groupicon" v-if="item.MsgReceiverType == 1">&#xe66c;</span>
                                        <span class="interphonefamily groupicon" v-if="item.MsgReceiverType == 2">&#xe649;</span>
                                            <div class="members_num" >
                                                <span class="members_name">{{item.Name}}</span>
                                            </div>  
                                    </div>
                                </div>
                            </div>
                            <div v-show="im_device_show">
                                <div class="im_group_list">{{ $t("im.member_list") }}</div>   
                                <div class="im_overheigh">           
                                    <div v-for="(item,index) in local_device_list" :key="item.id" class="im_members"
                                        @click="select_im_device(index,item.id)" 
                                        @mouseover="enter_im_device(index)"
                                        @mouseout="leave_im_device(index)"
                                        :class="{active_on:im_device_active === index,active_hover:im_device_actived === index}"> 
                                            <span class="interphonefamily groupicon">&#xe66c;</span>
                                            <div class="members_num" >
                                                <span class="members_name">{{item.user_name}}</span>
                                            </div>  
                                    </div>
                                </div> 
                            </div>
                            <div v-show="im_group_show">  
                                <div class="im_group_list">{{ $t("im.group_list") }}</div>   
                                <div class="im_overheigh">              
                                    <div v-for="(item,index) in group_list" :key="item.group_info.id" class="im_members"
                                        @click="select_im_group(index,item.group_info.id)" 
                                        @mouseover="enter_im_group(index)"
                                        @mouseout="leave_im_group(index)"
                                        :class="{active_on:im_group_active === index,active_hover:im_group_actived === index}"> 
                                            <span class="interphonefamily groupicon">&#xe649;</span>
                                            <div class="members_num" >
                                                <span class="members_name">{{item.group_info.group_name}}</span>
                                            </div>  
                                        </div>
                                </div>  
                            </div>
                        </div>
                        
                        <span class="interphonefamily" style="font-size: 42px;">&#xe660;</span>
                        <span>{{talk_name}}</span>  
                        <span class="interphonefamily" style="font-size:18px;float:right;cursor:pointer;" @click="im_close">&#xe60b;</span>
                    
                        <div class="im_content" v-if="im_content_show">
                            <!-- 单人聊天 -->
                            <div class="im_message" v-if="im_self_show" >
                                <div  class="chat_box" v-for="(item) in im_total_message" :key="item.time_id">
                                    <div class="clearfix" :class='item.im_send_obj==1?"im_self":"im_opposite"'  >
                                        <div  :class='item.im_send_obj==1?"im_self_pic":"im_opposite_pic"'>
                                            <img v-show='item.im_send_obj==1'  class="im_self_img" src="../../assets/img/computer.png" alt=""> 
                                            <img  v-show='item.im_send_obj!==1' class="im_opposite_img" src="../../assets/img/inter.png" alt=""> 
                                        </div>
                                        <div :class='item.im_send_obj==1?"im_self_msg":"im_opposite_msg"'>
                                            <div  :class='item.im_send_obj==1?"im_self_news":"im_opposite_news"' v-show='item.MsgType==1'>{{item.ResourcePath}}</div>
                                            <div  :class='item.im_send_obj==1?"im_self_corner":"im_opposite_corner"'></div>
                                            <img  class="im_get_img" v-show='item.MsgType==2'  :src="item.ResourcePath" alt=""> 
                                            <div class="im_filebox" v-show='item.MsgType==3'>
                                                <audio :src="item.ResourcePath"  controls="controls"></audio>
                                                <!-- <video :src="item.ResourcePath" autoplay="autoplay" controls="controls" style="max-width: 190px;"></video> -->
                                            </div>
                                            <div class="im_filebox" v-show='item.MsgType==4'>
                                                <video :src="item.ResourcePath" autoplay="autoplay" controls="controls" style="max-width: 190px;"></video>
                                            </div>
                                            <div class="im_filebox" v-show='item.MsgType==5'>
                                                    <span class="im_file_text" style="float: left;" >{{item.file_name}}</span>
                                                    <a :href="item.ResourcePath" download="">
                                                            <span class="interphonefamily" style="float: right;cursor: pointer;margin-left:10px" >&#xe69c;</span>
                                                    </a>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <!-- 群主聊天 -->
                            <div class="im_message im_group_box"  v-if="im_record_show">
                                    <div  class="chat_box" v-for="(item) in im_total_group" :key="item.time_id">
                                        <div class="clearfix" :class='item.im_send_obj==1?"im_self":"im_opposite"'  >
                                            <div class="im_group_name"   v-if="item.im_send_obj !=1"  v-show="item.ReceiverType == 2 ">{{item.SenderName}}</div>
                                            <div  :class='item.im_send_obj==1?"im_self_pic":"im_opposite_pic"'>
                                                <img v-show='item.im_send_obj==1'  class="im_self_img" src="../../assets/img/computer.png" alt=""> 
                                                <img  v-show='item.im_send_obj!==1' class="im_opposite_img" src="../../assets/img/inter.png" alt=""> 
                                            </div>
                                            <div :class='item.im_send_obj==1?"im_self_msg":"im_opposite_msg"'>
                                                <div  :class='item.im_send_obj==1?"im_self_news":"im_opposite_news"' v-show='item.MsgType==1'>{{item.ResourcePath}}</div>
                                                <div  :class='item.im_send_obj==1?"im_self_corner":"im_opposite_corner"'></div>
                                                <img  class="im_get_img" v-show='item.MsgType==2'  :src="item.ResourcePath" alt="">
                                                <div class="im_filebox" v-show='item.MsgType==3'>
                                                        <audio :src="item.ResourcePath"  controls="controls"></audio>
                                                        <!-- <video :src="item.ResourcePath" autoplay="autoplay" controls="controls" style="max-width: 190px;"></video> -->
                                                </div> 
                                                <div class="im_filebox" v-show='item.MsgType==4'>
                                                    <video :src="item.ResourcePath" autoplay="autoplay" controls="controls" style="max-width: 190px;"></video>
                                                </div>
                                                <div class="im_filebox" v-show='item.MsgType==5'>
                                                        <span class="im_file_text" style="float: left;" >{{item.file_name}}</span>
                                                        <a :href="item.ResourcePath" download="">
                                                                <span class="interphonefamily" style="float: right;cursor: pointer;margin-left:10px" >&#xe69c;</span>
                                                        </a>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            <div class="im_media">
                                <label for="im_resource"><span class="interphonefamily im_media_icon" :title="$t('im.media')">&#xe7d8;</span></label>
                                <span class="interphonefamily im_media_icon" @click="map_show" >&#xe643;</span>
                                <span class="interphonefamily im_media_icon" @click="im_audio" :title="$t('video.audio_text')"   v-show="im_self_show">&#xe621;</span>
                                <span class="interphonefamily im_media_icon" @click="im_video" :title="$t('video.video_text')" v-show="im_self_show">&#xe616;</span>
                                <span class="interphonefamily im_media_icon" @click="im_audio_bridge" :title="$t('talkback.audio_bridge')" v-show="!im_self_show">&#xe6a5;</span>
                                <div class="im_record"  @click="im_unfold_history">
                                <span class="interphonefamily im_record_icon" >&#xe60d;</span>
                                <span>{{ $t("im.chat_history") }}</span>
                                <i class="el-icon-caret-right"></i>
                                </div>
                            </div>
                            <div class="im_send">
                                <!-- <div    class="im_send_message"  contenteditable="true"  ></div> -->
                                <textarea class="im_send_message" v-model="im_send_news" v-if="im_box_show" @keyup.enter="im_send_subimt" ></textarea>
    
                                <div class="upload_file" v-show="upload_show">
                                    <form id="im_formid">
                                    <input type="file" id="im_resource" name="resource" ref="resource" @change="filechange"  style="display: none">
                                    <!-- <span >{{file_name}}</span> -->
                                    <!-- <div class="dele_file">删除文件</div> -->
                                        <img src=""  id="preview_img" alt="">
                                        <textarea id="dele_area"  @keyup.delete="dele_file" readonly="readonly"></textarea>
                                    <!-- <button type="button" @click="dele_file">删除文件</button> -->
                                    
                                </form> 
                                </div>
                                <!-- <input v-show="fileupView" id="fileSelect" name="fileSelect"  ref="inputer"  type="file"/> -->
                                <div class="im_send_button">
                                    <span class="im_send_close" @click="im_send_close">{{ $t("im.close") }}</span>
                                    <!-- <span @click="duankai">   断开</span>
                                    <span @click="lianjie">lianjie</span> -->
                                    <span class="im_send_subimt" @click="im_send_subimt">{{ $t("im.send") }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="im_talkback" v-if="im_talkback_show"> 
                              <!-- <span>点击加入群对讲</span> -->
                              <span class="out_talkback"  @click="out_talkback">{{ $t("talkback.out") }}</span>
                              <div class="talkactive_tip" v-show="mouse_show">{{ $t("talkback.down") }}</div>
                              <div class="talkactive_tip" v-show="!mouse_show">{{ $t("talkback.up") }}</div>
                              <div class="talkback_div" @mousedown="talk_active" @mouseup="talk_actived" :class="{talkback_active:talkactive}" :title="$t('talkback.touch')">
                                    <span class="interphonefamily talkback_icon">&#xe6a5;</span>
                              </div>


                        </div>
                        <div class="im_history" :style="'right:'+history_right+';'+'z-index:'+history_index">
                            <div class="im_history_top">
                                <div class="clearfix">
                                    <i class="interphonefamily"> &#xe60d;</i>
                                    <span  class="im_history_title">{{ $t("im.chat_history") }}</span>
                                    <i class="interphonefamily" @click="im_collapse_history" style="float:right;cursor:pointer;"> &#xe60e;</i>
                                </div>
                                <div class="clearfix">
                                    <div class="im_history_text" @click="select_im_text" v-bind:class="{active_history:im_select_text}">{{ $t("im.text") }}</div>
                                    <div class="im_history_file" @click="select_im_file" v-bind:class="{active_history:im_select_file}">{{ $t("im.file") }}</div>
                                    <div class="im_history_talk" @click="select_im_talk" v-if="im_record_show"  v-bind:class="{active_history:im_select_talk}">{{ $t("im.talk") }}</div>
                                </div>
    
                            </div>
                            <div class="im-history-wrap">
                                <div  v-if="select_text">
                                    <div v-for="(item) in im_total_history" :key="item.time_id">
                                        <div class="im_history_info"  v-if='item.MsgType ==1'>
                                            <div class="im_history_name clearfix">
                                                <div class="im_record_name">{{item.SenderName}}</div>
                                                <div class="im_record_time">{{item.SendTime}}</div>
                                            </div>
                                            <div class="im_history_comtent">{{item.ResourcePath}}</div>
                                        </div>
                                    </div>
                               </div>
                               <div  v-if="select_file">
                                    <div    v-for="(item) in im_total_history" :key="item.time_id">
                                        <div v-if='item.MsgType !=1' class="im_history_info">
                                            <div class="im_history_name clearfix">
                                                <div class="im_record_name">{{item.SenderName}}</div>
                                                <div class="im_record_time">{{item.SendTime}}</div>
                                            </div>
                                            <img  class="im_get_img" v-if='item.MsgType==2'  :src="item.ResourcePath" alt=""> 
                                            <div class="im_filebox" v-if='item.MsgType==4'>
                                                <video :src="item.ResourcePath" autoplay="autoplay" controls="controls" style="max-width: 190px;"></video>
                                            </div>
                                            <div  v-if='item.MsgType==5' style="height: 30px" >
                                                    <span class="im_file_text" style="float: left;" >{{item.file_name}}</span>
                                                    <a :href="item.ResourcePath" download="">
                                                            <span class="interphonefamily" style="float: right;cursor: pointer;margin-left:10px" >&#xe69c;</span>
                                                    </a>
                                            </div>
                                       </div>
                                    </div>
                               </div>
                               <div  v-if="select_talk">
                                    <!-- <div v-for="(item) in im_total_history" :key="item.time_id"> -->
                                            <!-- <div class="im_history_info"  v-if='item.MsgType ==1'> -->
                                                <!-- <div class="im_history_name clearfix">
                                                    <div class="im_record_name">{{item.SenderName}}</div>
                                                    <div class="im_record_time">{{item.SendTime}}</div>
                                                </div>
                                                <div class="im_history_comtent">{{item.ResourcePath}}</div>
                                            </div> -->
                                        <!-- </div> -->
                                        <div class="im_history_info" v-for="(item,index) in ssaudio" :key="item.id" >
                                                <div class="im_history_name clearfix">
                                                        <div class="im_record_name">{{item.name}}</div>
                                                        <div class="im_record_time">{{item.time}}</div>
                                                </div>
                                                <div class="im_talk_content clearfix">
                                                    <div class="im_voice-wrap">
                                                        <div class="im_voice_div">
                                                            <img class="im_voice_img" :filename="item.file" 
                                                            :src="pause_img" alt="" @click="play_audio(item,index)"  v-show="audio_play">
                                                            <img  class="im_voice_gif"  :src="play_img" alt="" v-show="img_gif===index">
                                                            <!-- <audio src="../../assets/img/qwe.mp3"></audio> -->
                                                        </div>
    
                                                        <span class="im_voice_num">{{item.long}}</span>
                                                    </div>
                                                </div>
                                            </div>
                                           <div class="im_history_info"  >
                                                <div class="im_history_name clearfix">
                                                        <div class="im_record_name">leikun</div>
                                                        <div class="im_record_time">2019年04月29日 19:07:46</div>
                                                    </div>
                                                    <!-- <div class="im_history_comtent">{{item.ResourcePath}}</div> -->
                                            </div>
                                            <div class="im_history_info"  >
                                                    <div class="im_history_name clearfix">
                                                            <div class="im_record_name">leikun</div>
                                                            <div class="im_record_time">2019年04月29日 20:07:46</div>
                                                        </div>
                                                        <!-- <div class="im_history_comtent">{{item.ResourcePath}}</div> -->
                                            </div>
                               </div>                          
                            </div>
                        </div>
                    </div>
                    <!-- 创建组 -->
                    <el-dialog :title="$t('group.add_group')" :visible.sync="group_div" :show-close="false">
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
                        <el-dialog    style="text-align: center;"   width="42%" :title="$t('group.member_title')" :visible.sync="members_div" append-to-body :show-close="false">
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
                    <div class="device_detail" v-if="group_div_show"  @click="hide"  >
                        <div class="device_detail_tittle">
                            <span class="device_detail_name">{{device_group_name}}</span>
                            <span class="device_detail_close" @click="close_device_detail">x</span>
                        </div>
                        <div class="device_detail_num" >
                            <ul>
                                <li>       
                                    <div class='device_pic' >
                                        <img src="../../assets/img/computer.png" alt="">
                                    </div>
                                    <span class="device_name" :title="computer_name">
                                    {{computer_name}}
                                    </span> 
                                </li>
                                <li  @contextmenu.prevent='media_control(index,item)'  v-for="(item,index) in select_group_num" :key="item.id" :class="{off_color:item.online == 1}">
                                    <div>
                                            <div class='device_pic' >
                                                    <img src="../../assets/img/inter.png" alt="">
                                                </div>
                                                <span class="device_name" :title="item.user_name">
                                                    {{item.user_name}}
                                                </span> 
                                                <div class="control_menum" v-if="media_show === index">
                                                    <div class="control_voice" @click="audio_begin(item)">{{ $t("group.voice") }}</div>
                                                    <div class="control_vedio" @click="video_begin(item)">{{ $t("group.video") }}</div>
                                                    <!-- <div class="control_look">视频查看</div> -->
                                                    <div class="control_text" @click="text_begin(item)">{{ $t("group.messaging") }}</div>
                                                </div>
    
                                    </div>
    
                                </li>                             
                            </ul>
                        </div>
                    </div>
                    <!-- 左侧删除组弹出框 -->
                    <el-dialog :title="$t('control.hint')" :visible.sync="dialogVisible" width="30%" :show-close="false">
                        <span>{{$t('control.delete_this')}}</span>
                        <span slot="footer" class="dialog-footer">
                            <el-button @click="dialogVisible = false">{{$t('button_message.cancel')}}</el-button>
                            <el-button type="primary" @click="dele_submit">{{$t('button_message.confirm')}}</el-button>
                        </span>
                    </el-dialog>
                    <!-- 左侧修改组弹出框 -->
                    <el-dialog :title="$t('control.Modify_group')" :visible.sync="amend_show"  :show-close="false">
                        <el-form :model="Modify_group_form">
                            <el-form-item :label="$t('control.group_num')" label-width="110px">
                            <el-input v-model="Modify_group_form.num" autocomplete="off" :disabled="dis_control"></el-input>
                            </el-form-item>
                            <el-form-item :label="$t('control.group_name')" label-width="110px">
                            <el-input v-model="Modify_group_form.name" autocomplete="off"></el-input>
                            </el-form-item>
                        </el-form>
                        <div slot="footer" class="dialog-footer">
                            <el-button @click="amend_cancel">{{$t('button_message.cancel')}}</el-button>
                            <el-button type="primary" @click="amend_group_submit">{{$t('button_message.confirm')}}</el-button>
                        </div>
                    </el-dialog>
                    <!-- 修改组成员弹出框 -->
                    <el-dialog width="42%" style="text-align: center;" :title="$t('group.modified_member')" :visible.sync="modified_member_show" append-to-body  :show-close="false">
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
                            </el-transfer>
                            <div>
                            <div class="select_cancel" @click="modified_cancel">{{$t("button_message.cancel")}}</div>
                            <div class="select_add" @click="modified_add">{{$t("button_message.confirm")}}</div>
                            </div>
                    </el-dialog>
                    <!-- 加入对讲弹出框 -->
                    <el-dialog :title="$t('control.hint')" :visible.sync="talkbackVisible" width="30%" :show-close="false">
                            <span>{{$t('talkback.talkback')}}</span>
                            <span slot="footer" class="dialog-footer">
                                <el-button @click="talkbackVisible = false">{{$t('button_message.cancel')}}</el-button>
                                <el-button type="primary" @click="talkback_submit">{{$t('button_message.confirm')}}</el-button>
                            </span>
                        </el-dialog>

        
             </div>
            <!-- 语音or视频请求 -->
        <!-- 背呼方 -->
        <div v-show="video_apply" class="video_apply" >
            <div class="video_apply_title">
                <span>{{your_callname}}</span> &nbsp;<span>{{$t('apply.request')}}</span>&nbsp;<span>{{your_calltype}}</span>&nbsp;<span>{{$t('apply.connection')}}</span>
            </div>
            <div class="video_apply_body">
                    <el-button type="info" @click="refuse">{{$t('apply.refuse')}}</el-button>
                    <el-button type="primary" @click="video_accept">{{$t('apply.accept')}}</el-button>
            </div>
        </div>
        <!-- 主呼方收到同意回复 -->
        <div v-show="video_anwer" class="video_apply" >
                <div class="video_apply_title">
                    <!-- <span>{{your_callname}}</span> &nbsp;<span>{{$t('apply.request')}}</span>&nbsp;<span>{{your_calltype}}</span>&nbsp;<span>{{$t('apply.connection')}}</span> -->
                    <span>{{$t('apply.agree')}}</span>
                </div>
                <div class="video_apply_body">
                        <el-button type="info" @click="cancle_call">{{$t('button_message.cancel')}}</el-button>
                        <el-button type="primary" @click="video_initial">{{$t('apply.call')}}</el-button>
                </div>
        </div>
        <!-- 主呼收到拒绝回复 -->
        <div v-show="video_reject" class="video_apply" >
            <div class="video_apply_title">
                <!-- <span>{{your_callname}}</span> &nbsp;<span>{{$t('apply.request')}}</span>&nbsp;<span>{{your_calltype}}</span>&nbsp;<span>{{$t('apply.connection')}}</span> -->
                <span>{{$t('apply.reject')}}</span>
            </div>
            <div class="video_apply_body">
                    <el-button type="info" @click="cancle_reject">{{$t('button_message.cancel')}}</el-button>
                    <el-button type="primary" @click="confirm_reject">{{$t('button_message.confirm')}}</el-button>
            </div>
        </div>
            <!-- 主呼收到对方不在线或者不存在 -->
        <div v-if="video_offline" class="video_apply" >
                <div class="video_apply_title">
                    <span v-show="offline_span">{{$t('apply.offline')}}</span>
                    <span v-show="exist_span">{{$t('apply.exist')}}</span>
                </div>
                <div class="video_apply_body">
                        <el-button type="primary" @click="confirm_offline">{{$t('button_message.confirm')}}</el-button>
                </div>
        </div>
    </div>
    </template>
    
    <script>
    import adapter from 'webrtc-adapter'
    import $ from 'jquery'
    import Janus from '../../assets/videocall/janus.js'
    // import { win32 } from 'path';
    import intermap from '../inter-components/inter-baidumap.vue'
    import intergooglemap from '../inter-components/inter-googlemap'
    export default {
        components:{
                intermap,intergooglemap
            // testComponent:require('./testComponent.vue').default
            },
        data() {   
            return {
                    mapshow:true,
                    local_group_list:[],
                    get_device_list:[],
                    yesData: [],
                    select_Data:[],
                    show:true,
                    // online:6,
                    // totalnum:20,
                    active:null,
                    actived:null,
                    members_show:true,
                    down_show:true,
                    right_show:false,
                    computer_name:sessionStorage.getItem('loginName'),
                    group_div:false,
                    members_div:false, 
                    media_show:false,
                    editorial_show:false,
                    dialogVisible: false,
                    talkbackVisible:false,
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
                    Modify_group_info :{},
                    updata_group_info:{},
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
                    // 视频
                    video_show:false,
                    aduio_show:false,
                    video_name:'',
                    call_name:'',
                    videocall : null,
                    audio_answer:false,
                    video_hang:'',
                    audio_bridge:'',
                    audio_hang:'',
                    audio_logo:true,
                    // calling_Visible:false,
                    jesp_code:'',
                    // 即时聊天
                    im_show:false,
                    im_recent_show:false,
                    im_device_show:true,
                    im_group_show:false,
                    im_group_active:null,
                    im_group_actived:null,
                    im_device_active:null,
                    im_device_actived:null,
                    im_line_active:null,
                    im_line_actived:null,
                    talk_name:'',
                    talk_id:'',
                    im_select_text:true,
                    im_select_file:false,
                    im_select_talk:false,
                    history_right:'0px',
                    history_index:-1,
                    im_send_news:'',
                    im_send_obj:1,
                    receiver_type:1,
                    receiver_id:'',
                    im_now_date:'',
                    // im_my_img:require('../../assets/img/computer.png'),
                    // im_you_img:require( '../../assets/img/inter.png'),
                    pause_img:require('../../assets/img/voiceStatic.png'),
                    play_img:require('../../assets//img/voiceDynamic.gif'),
                    audio_img:this.pause_img,
                    local_mymessage : [],
                    local_mynews:[],
                    local_groupnews:[],
                    local_history:[],
                    im_content_show:true,
                    im_talkback_show:false,
                    network_message:{
                    },
                    local_message:[
                    ],
                    // websocket
                    websocket: null,
                    im_box_show:true,
                    file_name:'',
                    file_type:'',
                    upload_show:false,
                    other_show:false,
                    other_name:'',
                    sos_name:'',
                    other_send_type:'',
                    sos_show:false,
                    off_line_show:false,
                    online_show:false,
                    refresh_show:false,
                    off_line_name:'',
                    online_name:'',
                    // 单人、群
                    im_self_show:true,
                    im_record_show:false,
                    other_send_id:'',
                    other_group_id:'',
                    select_text :true,
                    select_file:false,
                    select_talk:false,
                        // 视频申请
                     video_apply:false,
                    your_callname:'',
                    your_calltype:'',
                    your_callid:'',
                    video_anwer:false,
                    video_audio:true,
                    is_Video:'',
                    call_id:'',
                    video_reject:false,
                    video_offline:false,
                    offline_span:false,
                    exist_span:false,
                    talkback_id:'',
                    talkactive:false,
                    mouse_show:true,
                    // jiadata:'',
    
                    // sosGps
                    // sosdata:'',
                    audio_file:'',
                    audio_play:true,
                    img_gif:false,
                    set_time:null,
                    time_long:'',
    
                    ssaudio:[
                        {id:1,file:'http://yss.yisell.com/yisell/pays2018050819052088/sound/yisell_sound_2008041415090410973_88011.mp3',name:'chat1',time:'2019年04月28日 18:07:46',long:12},
                        {id:2,file:'http://yss.yisell.com/yisell/pays2018050819052088/sound/yisell_sound_2008041415063070478_88011.mp3',name:'chat1',time:'2019年05月29日 18:07:46',long:10},
                        {id:3,file:'http://yss.yisell.com/yisell/pays2018050819052088/sound/yisell_sound_2008041415055629975_88011.mp3',name:'chat2',time:'2019年06月29日 18:07:46',long:5},
                        {id:4,file:'http://yss.yisell.com/yisell/pays2018050819052088/sound/yisell_sound_2008041415080752918_88011.mp3',name:'chat4',time:'2019年07月29日 18:07:46',long:8},
                        {id:5,file:'http://yss.yisell.com/yisell/pays2018050819052088/sound/yisell_sound_2008041415080752918_88011.mp3',name:'chat1',time:'2019年08月29日 18:07:46',long:8},
                    ]
        
            }
        },
        methods: {
                hide(){
                this.editorial_show = false.stop;
                this.media_show = false.stop
                },
                // 鼠标右击
                media_control(index){
                    this.media_show =index;
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
                // 即时聊天移入移出
                enter_im_group(index){
                    this.im_group_actived = index;     
                },
                leave_im_group(index){
                    index = null;
                    this.im_group_actived = index;
                },
                enter_im_device(index){
                    this.im_device_actived = index;     
                },
                leave_im_device(index){
                    index = null;
                    this.im_device_actived = index;
                },
                enter_im_outline(index){
                    this.im_line_actived = index;     
                },
                leave_im_outline(index){
                    index = null;
                    this.im_line_actived = index;
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
                    if(this.local_device_list == null){
                        return 
                    }else{
                        for( var i =0;i<this.local_device_list.length;i++){
                            for(var j=0;j<this.confirm_device.length;j++){
                                if(this.local_device_list[i].id == this.confirm_device[j]){
                                this.confirm_device_List.push(this.local_device_list[i])  
                                }
                            }
                        }
                    }
                    this.members_div = false;
                    this.group_memberdiv_show = true;
                },
                group_add_cancle(){
                    this.group_div = false;
                    this.group_memberdiv_show = false;
                    this.yesData =[];
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
                            group_info.account_id =  parseInt(sessionStorage.getItem('id')) ;
                            submit_form.group_info = group_info;
                            submit_form.device_ids = [-1];
                            submit_form.device_infos = this.confirm_device_List;
                            window.console.log(submit_form)
                            this.$axios.post('/group',submit_form,
                            { headers: 
                            {"Authorization" : sessionStorage.getItem('setSession_id')}
                            })
                            .then((response) =>{
                                window.console.log(response)
                            this.$message({
                            message: this.$t('establish.success'),
                            type: 'success'
                            });
                            this.get_new_group();
                            this.group_add_cancle();
                            window.console.log(response.data.group_info.gid)
                            var room_id= parseInt(response.data.group_info.gid);
                            var new_room={"request":"create","room":room_id,"permanent":true,"is_private":false}
                            this.mixertest.send({'message':new_room})
    
                            })
                            .catch( (error) => {
                                window.console.log(error)
                            // if( error.response.data.code == 422){
                            //         this.$message({                                   
                            //         message: this.$t('group.name'),
                            //         type: 'warning'
                            //         }); 
                            //     }else{
                            //         this.$message({                                   
                            //         message: this.$t('establish.failed'),
                            //         type: 'warning'
                            //         }); 
                            //     }
                            });
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
                    let dele_info = this.group_list[dele_mumber].group_info;
                    window.console.log(dele_info.id)
                    this.$axios.post('/group/delete',dele_info,
                    { headers: 
                    {"Authorization" : sessionStorage.getItem('setSession_id')}
                    })
                    .then(() =>{
                    this.$message({                                   
                    message: this.$t('group.dele_success'),
                    type: 'success'
                    }); 
                    this.dialogVisible = false;
                    this.get_new_group();
                    var dele_id= dele_info.id
                    var del_room={"request":"destroy","room":dele_id,"permanent":true,}
                    this.mixertest.send({'message':del_room});
                    this.group_div_show=false
                
                    })
                    .catch(function (error) {
                    window.console.log(error);
                    });
                },
                // 左侧修改组
                amend(){
                    this.amend_show = true;
                    let  amend_num = this.dele_num;
                    // let amend_group_list = JSON.parse(localStorage.getItem('group_list'));
                    // this.Modify_group_form.num = amend_group_list[amend_num].group_info.id
                    this.Modify_group_form.num = this.local_group_list[amend_num].group_info.id;
                    this.Modify_group_info = this.local_group_list[amend_num].group_info
    
                },
                // 左侧修改组提交
                amend_group_submit(){
                    if(this.Modify_group_form.name !== ''){
                    // window.console.log(this.Modify_group_info)
                    let  new_group_info =JSON.parse(JSON.stringify(this.Modify_group_info));
                    // new_group_info = this.Modify_group_info;
                    new_group_info.group_name = this.Modify_group_form.name;
                    window.console.log(this.Modify_group_info);
                    window.console.log(new_group_info);
                    this.$axios.post('/group/update',new_group_info,
                    { headers: 
                    {"Authorization" : sessionStorage.getItem('setSession_id')}
                    })
                    .then((response) =>{
                    window.console.log(response);
                    this.get_new_group();
                    this.$message({                                   
                    message: this.$t('group.modify_success'),
                    type: 'success'
                    }); 
                    })
                    .catch( (error)=>{
                        window.console.log(error.response.data.code);
                    if( error.response.data.code == 422){
                        this.$message({                                   
                        message: this.$t('group.name'),
                        type: 'warning'
                        }); 
                    }else{
                        this.$message({                                   
                        message: this.$t('group.modify_failed'),
                        type: 'warning'
                        }); 
                    }
                
                    
                    });            
                    this.amend_cancel();
                    }else{
                        this.$message({
                            message: this.$t('group.login_error'),
                            type: 'warning'
                        });
                    }
                },
                // 关闭修改组清空
                amend_cancel(){
                this.amend_show = false ;
                this.Modify_group_form.name = '';
                },
                // 左侧成员编辑
                modified_member(){
                    this.modified_member_show=true;
                    let  modified_num = this.dele_num;
                    window.console.log(modified_num);
                    let modified_member_data =[];
                    modified_member_data = this.local_group_list[modified_num].device_infos;
                    this.updata_group_info = this.local_group_list[modified_num].group_info;
                    window.console.log(this.local_group_list[modified_num]);
                    window.console.log(this.local_group_list[modified_num].group_info);
                    let selected_member_data=modified_member_data.map(e =>{
                        if(e.hasOwnProperty('id')){
                            return e.id
                        }
                    })
                    this.select_Data=selected_member_data;
                        window.console.log(selected_member_data)
                    window.console.log(this.select_Data)
                },
                modified_cancel(){
                    this.modified_member_show=false;
                },
                // 左侧成员编辑提交
                modified_add(){
                        window.console.log(this.select_Data);
                    this.modified_add_member = []      
                        for( var i=0;i<this.local_device_list.length;i++){
                            for(var j=0;j< this.select_Data.length;j++){
                                // if(modified_divice[i].id == this.select_Data[j]){
                                if(this.local_device_list[i].id == this.select_Data[j]){
                                    // this.modified_add_member.push(modified_divice[i])
                                    this.modified_add_member.push(this.local_device_list[i])
                                }
                            }
                        }
                        window.console.log(this.modified_add_member)
                        let group_info_new = {};
                        group_info_new.group_name = this.updata_group_info.group_name;
                        group_info_new.account_id = this.updata_group_info.account_id;
                        group_info_new.id = this.updata_group_info.id;
                        let sent_group = {};
                        sent_group.group_info = group_info_new;
                        sent_group.device_ids = [-1];
                        sent_group.device_infos = this.modified_add_member;
                        window.console.log(sent_group);
                        // sent_group.device_infos = [];
                        this.$axios.post('/group/devices/update',sent_group,
                            { headers: 
                            {"Authorization" : sessionStorage.getItem('setSession_id')}
                            })
                            .then(() =>{
                            this.$message({
                            message: this.$t('group.modify_success'),
                            type: 'success'
                            });
                            this.get_new_group();
                            this.group_add_cancle();
                            })
                            .catch( (error) => {
                                window.console.log(error)
                                // if( error.response.data.code == 422){
                                //     this.$message({                                   
                                //     message: this.$t('group.name'),
                                //     type: 'warning'
                                //     }); 
                                // }else{
                                //     this.$message({                                   
                                //     message: this.$t('establish.failed'),
                                //     type: 'warning'
                                //     }); 
                                // }
                        });
                        this.modified_member_show=false;
                        this.group_div_show=false;
                },
                transfer_cancel(){
                    this.confirm_device = this.yesData;
                    this.members_div = false;
                },
                // 获得组
                get_group_list(){
                    this.local_group_list = JSON.parse(localStorage.getItem('group_list'));
                    this.get_device_list = JSON.parse(localStorage.getItem('device_list'))
                },
                // 获取新群组
                get_new_group(){
                    this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
                    { headers: 
                    {
                    "Authorization" : sessionStorage.getItem('setSession_id')
                    }
                    })
                    .then((response) =>{
                        this.local_group_list = response.data.group_list;
                        localStorage.setItem('group_list', JSON.stringify(response.data.group_list));
                        // this.group_add_cancle()
                    })
                    .catch( (error) => {
                    window.console.log(error);
                    }); 
                },
                //视频通话
                video_server(){
                    var self =this;
                    // var server = [
                    // "ws://" + "113.105.153.240" + ":9188",
                    // "/janus"
                    // ];
                    // var server = [
                    // "ws://" + "ptt.jimilab.com" + ":9188",
                    // "/janus"
                    // // ];       
                    // var server = [
                    // "wss://" + "ptt.jimilab.com" + ":8989",
                    // "/janus"
                    // ];   
                    var server = [
                    "wss://" + "test.jimilab.com" + ":8989",
                    "/janus"
                    ];                       
                    // var server = "https://" + "ptt.jimilab.com" + ":8089/janus"
                    var janus = null;
                    var opaqueId = "videocalltest-"+Janus.randomString(12);
                    var myusername = null;
                    var yourusername = null;
                    var incoming  =null;
                    if(!Janus.isWebrtcSupported()) {
                                alert("No WebRTC support... ");
                            return;
                    }
                    Janus.init({debug: true, callback: function() {
                    janus = new Janus({
                        server: server,
                        success: function() {
                        janus.attach({
                            plugin: "janus.plugin.videocall",
                            opaqueId: opaqueId,
                            success: function(pluginHandle) {
                                self.videocall = pluginHandle;
                                Janus.log("Plugin attached! (" + self.videocall.getPlugin() + ", id=" + self.videocall.getId() + ")");
                                // $('#videocall').removeClass('hide').show();
                                var user_register =sessionStorage.getItem('id').toString();
                                var viedo_register = { "request": "register", "username": user_register };
                                self.videocall.send({"message": viedo_register});   
                                window.console.log('````````````````````````1')
                            },
                            error: function(error) {
                                Janus.error("  -- Error attaching plugin...", error);
                                alert("  -- Error attaching plugin... " + error);
                                window.console.log('````````````````````````2')
                            },
                            consentDialog: function(on) {
                                Janus.debug("Consent dialog should be " + (on ? "on" : "off") + " now");
                                window.console.log('````````````````````````3')
                            },
                            mediaState: function(medium, on) {
                                Janus.log("Janus " + (on ? "started" : "stopped") + " receiving our " + medium);
                                window.console.log('````````````````````````4')
                            },
                            webrtcState: function(on) {
                                Janus.log("Janus says our WebRTC PeerConnection is " + (on ? "up" : "down") + " now");
                                window.console.log(on)
                                if(on == true&&self.video_audio == false){
                                    window.console.log('............0')
                                    setTimeout(()=>{self.close_video()},5000)
                                    // self.videocall.send({"message": { "request": "set", "video": false}});  
                                }
                                window.console.log('````````````````````````5')
                            },
                            onmessage: function(msg, jsep) {
                                window.console.log('````````````````````````6')
                                Janus.debug(" ::: Got a message :::");
                                Janus.debug(msg);
                                window.console.log(msg)
                                var result = msg["result"];
                                if(result !== null && result !== undefined) {
                                    if(result["list"] !== undefined && result["list"] !== null) {
                                        window.console.log('````````````````````````7')
                                        var list = result["list"];
                                        Janus.debug("Got a list of registered peers:");
                                        Janus.debug(list);
                                        for(var mp in list) {
                                        Janus.debug("  >> [" + list[mp] + "]");
                                        }
                                    } else if(result["event"] !== undefined && result["event"] !== null) {
                                    var event = result["event"];
                                    if(event === 'registered') {
                                        window.console.log('````````````````````````8')
                                        myusername = result["username"];
                                        self.videocall.send({"message": { "request": "list" }});
                                    } else if(event === 'calling') {
                                        Janus.log("Waiting for the peer to answer...");
                                        // alert("Waiting for the peer to answer...");
                                        window.console.log('````````````````````````9')
                                    } else if(event === 'incomingcall') {
                                        window.console.log('````````````````````````10')
                                        Janus.log("Incoming call from " + result["username"] + "!");
                                        yourusername = result["username"];
                                        self.$confirm("Incoming call from " + result["username"] + "!" , {
                                        confirmButtonText: self.$t("button_message.confirm"),
                                        cancelButtonText: self.$t("button_message.cancel"),
                                        type: 'warning'
                                        }).then(() => {
                                        // window.console.log(self.jesp_code)
                                        incoming=null;
                                        self.videocall.createAnswer(
                                            {
                                                jsep: jsep,
                                                // No media provided: by default, it's sendrecv for audio and video
                                                media: { data: true },	
                                                // Let's negotiate data channels as well
                                                // If you want to test simulcasting (Chrome and Firefox only), then
                                                // pass a ?simulcast=true when opening this demo page: it will turn
                                                // the following 'simulcast' property to pass to janus.js to true
                                                // simulcast: true,
                                                success: function(jsep) {
                                                    Janus.debug("Got SDP!");
                                                    Janus.debug(jsep); 
                                                    window.console.log(self.video_audio)
                                                            if(self.video_audio == true){
                                                                self.video_show= true;
                                                            }else{
                                                                self.aduio_show= true;
                                                                self.audio_logo=false  
                                                    }
                                                    var accept_body = { "request": "accept"};
                                                    self.videocall.send({"message": accept_body, "jsep": jsep});
                                                    $('#videos').show();
                                                },
                                                error: function(error) {
                                                    Janus.error("WebRTC error:", error);
                                                    // window.console.log(jsep)
                                                    alert("WebRTC error... " + JSON.stringify(error));
                                                }
                                            });	  
                                        }).catch(() => {
                                            
                                        });
                                    } else if(event === 'accepted') {
                                        window.console.log('````````````````````````11')
                                        if(jsep){
                                            self.videocall.handleRemoteJsep({jsep: jsep});
                                        }
                                        
                                        if(self.video_audio == true){
                                            window.console.log('````````````````````````12')   
                                                    //    self.close_video()      
                                            }else{
                                                window.console.log('````````````````````````13')
                                            
                                                //  self.videocall.send({"message": { "request": "set", "video": false}});                    
                                                //  self.videocall.send({"message": { "request": "set", "audio": false}});  
                                                // setTimeout(()=>{self.close_video()},5000)
                                                // self.setTimeout(self.close_video(),5000);                  
                                        }
                                  
                                    } else if(event === 'update') {
                                        window.console.log('````````````````````````13')
                                        if(jsep) {
                                        if(jsep.type === "answer") {
                                            window.console.log('````````````````````````14')
                                            self.videocall.handleRemoteJsep({jsep: jsep});
                                        } else {
                                            window.console.log('````````````````````````15')
                                            self.videocall.createAnswer(
                                            {
                                                jsep: jsep,
                                                media: { data: true },	// Let's negotiate data channels as well
                                                success: function(jsep) {
                                                Janus.debug("Got SDP!");
                                                Janus.debug(jsep);
                                                var body = { "request": "set",
                                                };
                                                self.videocall.send({"message": body, "jsep": jsep});
                                                },
                                                error: function(error) {
                                                Janus.error("WebRTC error:", error);
                                                // alert("WebRTC error... " + JSON.stringify(error));
                                                }
                                            });
                                        }
                                        }
                                    } else if(event === 'hangup') {
                                        window.console.log('````````````````````````16')
                                        Janus.log("Call hung up by " + result["username"] + " (" + result["reason"] + ")!");
                                        // Reset status
                                        if(self.video_audio == true){
                                                self.video_show= false;
                                            }else{
                                                self.aduio_show= false;
                                                self.audio_logo=false  
                                    }
                                        self.videocall.hangup();
                                        $('#videos').hide();
    
                                    }
                                    }
                                } else {
                                    window.console.log('````````````````````````17')
                                    var no_video =msg
                                        if(no_video.event === 'notified'){
                                            window.console.log('5--------------------------------------')
                                        }else{
                                            // FIXME Error?
                                            window.console.log('6--------------------------------------')
                                            
                                    var error = msg["error"];
                                    alert(error);
                                    self.$router.go(0)
                                    // self.video_hang.destroy();
                                    self.video_show=false;
                                    self.videocall.hangup();
                                    $('#videos').hide();
                                     }
                                }
                            },
                            onlocalstream: function(stream) {
                                window.console.log('````````````````````````18')
                                Janus.debug(" ::: Got a local stream :::");
                                Janus.debug(stream);
                                // self.video_show= true
                                var videoTracks =''
                                if(self.audio_answer == true){
                                    $('#audios').removeClass('hide').show();
                                    if($('#myaudio').length === 0)
                                    $('#audiocall').append('<video class="rounded centered" id="myaudio"  style="display:none" width=131 height=83 autoplay playsinline muted="muted"  />');
                                    Janus.attachMediaStream($('#myaudio').get(0), stream);
                                    $("#myaudio").get(0).muted = "muted";
                                    // if(audio_self.videocall.webrtcStuff.pc.iceConnectionState !== "completed" &&
                                    //     audio_self.videocall.webrtcStuff.pc.iceConnectionState !== "connected") {
                                    //     $('#audioright').append('<video class="rounded centered" id="waitingvideo" width=406 height=271    />');
                                    // }
                                    videoTracks = stream.getVideoTracks();
                                    if(videoTracks === null || videoTracks === undefined || videoTracks.length === 0) {
                                        // No webcam
                                        $('#myaudio').hide();
                                        if($('#audiocall .no-video-container').length === 0) {
                                        $('#audiocall').append(
                                            '<div class="no-video-container">' +
                                            '<span class="no-video-text">No webcam available</span>' +
                                            '</div>');
                                        }
                                    } else {
                                        $('#audiocall .no-video-container').remove();
                                        // $('#myvideo').removeClass('hide').show();
                                    }
                                }else{
                                    window.console.log('````````````````````````19')
                                    $('#videos').show();
                                    if($('#myvideo').length === 0)
                                        $('#videoleft').append('<video class="rounded centered" id="myvideo" width=131 height=83 autoplay playsinline muted="muted" />');
                                    Janus.attachMediaStream($('#myvideo').get(0), stream);
                                    $("#myvideo").get(0).muted = "muted";
                                    videoTracks = stream.getVideoTracks();
                                    if(videoTracks === null || videoTracks === undefined || videoTracks.length === 0) {
                                        // No webcam
                                        $('#myvideo').hide();
                                        if($('#videoleft .no-video-container').length === 0) {
                                        $('#videoleft').append(
                                            '<div class="no-video-container">' +
                                            '<span class="no-video-text">No webcam available</span>' +
                                            '</div>');
                                        }
                                    } else {
                                        $('#videoleft .no-video-container').remove();
                                        $('#myvideo').removeClass('hide').show();
                                    }
                                }
                            },
                            onremotestream: function(stream) {
                                window.console.log('````````````````````````20')
                                Janus.debug(" ::: Got a remote stream :::");
                                Janus.debug(stream);
                                var addButtons = false;
                                var videoTracks =''
                                if(self.audio_answer == true){
                                    if($('#remoteaudio').length === 0) {
                                    addButtons = true;
                                    $('#audioright').append('<video class="rounded centered hide" id="remoteaudio" style="display:none"  width=406 height=271   autoplay playsinline />');
                                    // $('.audio_loding').hide()
                                    self.audio_logo=false                 
                                    // $('#audio_box').show();
                                    $("#remoteaudio").bind("playing", function () {
                                        if(this.videoWidth)
                                        // $('#remotevideo').removeClass('hide').show();
                                        var width = this.videoWidth;
                                        var height = this.videoHeight;
                                    });
                                    }
                                    Janus.attachMediaStream($('#remoteaudio').get(0), stream);
                                    videoTracks = stream.getVideoTracks();
                                    if(videoTracks === null || videoTracks === undefined || videoTracks.length === 0) {
                                    // No remote video
                                    $('#remoteaudio').hide();
                                    if($('#audioright .no-video-container').length === 0) {
                                        $('#audioright').append(
                                        '<div class="no-video-container">' +
                                            '<span class="no-video-text">No remote video available</span>' +
                                        '</div>');
                                    }
                                    } else {
                                    $('#audioright .no-video-container').remove();
                                    }
    
                                }else{
                                    window.console.log('````````````````````````21')
                                if($('#remotevideo').length === 0) {
                                    addButtons = true;
                                    $('#videoright').append('<video class="rounded centered hide" id="remotevideo" width=406 height=271   autoplay playsinline/>');
                                    $("#remotevideo").bind("playing", function () {
                                        if(this.videoWidth)
                                        $('#remotevideo').removeClass('hide').show();
                                        var width = this.videoWidth;
                                        var height = this.videoHeight;
                                    });
                                    }
                                    Janus.attachMediaStream($('#remotevideo').get(0), stream);
                                    videoTracks = stream.getVideoTracks();
                                    if(videoTracks === null || videoTracks === undefined || videoTracks.length === 0) {
                                    // No remote video
                                    $('#remotevideo').hide();
                                    if($('#videoright .no-video-container').length === 0) {
                                        $('#videoright').append(
                                        '<div class="no-video-container">' +
                                            '<span class="no-video-text">No remote video available</span>' +
                                        '</div>');
                                    }
                                    } else {
                                    $('#videoright .no-video-container').remove();
                                    $('#remotevideo').removeClass('hide').show();
                                    }
                                }
                            },
                            });
                        },
                        error: function(error) {
                            window.console.log('````````````````````````22')
                            Janus.error(error);
                            alert(error, function() {
                            window.location.reload();
                            });
                            self.$router.go(0)
                        },
                        mycallmessage:function(e){
                            window.console.log('````````````````````````23')
                            window.console.log(e);
                            window.console.log('--------------------p');
                            if(e.code == 1){
                                self.video_offline = true;
                                self.offline_span=true
                                
                            }else if(e.code ==2){
                                self.video_offline = true;
                                self.exist_span=true
    
                            }
                        },
                        youcallmessage:function(e){
                            window.console.log('````````````````````````24')
                            window.console.log(e);
                            window.console.log('````````````````````````````p');
                            if(e.type == 0&&e.isAccept == 0){
                                self.video_apply =true;
                                self.your_callname = e.name;
                                self.your_callid = e.myId;
                                if(e.isVideo == "true"){
                                    self.your_calltype = 'video';
                                    self.is_Video= 'true';
                                    self.video_audio=true;
                                    window.console.log('````````````````````````0')
                                }else{
                                    self.your_calltype = 'audio';
                                    self.is_Video= 'false';
                                    self.video_audio=false
                                    window.console.log('````````````````````````5')
                                }
                            }
                            if(e.type == 1){
                                if(e.isAccept == 2){
                                  self.video_anwer =true;
                                  self.call_id =e.myId
                                   if(e.isVideo == 'true'){
                                    self.video_audio=true;
                                   }else{
                                    self.video_audio;
                                    self.video_audio=false
                                   }
                                } else if(e.isAccept == 1){
                                    self.video_reject =true
                                }
                            }
                            
                        },
                        destroyed: function() {
                            // window.location.reload();
                            window.console.log('````````````````````````25')
                        }
                        });
                    }})
                    this.video_hang=janus;
                },
                video_begin(item){
                    this.audio_answer = false;
                    var body ={
                    'uid':item.id,'type':0,"name":sessionStorage.getItem('username'),"isVideo":'true','isAccept':0,
                     'myId':sessionStorage.getItem('id')
                    }
                    window.console.log(body)
                    this.videocall.send({'user_call':body})
                    
            
                },
                // 语音通话
                audio_begin(item){
                    // var audio_self=this 
                    var body ={
                    'uid':item.id,'type':0,"name":sessionStorage.getItem('username'),"isVideo":'false','isAccept':0,
                     'myId':sessionStorage.getItem('id')
                    }
                    window.console.log(body)
                    // audio_self.videocall.send({'user_call':body})
                    this.videocall.send({'user_call':body})
                },
                 //    对讲服务拉起
                audiobridge_serve(){
                    var bridge =this;
                    var server = [
                    "wss://" + "test.jimilab.com" + ":8989",
                    "/janus"
                    ];   
                    // var server = [
                    // "wss://" + "ptt.jimilab.com" + ":8989",
                    // "/janus"
                    // ];  
                    var janus = null;
                    // var mixertest = null;
                    var opaqueId = "audiobridgetest-"+Janus.randomString(12);
                    var spinner = null;
                    var myroom = '';	// Demo room
                    var myusername = null;
                    var myid = null;
                    var webrtcUp = false;
                    var audioenabled = false;
                    Janus.init({debug: "all", callback: function() {
                        // Use a button to start the demo
                        // Make sure the browser supports WebRTC
                        if(!Janus.isWebrtcSupported()) {
                            alert("No WebRTC support... ");
                            return;
                        }
                        // Create session
                        janus = new Janus(
                            {
                                server: server,
                                success: function() {
                                    // Attach to Audio Bridge test plugin
                                    janus.attach(
                                        {
                                            plugin: "janus.plugin.audiobridge",
                                            opaqueId: opaqueId,
                                            success: function(pluginHandle) {
                                                bridge.mixertest = pluginHandle;
                                                Janus.log("Plugin attached! (" + bridge.mixertest.getPlugin() + ", id=" + bridge.mixertest.getId() + ")");
                                            },
                                            error: function(error) {
                                                Janus.error("  -- Error attaching plugin...", error);
                                            alert("Error attaching plugin... " + error);
                                            },
                                            consentDialog: function(on) {
                                                Janus.debug("Consent dialog should be " + (on ? "on" : "off") + " now");
                                            },
                                            onmessage: function(msg, jsep) {
                                                Janus.debug(" ::: Got a message :::");
                                                Janus.debug(msg);
                                                var event = msg["audiobridge"];
                                                Janus.debug("Event: " + event);
                                                if(event != undefined && event != null) {
                                                    if(event === "joined") {
                                                        // Successfully joined, negotiate WebRTC now
                                                        myid = msg["id"];
                                                        Janus.log("Successfully joined room " + msg["room"] + " with ID " + myid);
                                                        if(!webrtcUp) {
                                                            webrtcUp = true;
                                                            // Publish our stream
                                                            bridge.mixertest.createOffer(
                                                                {
                                                                    media: { video: false},	// This is an audio only room
                                                                    success: function(jsep) {
                                                                        Janus.debug("Got SDP!");
                                                                        Janus.debug(jsep);
                                                                        var publish = { "request": "configure", "muted": false };
                                                                        bridge.mixertest.send({"message": publish, "jsep": jsep});
                                                                    },
                                                                    error: function(error) {
                                                                        Janus.error("WebRTC error:", error);
                                                                    alert("WebRTC error... " + JSON.stringify(error));
                                                                    }
                                                                });
                                                        }
                                                        // Any room participant?
                                                        if(msg["participants"] !== undefined && msg["participants"] !== null) {
                                                            var lists = msg["participants"];
                                                            Janus.debug("Got a list of participants:");
                                                            Janus.debug(lists);
                                                        }
                                                    } else if(event === "roomchanged") {
                                                        // The user switched to a different room
                                                        myid = msg["id"];
                                                        Janus.log("Moved to room " + msg["room"] + ", new ID: " + myid);
                                                        // Any room participant?
                                                        if(msg["participants"] !== undefined && msg["participants"] !== null) {
                                                            var listss = msg["participants"];
                                                            Janus.debug("Got a list of participants:");
                                                            Janus.debug(listss);
                                                        }
                                                    } else if(event === "destroyed") {
                                                        // The room has been destroyed
                                                        Janus.warn("The room has been destroyed!");
                                                        alert("The room has been destroyed", function() {
                                                            window.location.reload();
                                                        });
                                                    } else if(event === "event") {
                                                        if(msg["participants"] !== undefined && msg["participants"] !== null) {
                                                            var listsss = msg["participants"];
                                                            Janus.debug("Got a list of participants:");
                                                            Janus.debug(listsss);
                                                        } else if(msg["error"] !== undefined && msg["error"] !== null) {
                                                            if(msg["error_code"] === 485) {
                                                                // This is a "no such room" error: give a more meaningful description
                                                            alert(
                                                                    "<p>Apparently room <code>" + myroom + "</code> (the one this demo uses as a test room) " +
                                                                    "does not exist...</p><p>Do you have an updated <code>janus.plugin.audiobridge.cfg</code> " +
                                                                    "configuration file? If not, make sure you copy the details of room <code>" + myroom + "</code> " +
                                                                    "from that sample in your current configuration file, then restart Janus and try again."
                                                                );
                                                            } else {
                                                                alert(msg["error"]);
                                                            }
                                                            return;
                                                        }
                                                        // Any new feed to attach to?
                                                    }
                                            }
                                                if(jsep !== undefined && jsep !== null) {
                                                    Janus.debug("Handling SDP as well...");
                                                    Janus.debug(jsep);
                                                    bridge.mixertest.handleRemoteJsep({jsep: jsep});
                                                }
                                            },
                                            onlocalstream: function(stream) {
                                                Janus.debug(" ::: Got a local stream :::");
                                                Janus.debug(stream);
                                            },
                                            onremotestream: function(stream) {
                                        
                                                // Mute button
                                                audioenabled = true;
                                            },
                                            oncleanup: function() {
                                                webrtcUp = false;
                                                Janus.log(" ::: Got a cleanup notification :::");
                                            }
                                        });
                                },
                                error: function(error) {
                                    Janus.error(error);
                                    alert(error, function() {
                                        window.location.reload();
                                    });
                                    bridge.$router.go(0)

                                },
                                destroyed: function() {
                                    window.location.reload();
                                }
                            }); 
}});
                this.audio_bridge=janus
            },

                // 即时链接
                text_begin(item){
                    window.console.log(item.id);
                    window.console.log(sessionStorage.getItem('id'));
                    this.group_div_show = false;
                    this.im_show = true;
                    this.receiver_id = item.id;
                    this.talk_id = this.receiver_id
                    window.console.log(this.receiver_id)
                    this.talk_name = item.user_name;  
                    // this.websocket = new WebSocket('ws://10.0.18.132:10000/im-server/'+sessionStorage.getItem('id'));
                    // this.websocket = new WebSocket('ws://113.105.153.240:8888/im-server/'+sessionStorage.getItem('id'));
                    // this. initWebSocket();
                    this.im_self_show =true;
                    this.receiver_type = 1;
                    var text_local_num= this.receiver_type+'.'+sessionStorage.getItem('id')+'_'+this.talk_id
                    window.console.log(text_local_num)
                    this.local_mynews =JSON.parse(localStorage.getItem(text_local_num))  
                }, 
                createWebSocket(){
                    // this.websocket = new WebSocket('ws://10.0.18.132:10000/im-server/'+sessionStorage.getItem('id'));
                    // this.websocket = new WebSocket('ws://10.0.18.132:8888/im-server/'+sessionStorage.getItem('id'));
                    // this.websocket = new WebSocket('ws://113.105.153.240:8888/im-server/'+sessionStorage.getItem('id'));
                    // this.websocket = new WebSocket('wss://ptt.jimilab.com/websocket/'+sessionStorage.getItem('id'));
                    // this.websocket = new WebSocket('wss://ptt.jimilab.com:8888/im-server/'+sessionStorage.getItem('id'));
                    this.websocket = new WebSocket('wss://test.jimilab.com:10000/im-server/'+sessionStorage.getItem('id'));
                    this. initWebSocket();
                },
                initWebSocket() {
                    //连接错误
                    this.websocket.onerror = this.setErrorMessage
                
                    // //连接成功
                    this.websocket.onopen = this.setOnopenMessage
                
                    //收到消息的回调
                    this.websocket.onmessage = this.setOnmessageMessage
                
                    //连接关闭的回调
                    this.websocket.onclose = this.setOncloseMessage
                
                    //监听窗口关闭事件，当窗口关闭时，主动去关闭websocket连接，防止连接还没断开就关闭窗口，server端会抛异常。
                    window.onbeforeunload = this.onbeforeunload
                },
                // 获取消息
                //挂断
                video_hangup(){
                    var hangup = { "request": "hangup" };
                    window.console.log(this.videocall)
                    this.videocall.send({"message": hangup});
                    this.videocall.hangup();
                    this.yourusername = null;
                    // this.video_hang.destroy();
                    this.video_show=false;
                    $('#remotevideo').hide()
                    // this.aduio_show=false;
                },
                audio_hangup(){
                    var hangup = { "request": "hangup" };
                    window.console.log(this.videocall)
                    this.videocall.send({"message": hangup});
                    this.videocall.hangup();
                    this.yourusername = null;
                    // this.video_hang.destroy();
                    this.aduio_show=false;
                    $('#remotevideo').hide();
                    this.audio_logo=true;        
                },
                refuse(){
                    this.video_apply=false;
                        var body ={
                        'uid':this.your_callid,'type':1,"name":sessionStorage.getItem('username'),"isVideo":this.is_Video,'isAccept':1,
                        'myId':sessionStorage.getItem('id')
                        }
                        // var body ={
                        // 'uid':parseInt(this.your_callid),'type':1,"name":'text001',"isVideo":'true','isAccept':1,
                        //  'myId':1501
                        // }
                        window.console.log(body)
                        // this.videocall.send({'user_call':body})
                },
                video_accept(){
                        this.video_apply=false;
                        var body ={
                        'uid':this.your_callid,'type':1,"name":sessionStorage.getItem('username'),"isVideo":this.is_Video,'isAccept':2,
                        'myId':sessionStorage.getItem('id')
                        }
                        // var body ={
                        // 'uid':parseInt(this.your_callid),'type':1,"name":'text001',"isVideo":'true','isAccept':2,
                        //  'myId':1501
                        // }
                        window.console.log(body)
                        this.videocall.send({'user_call':body})
                },
                // 取消通话
                cancle_call(){
                      this.video_anwer=false;
                },
                // 拒绝取消
                cancle_reject(){
                    this.video_reject =false
                },
                confirm_reject(){
                    this.video_reject =false
                },
                confirm_offline(){
                    this.video_offline =false;
                    this.exist_span =false;
                    this.offline_span =false;
                },
                // 发起通话
                video_initial(){
                    var selfs =this;
                    this.video_anwer =false
                    if(this.video_audio == true){
                        this.aduio_show = false;
                        this.video_show=true;
                        this.audio_answer = false;
                        var userpeer = this.call_id.toString();
                        window.console.log(userpeer)
                            this.videocall.createOffer(
                            {
                                media: { data: true },	
                                success: function(jsep) {
                                var body = { "request": "call", "username": userpeer };
                                selfs.videocall.send({"message": body, "jsep": jsep});
                                },
                                error: function(error) {
                                Janus.error("WebRTC error...", error);
                                }
                            });	 
                    }else{
                        this.video_show=false;
                        this.aduio_show=true; 
                        this.audio_answer = true
                        var audiopeer = this.call_id.toString();
                            this.videocall.createOffer(
                            {
                                media: { data: true },	
                                success: function(jsep) {
                                var body = { "request": "call", "username": audiopeer };
                                selfs.videocall.send({"message": body, "jsep": jsep});
                                // selfs.videocall.send({"message":{"request":"set","video":false}})
                                // selfs.videocall.send({"message": { "request": "set", "video": false}});
                                // selfs.videocall.send({"message":{"request":"set","audio":false}})
                                },
                                error: function(error) {
                                Janus.error("WebRTC error...", error);
                                }
                            });	 
                    }
                },
                // 挂断视频
                close_video(){
                    this.videocall.send({"message": { "request": "set", "video": false}});
                    window.console.log('````````````````````````````````````````````````1')
                },
                // 视频移动
                video_move(e){
                    let odiv = e.target;    //获取目标元素
                    //算出鼠标相对元素的位置
                    let disX = e.clientX - odiv.offsetLeft;
                    let disY = e.clientY - odiv.offsetTop;
                    document.onmousemove = (e)=>{    //鼠标按下并移动的事件
                        //用鼠标的位置减去鼠标相对元素的位置，得到元素的位置
                        let left = e.clientX - disX;  
                        let top = e.clientY - disY;
                        
                        //绑定元素位置到positionX和positionY上面
                        // this.positionX = top;
                        // this.positionY = left;
                        
                        //移动当前元素
                        odiv.style.left = left + 'px';
                        odiv.style.top = top + 'px';
                    };
                    document.onmouseup = (e) => {
                        document.onmousemove = null;
                        document.onmouseup = null;
                    };
                },
                // 即时聊天
                select_imrecent(){
                    this.im_recent_show = true;
                    this.im_device_show = false;
                    this.im_group_show = false;
                    this.im_line_active = 'a' ;
                    this.talk_name ='';
                    this.talk_id ='';
                },
                select_imdevice(){
                    this.im_device_active = 'a' ;
                    this.im_recent_show = false;
                    this.im_device_show = true;
                    this.im_group_show = false;
                    this.im_record_show = false;
                    this.im_self_show = true;
                    this.local_mynews =[];
                    this.receiver_type = 1;
                    this.talk_name ='';
                    this.talk_id ='';
                    if(this.history_index == 1){
                        this.history_index = -1;
                        this.history_right='0px'
                   }
                },
                select_imgroup(){
                    this.im_group_active = 'a' ;
                    this.im_recent_show = false;
                    this.im_device_show = false;
                    this.im_group_show = true;
                    this.im_self_show = false;
                    this.im_record_show = true;
                    this.local_groupnews = [];
                    this.receiver_type = 2;
                    this.talk_name ='';
                    this.talk_id ='';
                    if(this.history_index == 1){
                        this.history_index = -1;
                        this.history_right='0px'
                   }
                },
                // 即时聊天点击变色事件
                select_im_group(index,select_id){
                    this.im_group_active = index ;
                    this.im_content_show =false;
                    this.im_content_show =true;
                    this.im_self_show =false
                    this.im_record_show =true
                    window.console.log(index);
                    window.console.log(select_id);
                    window.console.log(this.group_list);
                    window.console.log(this.group_list[index].group_info.id);
                    this.talkback_id=select_id;
                    this.talk_name = this.group_list[index].group_info.group_name;
                    this.talk_id = this.group_list[index].group_info.id;
                    this.receiver_id = this.talk_id
                    this.receiver_type = 2;
                    var im_group_num= this.receiver_type+'.'+sessionStorage.getItem('id')+'_'+this.talk_id
                    window.console.log(im_group_num);
                    window.console.log(JSON.parse(localStorage.getItem(im_group_num)));
                    this.local_groupnews =JSON.parse(localStorage.getItem(im_group_num));
                    if(this.history_index == 1){
                        this.history_index = -1;
                        this.history_right='0px'
                   }  
                },
                select_im_device(index,select_id){
                    this.im_device_active = index ;
                    this.im_content_show =false;
                    this.im_content_show =true;
                    this.im_record_show =false;
                    this.im_self_show =true;
                    window.console.log(index);
                    window.console.log(select_id);
                    // this.local_mymessage =[]
                    // this.local_mynews =[]
                    window.console.log(this.local_device_list[index].id);
                    this.talk_name = this.local_device_list[index].user_name;
                    this.talk_id = this.local_device_list[index].id;
                    this.receiver_id = this.talk_id
                    this.receiver_type = 1;
                    var im_device_num= this.receiver_type+'.'+sessionStorage.getItem('id')+'_'+this.talk_id
                    window.console.log(im_device_num)
                    this.local_mynews =JSON.parse(localStorage.getItem(im_device_num));
                    if(this.history_index == 1){
                        this.history_index = -1;
                        this.history_right='0px'
                   }  
    
                },
                select_im_outline(index,id,name,type,data,local){
                    window.console.log(index)
                    window.console.log(id)
                    window.console.log(name)
                    window.console.log(type)
                    window.console.log(data)
                    window.console.log(local)
    
                    this.im_line_active = index ;
                    this.im_content_show =false;
                    this.im_content_show =true;
                    this.im_record_show =false;
                    this.im_self_show =true;
                    this.talk_name=name
                    
                    if(type==1){
                        this.id=local.SenderId
                        this.im_record_show = false;
                        this.im_self_show = true;
                        var your_out_name=type+'.'+sessionStorage.getItem('id')+'_'+local.SenderId;
                        window.console.log(your_out_name);
                        window.console.log(JSON.parse(localStorage.getItem(your_out_name)));
                        this.local_mynews = JSON.parse(localStorage.getItem(your_out_name));
                                if(this.local_mynews  ==null){
                                    this.local_mynews = [] ;
                                    for(var v=0;v<data.length;v++){
                                        this.local_mynews.push(data[v]);
                                    }
                                }else{
                                    for(var s=0;s<data.length;s++){
                                        this.local_mynews.push(data[s]);
                                    }
                                   
                                }
                        // localStorage.setItem(your_out_name, JSON.stringify(this.local_mynews));  
    
                    }else{
                        this.id=local.GroupId
                        this.im_self_show = false;
                        this.im_record_show = true;
                        var group_out_name=type+'.'+sessionStorage.getItem('id')+'_'+local.GroupId;
                        window.console.log(group_out_name);
                        window.console.log(JSON.parse(localStorage.getItem(group_out_name)));
                        this.local_groupnews = JSON.parse(localStorage.getItem(group_out_name));
                                        if(this.local_groupnews  ==null){
                                            this.local_groupnews = [] ;
                                            for(var u=0;u<data.length;u++){
                                                this.local_groupnews.push(data[u]);
                                            }
                                            window.console.log( this.local_groupnews)
                                            
                                        }else{
                                            for(var q=0;q<data.length;q++){
                                                this.local_groupnews.push(data[q]);
                                                
                                            }
                                            window.console.log( this.local_groupnews)
                                         
                                        }
                    // localStorage.setItem(group_out_name, JSON.stringify(this.local_groupnews));  
                    }
                    
                },
                // 历史记录
                im_unfold_history(){
                   if(this.history_index == -1){
                        this.history_index=1;
                        this.history_right='-332px'
                   }else{
                        this.history_index = -1;
                        this.history_right='0px'
                   }
                
                   window.console.log(this.receiver_type)
                   if(this.receiver_type ==1){
                    
                    window.console.log(this.talk_id)
                    var im_alone_history = this.receiver_type +'.'+sessionStorage.getItem('id')+'_'+this.talk_id;
                    window.console.log(im_alone_history);
                    this.local_history = JSON.parse(localStorage.getItem(im_alone_history));
                    window.console.log( this.local_history);
                    window.console.log( this.im_total_history);
                   }else{
                    window.console.log(this.talk_id)
                    var im_group_history = this.receiver_type +'.'+sessionStorage.getItem('id')+'_'+this.talk_id;
                    window.console.log(im_group_history)
                    this.local_history = JSON.parse(localStorage.getItem(im_group_history));
                    window.console.log( this.local_history);
                    window.console.log( this.im_total_history);
                   }
    
                //    window.console.log(this.local_history)
                //    window.console.log(this.im_total_history)
                },
                im_collapse_history(){
                 this.history_index=-1;
                 this.history_right='0px'
                },
                select_im_text(){
                    this.im_select_text = true;
                    this.im_select_file =false;
                    this.im_select_talk =false;
                    this.select_text =true;
                    this.select_file=false;
                    this.select_talk=false;
                },
                select_im_file(){
                    this.im_select_text = false;
                    this.im_select_file =true;
                    this.im_select_talk =false;
                    this.select_file =true;
                    this.select_text =false;
                    this.select_talk=false;
                },
                select_im_talk(){
                    this.im_select_text = false;
                    this.im_select_file =false;
                    this.im_select_talk =true;
                    this.select_file =false;
                    this.select_text =false;
                    this.select_talk=true;
    
                },
                play_audio(item,index){
                     if(this.set_time != null){
                        clearTimeout(this.set_time)
                    }
                    if(this.audio_file.paused ==true ||this.audio_file.paused ==undefined){
                        this.img_gif =index
                        var time_this =this       
                        var ad=item.long
                        window.console.log(ad)
                        this.set_time=setTimeout(()=>{time_this.close_gif()},1000*ad)
                        this.audio_file = new Audio()
                        this.audio_file.src=item.file
                        // this.audio_file.load() 
                        this.audio_file.play()
                         window.console.log("`````2");
                    }else{
                        window.console.log("`````1")
                    }
                },
                close_gif(){ 
                    this.img_gif =false
                },
                // 发送消息
                im_send_subimt(){
                    if( this.im_box_show ==true){
                    this.im_create_time();
                    let message_body ={}
                    this.im_send_obj=1;
                    message_body.id = sessionStorage.getItem('id');
                    message_body.ReceiverType = this.receiver_type;
                    message_body.ReceiverId = this.receiver_id;
                    message_body.im_send_obj = this.im_send_obj;
                    message_body.time = this.im_now_date;
                    message_body.time_id = (new Date()).getTime();
                    message_body.MsgType = 1;
                    message_body.ResourcePath = this.im_send_news;
                    message_body.ReceiverName = this.talk_name;
                    message_body.SenderName = sessionStorage.getItem('username');
                    var send_body={}
                    send_body.id =  parseInt(sessionStorage.getItem('id'));
                    send_body.ReceiverType = this.receiver_type;
                    send_body.ReceiverId = this.receiver_id;
                    send_body.SenderName = sessionStorage.getItem('username');
    
                    // send_body.ReceiverId = 1517;
                    send_body.ResourcePath = this.im_send_news;
                    send_body.SendTime= this.im_now_date;
                    send_body.ReceiverName = this.talk_name;
                    send_body.MsgType = 1;
                    window.console.log(send_body);
                    this.websocket.send(JSON.stringify(send_body));
                    // if(this.receiver_type ==1){
                    //     var local_num=this.receiver_type+'.'+sessionStorage.getItem('id')+'_'+this.receiver_id;
                    //     window.console.log(local_num)
                    //     this.local_mynews = JSON.parse(localStorage.getItem(local_num));
                    //     window.console.log(this.local_mynews);
                    //     if(this.local_mynews ==null){
                    //         this.local_mynews = [];
                    //         this.local_mynews.push(message_body);
    
                    //     }else{
                    //         this.local_mynews.push(message_body);
                    //     }
                    //     localStorage.setItem(local_num, JSON.stringify(this.local_mynews));
                    //     window.console.log( this.local_mynews)
    
                    // }else{
                    //     window.console.log('wwwwwwwwwwwwwwwww')
                    //     // this.local_groupnews = JSON.parse(localStorage.getItem(local_num));
                    //     // window.console.log(this.local_groupnews);
                    //     // if(this.local_groupnews ==null){
                    //     //     this.local_groupnews = [];
                    //     //     this.local_groupnews.push(message_body);
    
                    //     // }else{
                    //     //     this.local_groupnews.push(message_body);
                    //     // }
                    //     // localStorage.setItem(local_num, JSON.stringify(this.local_groupnews));
                    //     // window.console.log( this.local_groupnews)
                    // }
    
                    this.im_send_news='';
    
                }else{  
                        this.im_create_time();
                        let file_body ={}
                        this.im_send_obj=1;
                        file_body.id = sessionStorage.getItem('id');
                        file_body.ReceiverType = this.receiver_type;
                        file_body.ReceiverId = this.receiver_id;
                        file_body.im_send_obj = this.im_send_obj;
                        file_body.SendTime = this.im_now_date;
                        file_body.time_id = (new Date()).getTime();
                        var formData = new FormData() // FormData 对象
                        window.console.log(this.$refs.resource.files[0])
                        formData.append("file", this.$refs.resource.files[0])
                        // formData.append("params", params)
                        // formData.append("id", 1503)
                        // formData.append("ReceiverId", 1500)
                        // formData.append("ReceiverType", 1)
                        formData.append("id", sessionStorage.getItem('id'))
                        formData.append("ReceiverId", this.receiver_id)
                        formData.append("ReceiverName", this.talk_name)
                        formData.append("SendTime", this.im_now_date)
                        formData.append("SenderName",sessionStorage.getItem('username'))
                        formData.append("ReceiverType", this.receiver_type)
                        var file_num=this.receiver_type+'.'+sessionStorage.getItem('id')+'_'+this.receiver_id;
                        window.console.log(file_num) 
                        file_body.MsgType = this.file_type;
                        file_body.file_name = this.file_name;
                        window.console.log( this.file_type)
                        // this.$axios.post('http://10.0.18.132:10000/upload', formData,{ headers: 
                        this.$axios.post('/upload', formData,{ headers: 
                                {
                                "Authorization" : sessionStorage.getItem('setSession_id'),
                                'Content-Type': 'multipart/form-data'
                                }
                                }).then((response) => {
                                this.upload_show = false;
                                this.im_box_show = true;
                                this.file_name =''
                             }).catch(function (error) {
                            window.console.log(error)
                        })
    
                }
                this.$nextTick(()=>{
                        var scrollHeight = $('.im_message').prop("scrollHeight");
                        // window.console.log(scrollHeight)
                        $('.im_message').scrollTop(scrollHeight,800);
                    })
                },
                im_send_close(){
                    this.im_show = false;
                    window.console.log('111'); 
//                     var audio_poc = { "request": "join", "room":371, "371": '1687' };
//  this.mixertest.send({"message": audio_poc});
                  
                },
            //     duankai(){
            //         window.console.log('duankai')
            //         var leave_room ={"request": "leave"}

            //         this.mixertest.send({"message": leave_room});


            //     },
            //     lianjie(){
            //         window.console.log('lianjie');
            //         var audio_poc2 = { "request": "join", "room":372, "display": '1687' };
		    // this.mixertest.send({"message": audio_poc2});
            //     },
                im_close(){
                    this.im_show = false;
                    // this.mapshow=true
                    // this.closeWebSocket();
                    // this.websocket.close()
                },
                socket_cloes(){
                    this.websocket.close()
                },
                filechange(){
                    this.im_box_show = false;
                    this.upload_show = true,
                  window.console.log(2222);
                  this.file_name = this.$refs.resource.files[0].name;
                //   this.file_type = 
                 var filepath= $('#im_resource').val();
                 var fileType = this.getFileType(filepath)
                 window.console.log(fileType)
                 if("jpg" == fileType || "jpeg" == fileType  || "png" == fileType || "gif" == fileType){
                    this.file_type = 2
                 }else if("mp3" == fileType ){
                    this.file_type = 3
                 }else if("mp4" == fileType || "rmvb" == fileType || "avi" == fileType || "ts" == fileType){
                    this.file_type = 4
                 }else{
                    this.file_type = 5
                 }
                 if(this.file_type == 2){
                     window.console.log(111)
                     var preview_image=document.getElementById('preview_img')
                     var reader =new FileReader()
                     reader.readAsDataURL(this.$refs.resource.files[0])
                     reader.onload =function(e){
                         window.console.log(e);
                         preview_image.src=e.target.result

                     }

                 }
             
                 
                },
                // 判断文件类型
                getFileType(filepath) {
                    var startIndex = filepath.lastIndexOf(".");
                    if(startIndex != -1)
                        return filepath.substring(startIndex + 1, filepath.length).toLowerCase();
                    else return "";
                },
                dele_file(){
                    this.upload_show = false,
                    this.im_box_show = true;
                },
                setOnopenMessage(e){
                    window.console.log('....................连接了') 
                    // this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
                    //             { headers: 
                    //             {
                    //             "Authorization" : sessionStorage.getItem('setSession_id')
                    //             }
                    //             })
                    //             .then((response) =>{
                    //              window.console.log(response);
                    //             localStorage.setItem('group_list', JSON.stringify(response.data.group_list));  
                    //             this.local_group_list = response.data.group_list    
                    //             })
                    //             .catch( (error) => {
                    //             window.console.log(error);
                    //             });  
                },
                setOncloseMessage(e){
                    window.console.log('....................关闭了')
                    window.console.log(e)
                    this.refresh_show = true;
                },
                setOnmessageMessage (e){ //数据接收
                    window.console.log(e)
                    window.console.log(e.data)
                    if(JSON.parse(e.data).DataType == 2){
                        window.console.log('---------------------------------p')
                        var off_line =JSON.parse(e.data).offlineImMsgResp;
                        
                        if(Object.keys(off_line).length==0){
                            window.console.log(333333)
                            window.console.log(off_line)
                        }else{
                            window.console.log(222222)
                            window.console.log(off_line.offlineImMsgs);
                            localStorage.setItem('off_line_local', JSON.stringify(off_line.offlineImMsgs));
                            
                        }
                        window.console.log(off_line)
                    }else if(JSON.parse(e.data).DataType == 3){
                        const redata = JSON.parse(e.data).imMsgData;
                        window.console.log(redata)
                        // if(redata.MsgType != 1){
                        //     redata.ResourcePath=  "http://"+ redata.ResourcePath
                        // }
                            redata.time_id = (new Date()).getTime();
                            // redata.time=this.im_now_date;
                            window.console.log(this.receiver_id)
                            window.console.log(this.local_mymessage);
                            if(redata.offlineImMsgs == undefined){
                                window.console.log(2222222222222)
                            }else{
                                window.console.log(1111111111111111111111)
                            }
                            if(redata.ReceiverType ==1){
                                window.console.log(sessionStorage.getItem('id'))
                                var your_local_name =''
                                if(redata.id == this.receiver_id ||redata.id == parseInt(sessionStorage.getItem('id'))){
                                      
                                       if(redata.id == parseInt(sessionStorage.getItem('id'))){
                                        redata.im_send_obj = 1;
                                        your_local_name=redata.ReceiverType+'.'+redata.id+'_'+redata.ReceiverId;
                                       window.console.log(your_local_name)
                                       }else{
                                        redata.im_send_obj = 2;
                                        your_local_name=redata.ReceiverType+'.'+redata.ReceiverId+'_'+redata.id;
                                          window.console.log(your_local_name)
                                       }
    
                                        this.local_mynews = JSON.parse(localStorage.getItem(your_local_name));
                                        if(this.local_mynews  ==null){
                                            this.local_mynews = [] ;
                                            this.local_mynews.push(redata);
                                        }else{
                                            this.local_mynews.push(redata);
                                        }
                                        localStorage.setItem(your_local_name, JSON.stringify(this.local_mynews));                       
                                }else{
                                    your_local_name=redata.ReceiverType+'.'+redata.ReceiverId+'_'+redata.id;
                                    window.console.log(redata.id)
                                    window.console.log(redata.ReceiverId)
                                        redata.im_send_obj = 2;
                                        window.console.log(888888888888888888)
                                        window.console.log(redata);
                                        var other_num = JSON.parse(localStorage.getItem(your_local_name));
                                        if(other_num   ==null){
                                            other_num = [] ;
                                            other_num .push(redata);
                                        }else{
                                            other_num.push(redata);
                                        }
                                        localStorage.setItem(your_local_name, JSON.stringify(other_num));
                                        this.other_name = redata.id;
                                        this.other_show =true;
                                        this.other_send_type = redata.ReceiverType;
                                        this.other_send_id = redata.id
                                        this.other_group_id =redata.ReceiverId
                                }
                            }else{
                                if(redata.id ==sessionStorage.getItem('id')){
                                    redata.im_send_obj = 1;
                                }else{
                                    redata.im_send_obj = 2;
                                }
                                // if(redata.MsgType !=1){
                                //     redata.ResourcePath = "http://"+redata.resourcePath;
                                // }
                                var group_local_name=redata.ReceiverType+'.'+sessionStorage.getItem('id')+'_'+redata.ReceiverId;
                                if(redata.ReceiverId ==this.receiver_id){
                                    window.console.log(555555)
                                    window.console.log(group_local_name)
    
                                    window.console.log(JSON.parse(localStorage.getItem(group_local_name)))
    
                                    this.local_groupnews = JSON.parse(localStorage.getItem(group_local_name));
                                        if(this.local_groupnews  ==null){
                                            this.local_groupnews = [] ;
                                            this.local_groupnews.push(redata);
                                        }else{
                                            this.local_groupnews.push(redata);
                                        }
                                        localStorage.setItem(group_local_name, JSON.stringify(this.local_groupnews));                                
                                }else{
                                    window.console.log(9999999)
                                    window.console.log(redata);
                                    var other_group_num = JSON.parse(localStorage.getItem(group_local_name));
                                    if(other_group_num  ==null){
                                        other_group_num = [] ;
                                        other_group_num .push(redata);
                                    }else{
                                        other_group_num.push(redata);
                                    }
                                    localStorage.setItem(group_local_name, JSON.stringify(other_group_num));
                                    this.other_name = redata.id;
                                    this.other_show =true;
                                    this.other_send_type = redata.ReceiverType;
                                    this.other_send_id = redata.id
                                    this.other_group_id =redata.ReceiverId                        
                                }
    
                            }
                            this.$nextTick(()=>{
                                var scrollHeight = $('.im_message').prop("scrollHeight");
                                window.console.log(scrollHeight)
                                $('.im_message').scrollTop(scrollHeight,200);
                                })
                    }else if(JSON.parse(e.data).DataType == 5){
                        this.off_line_show =true;
                        // var off_device=JSON.parse(e.data).logoutNotify.group_list[0].usr_list;
                        // window.console.log(off_device)
                        // for(var i=0;i<off_device.length;i++){
                        //     if(off_device[i].online == 2){
                        //         this.off_line_name =off_device[i].name;
                        //     }
                        // }
                        var off_device=JSON.parse(e.data).Notify.userInfo.name;
                        this.off_line_name =off_device;
                        window.console.log('````````````````````````78')
                        this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
                                { headers: 
                                {
                                "Authorization" : sessionStorage.getItem('setSession_id')
                                }
                                })
                                .then((response) =>{
                                 window.console.log(response);
                                localStorage.setItem('group_list', JSON.stringify(response.data.group_list));  
                                this.local_group_list = response.data.group_list    
                                })
                                .catch( (error) => {
                                window.console.log(error);
                                }); 
                                window.console.log('````````````````````````89')
                    }else if(JSON.parse(e.data).DataType == 6){
                        this.online_show =true;
                        var online_device=JSON.parse(e.data).Notify.userInfo.name;
                        this.online_name =online_device;
                        window.console.log('````````````````````````178')
                        this.$axios.get('/account/'+sessionStorage.getItem('loginName'),
                                { headers: 
                                {
                                "Authorization" : sessionStorage.getItem('setSession_id')
                                }
                                })
                                .then((response) =>{
                                 window.console.log(response);
                                localStorage.setItem('group_list', JSON.stringify(response.data.group_list));  
                                this.local_group_list = response.data.group_list    
                                })
                                .catch( (error) => {
                                window.console.log(error);
                                }); 
                                window.console.log('````````````````````````289')
                    }
                    // else if(JSON.parse(e.data).DataType == 7){
                    //     this.sos_show = true;
                    //     window.console.log(JSON.parse(e.data))
                    //     this.sos_name=JSON.parse(e.data).Notify.userInfo.name
                    //     var sos_location = JSON.parse(e.data).Notify.userLocation;
                    //     var gps_sos = {};
                    //     gps_sos.lng=sos_location.longitude;
                    //     gps_sos.lat=sos_location.latitude;
                    //     if(this.mapshow == true){
                    //         this.$refs.soschild.sos_point(gps_sos);
                    //     }else{
                    //          this.$refs.googlechild.sos_google(gps_sos);
                    //      }
                        
                    // }
    
                },
                look_other(){
                    window.console.log(this.other_send_type)
                    window.console.log(this.other_send_id)
                    window.console.log(this.group_list)
                    if(this.im_show ==false){
                        this.im_show =true
                    }
                    if(this.other_send_type ==1){
                        this.im_group_show = false;
                        this.im_device_show = true;
                        var local_decice_list = this.local_device_list;
                        for(var i=0;i<local_decice_list.length;i++){
                            if(local_decice_list[i].id == this.other_send_id){
                                this.talk_name = local_decice_list[i].user_name;
                            }
                        }
                        var look_other_self =this.other_send_type+'.'+this.other_group_id+'_'+this.other_send_id;
                        window.console.log(look_other_self)
                        this.im_record_show =false;
                         this.im_self_show =false;
                         this.im_self_show =true;
                        this.local_mynews = JSON.parse(localStorage.getItem(look_other_self));
                        this.receiver_id = this.other_send_id;
                        this.talk_id = this.other_send_id 
                    }else{
                        this.im_device_show = false;
                        this.im_group_show = true;
                        window.console.log(this.group_list)
                        var local_group_list = this.group_list;
                        for(var j=0;j<local_group_list.length;j++){
                            if(local_group_list[j].group_info.id == this.other_group_id){
                                this.talk_name = local_group_list[j].group_info.group_name;
                                window.console.log(222222)
                            }
                        }
                        var look_other_group =this.other_send_type+'.'+sessionStorage.getItem('id')+'_'+this.other_group_id;
                        window.console.log(look_other_group);
                        this.im_self_show =false;
                        this.im_record_show =false;
                        this.im_record_show =true;
                        this.local_groupnews =JSON.parse(localStorage.getItem(look_other_group));
                        this.receiver_id = this.other_group_id;
                        this.talk_id = this.other_group_id
                        
                    }
                    // this.talk_name=
                    
                    
                    // this.talk_id = this.other_send_id;
                    this.receiver_type = this.other_send_type;
                    this.other_show =false;
                },
                look_sos(){
                    var ccc= {
                            "localTime": 1557160225,
                            "longitude":114.1153786184,
                            "latitude": 22.4644910453,
                            "speed": 123.45646,
                            "course": 123
                    }
                    var a ={};
                    a.lng=ccc.longitude;
                    a.lat=ccc.latitude;
                    if(this.mapshow == true){
                        this.$refs.soschild.sos_point(a);
                    }else{
                        this.$refs.googlechild.sos_google(a);
                    }
                   
                    this.sos_show = false
                },
                off_line_button(){
                    this.off_line_show=false;
                },
                online_button(){
                    this.online_show =false;
                    
                },
                // websocket断开重连
                chat_refresh(){
                    this.refresh_show = false;
                    this.$router.go(0)
                },
                // 时间
                im_create_time(){
                    let timeStamp= new Date();
                    let year = timeStamp.getFullYear();
                    let month =timeStamp.getMonth() + 1 < 10? "0" + (timeStamp.getMonth() + 1): timeStamp.getMonth() + 1;
                    let date =timeStamp.getDate() < 10? "0" + timeStamp.getDate(): timeStamp.getDate();
                    let hh =timeStamp.getHours() < 10? "0" + timeStamp.getHours(): timeStamp.getHours();
                    let mm =timeStamp.getMinutes() < 10? "0" + timeStamp.getMinutes(): timeStamp.getMinutes();
                    let ss =timeStamp.getSeconds() < 10? "0" + timeStamp.getSeconds(): timeStamp.getSeconds();
                    this.im_now_date =year + "年" + month + "月" + date +"日"+" "+hh+":"+mm+":"+ss;
                },
                im_audio(){
                    window.console.log(this.talk_id)
                    window.console.log(this.receiver_type)
                    if(this.receiver_type == 1){
                        var body ={
                        'uid':this.talk_id,'type':0,"name":sessionStorage.getItem('username'),"isVideo":'false','isAccept':0,
                        'myId':sessionStorage.getItem('id')
                        }
                        window.console.log(body)
                        // this.videocall.send({'user_call':body})
                    }
    
                    
                },
                im_video(){
                    window.console.log(this.talk_id)
                    window.console.log(this.receiver_type);
                    if(this.receiver_type == 1){
                        this.audio_answer = false;
                            var body ={
                            'uid':this.talk_id,'type':0,"name":sessionStorage.getItem('username'),"isVideo":'true','isAccept':0,
                            'myId':sessionStorage.getItem('id')
                            }
                            window.console.log(body)
                            // this.videocall.send({'user_call':body})                   
                    }
                },
                im_audio_bridge(){
                    this.talkbackVisible=true
                            // this.im_content_show=false;
                            // this.im_talkback_show=true
                },
                talkback_submit(){
                    window.console.log(this.talkback_id);
                    var talkroom = this.talkback_id;
                    var my_roomname = sessionStorage.getItem('id')
                    var talk_request = { "request": "join", "room":talkroom, "display": my_roomname };
                    window.console.log(talk_request);
                    // this.mixertest.send({"message": talk_request}); 
                    this.im_content_show=false;
                    this.im_talkback_show=true;
                    this.talkbackVisible=false;

                },

                out_talkback(){
                    this.im_talkback_show=false;
                    this.im_content_show=true;
                    var out_talkroom ={"request": "leave"}
                //  this.mixertest.send({"message": out_talkroom});
                    this.talk_id = this.talkback_id;
                    this.receiver_id = this.talk_id
                    this.receiver_type = 2;
                    var im_group_num= this.receiver_type+'.'+sessionStorage.getItem('id')+'_'+this.talk_id
                    window.console.log(im_group_num);
                    window.console.log(JSON.parse(localStorage.getItem(im_group_num)));
                    this.local_groupnews =JSON.parse(localStorage.getItem(im_group_num));
                },
                talk_active(){
                     this.talkactive=true;
                    window.console.log(1)
                    this.mouse_show=false;
                },
                talk_actived(){
                     this.talkactive=false;
                    window.console.log(2)
                    this.mouse_show=true;
                },                
                map_show(){
                },
                // 地图切换
                baidu_selected(){
                  this.mapshow=true;
                },
                google_selected(){
                    this.mapshow=false;
                },
    
        },
        computed:{
                group_list(){
                    // return JSON.parse(localStorage.getItem('group_list'));
                    return this.local_group_list;
                },
                device_member(){
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
                    let transfer_newData = [];
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
                    let modified_newData = [];
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
                    let  group_select_device = [];
                    let  group_selected_num = this.select_group_nummber
                    if( group_selected_num !== 'a'){
                        group_select_device = this.local_group_list;
                        group_select_device = group_select_device[group_selected_num].device_infos;
                        let group_select_length = group_select_device.length;
                        if(group_select_length == 1 ){
                            group_select_device =[];
                        }else{
                           group_select_device = group_select_device.slice(1,group_select_length)
    
                        }
                    }
    
                    return   group_select_device
                },
                // 组名
                device_group_name(){
                    let group_selected_member = this.local_group_list;
                    let  group_selected_name = this.select_group_nummber
                    if( group_selected_name !== 'a'){
                        group_selected_member = group_selected_member[group_selected_name].group_info.group_name
                    }
                    return   group_selected_member
                },
                local_device_list(){
                    return this.get_device_list
                },
                // 即时消息
                im_total_message(){
                    return this.local_mynews
                },
                im_total_group(){
                    return this.local_groupnews
                },
                im_total_history(){
                    return this.local_history
                },
                im_off_line(){
                    // var a = JSON.parse(localStorage.getItem('off_line_local'));
                    window.console.log(localStorage.getItem('off_line_local'))
                    if(localStorage.getItem('off_line_local') == 'undefined'){
                        return []
                    }else{
                        var a = JSON.parse(localStorage.getItem('off_line_local'));
                        var b =[]
                    if (a !== null){
                        if(a.hasOwnProperty('offlineSingleImMsgs')){
                         window.console.log('```````456')
                         for(var i =0;i<a.offlineSingleImMsgs.length;i++){
                             b.push(a.offlineSingleImMsgs[i])
                         }
                     }
                     if(a.hasOwnProperty('offlineGroupImMsgs')){
                        for(var j =0;j<a.offlineGroupImMsgs.length;j++){
                             b.push(a.offlineGroupImMsgs[j])
                         }                
                        }
                     for (var k = 0; k<b.length;k++){
                         b[k].id=k
                     }   
                    }
                     
                    //   return this.jiadata
                    return b
                    }

                },
                // sosGps(){
                //     this.sosdata = {l:114.07,s:22.62}
                //     return this.sosdata
                // }
                // 在线
                
        },
        created(){
          
        },
        beforeMount(){
              this.get_group_list();
              window.console.log( JSON.parse(localStorage.getItem('group_list')))
              window.console.log( JSON.parse(localStorage.getItem('device_list')))
               
        },
        mounted(){ 
            this.audiobridge_serve(); 
            window.console.log( JSON.parse(localStorage.getItem('group_list')))
            this.video_server();    
            this.createWebSocket()
            window.console.log(this.group_list)
        },
        updated() {
                
        },
        beforeDestroy(){
            this.socket_cloes()
        }
    }
    </script>
    
    <style >
    #big_box{
        position: relative;
        height: 100%;
    }
    /* 地图切换 */
    .map_select{
        z-index: 888;
        position: absolute;
        right: 9px;
        top: 57px;
        background-color: white;
        height: 29px;
        line-height: 29px;
    }
    .map_select div{
        display: inline-block;
        padding-left: 5px;
        padding-right: 5px;
    }
    .baidumap{
        padding-right: 10px;
        font-size: 12px;
        cursor: pointer;
    }
    .googlemap{
        font-size: 12px;
        cursor: pointer;
    }
    .map_active{
        color: white;
        background-color: #8ea8e0;
    }
    #monotor_content{
        /* z-index: 666; */
        position: absolute;
        top: 0px;
    }
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
        overflow: auto;
        height: 600px;
    }
    .members{
        height:35px;
        position: relative;
    }
    .im_members {
        height:37px;
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
        cursor: pointer;
    
    }
    .device_detail_close{
        float: right;
        padding-right: 8px;
        cursor: pointer;
    }
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
        margin: 7px;
        border: 1px solid #ccc;
        width: 80px;
        height: 80px;
        position: relative;
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
        text-align: center;
        position: absolute;
        bottom: 0px;
        width: 80px;
        overflow: hidden;
     
    }
    .off_color{
        background-color: #ccc !important;
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
        border: 1px solid #ccc;
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
    .device_pic,.group_pic{
        width: 80px;
        height: 60px;
        text-align: center;
    }
    .control_menum{
        background-color: white;
        position: absolute;
        border: 1px solid #c4c4c4;
        height: 87px;
        width: 78px;
        top: 45px;
        left: 27px;
        z-index: 6;
        
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
    .device_detail_name{
        display: inline-block;
        margin-left: 10px;
        margin-top: 5px;
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
    .ungroup:hover,.modification:hover,.editor_member:hover,.control_voice:hover,.control_vedio:hover,.control_look:hover,.control_text:hover{
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
    .video_room{
        width: 665px;
        height: 314px;
        background-color: rgba(94, 94, 94, 0.3);
        position: absolute;
        top: 5%;
        left: 50%;
        z-index: 666;
    }
    #remotevideo{
        position: absolute;
        left:127px
        /* top: 10px; */
    }
    #myvideo{
    position: absolute;
        bottom: 43px;
        z-index: 888;
        right: 4px;
        z-index: 777;
    }
    .video_close{
        position: absolute;
        bottom: 6px;
        text-align: center;
        left: 292px;
    }
    .audio_close{
        position: absolute;
        bottom: 18px;
        text-align: center;
        left: 74px; 
    }
    #myaudio{
     position: absolute;
        bottom: 43px;
        z-index: 888;
        right: 4px;
        z-index: 777;   
    }
    #remoteaudio{
            position: absolute;
        left:127px
    }
    .audio_room{
        width: 219px;
        height: 266px;
        background-color: black;
        position: absolute;
        top: 5%;
        left: 50%;
        z-index: 666;
    }
    .audio_tittle{
    color: white;
    font-size: 14px;
    text-align: center;
    margin-top: 10px;
    }
    .video_tittle{
        position: absolute;
        font-size: 14px;
        margin-left: 17px;
        margin-top: 14px;
    }
    #audio_box{
         width: 86px;
         height: 91px;
         border: 1px solid #ccc;
         margin-top: 54px;
         margin-left: 68px;
         background-color: white;
         /* display: none */
    }
    .audio_img{
        width: 83px;
        height: 86px;
    }
    .audio_loding{
        color: white;
        font-size: 14px;
        text-align: center;
        margin-top: 50px;
    }
    .im_box{
        width: 710px;
        height: 448px;
        position: absolute;
        top: 50px;
        left: 400px;
        box-shadow: 0 0 10px 5px rgba(207, 207, 209, 0.8);
        background-color: #f0f0f0;
        /* z-index: 662; */
        /* border: 1px solid #ccc;
        box-sizing: border-box; */
    }
    .im_aside{
        height: 448px;
        width: 40px;
        background-color: #206ba2;
        float: left;
        border-right: 1px solid #ccc;
        box-sizing: border-box;
    }
    .im_group{
        height: 448px;
        width: 233px;
        background-color: #f2f2f2;
        float: left; 
    }
    .im_content{
        background-color: white;
        height: 403px;
        width: 437px;
        float: right;
        margin-top: -1px;
    }
    .im_talkback{
        background-color: white;
        height: 403px;
        width: 437px;
        float: right;
        margin-top: -1px;    
    }
    .im_aside_logo{
        margin-top: 10px;
        width: 40px;
        height: 40px;
        text-align: center;
    }
    .im_aside_img{
        width: 30px;
        height: 30px;
    }
    .im_recent_icon,.im_device_icon,.im_group_icon{
        margin-top: 10px;
        width: 40px;
        height: 40px;
        /* background-color: brown; */
        text-align: center;
        color: white;
        cursor: pointer;
    }
    .im_recent_icon:hover,.im_device_icon:hover,.im_group_icon:hover{
        color: white
    }
    .im_size{
        font-size: 34px; 
        color: black;   
    }
    .aside_active{
        color: white !important;
    }
    .im_group_list{
        background-color: #206ba2;
        height: 45px;
        text-align: center;
        font-size: 19px;
        line-height: 45px;
        color: white;
        border-bottom: 1px solid #ccc;
        box-sizing: border-box;
    }
    .im_overheigh{
        overflow: auto;
        height: 402px;
    }
    .im_message{
        height:278px;
        background-color:#fafafa;
        overflow-y: scroll;
    }
    .im_media{
        height: 31px;
        /* background-color: aqua; */
        border-top: 1px solid #ccc;
    }
    .im_media_icon{
        padding-left: 10px;
        cursor: pointer;
    }
    .im_record{
        display: inline-block;
        float: right;
        margin-right: 12px;
        cursor: pointer;
    }
    .im_record_icon{
        font-size: 21px;
    }
    .el-icon-caret-right{
        font-size: 20px;
    }
    .im_send_message{
        font-size: 13px;
        overflow: auto;
        width: 424px;
        height: 62px;
        resize: none;
        outline: none;
        border: none;
    }
    .upload_file{
        font-size: 13px;
        overflow: auto;
        width: 424px;
        height: 62px;
        resize: none;
        outline: none;
        border: none;  
    }
    #preview_img{
        height: 52px;
        width: 150px;
    }
    #dele_area{
        height: 15px;
        border: none;
        resize: none;
        outline: none;

    }
    .im_send_button{
       float: right; 
    }
    .im_send_close,.im_send_subimt{
        font-size: 14px;
        border: 1px solid #ccc;
        padding-left: 17px;
        padding-right: 17px;
        padding-top: 4px;
        padding-bottom: 4px;
        border-radius: 5px 5px 5px 5px;
        cursor: pointer;
        margin-right: 10px;
    }
    .im_send_subimt{
    
    }
    /* 消息队列 */
    .im_opposite_name{
        margin-left: 10px;
    }
    .im_opposite_pic{
        width: 35px;
        height: 35px;
        float: left;
        margin: 8px 10px 8px 10px;
        background-color: #ccc;
    }
    .im_group_name{
        margin-left: 10px;
        font-size: 12px;
    }
    .im_opposite_img,.im_self_img{
        width: 35px;
        height: 35px;
    }
    .im_opposite_msg{
        float: left;
        font-size: 13px;
        max-width: 316px;
        padding: 10px;
        background-color: #fff;
        position: relative;
        border: 1px solid #e5e5e5;
        border-radius: 3px;
        line-height: 13px;
        margin-top: 8px;
        white-space: normal;
        word-break: break-all;
    }
    .im_opposite_corner{
        width: 7px;
        height: 7px;
        border: 1px solid #e5e5e5;
        transform: rotate(45deg);
        position: absolute;
        top: 12px;
        left: -4px;
        background-color: #fff;
        border-top: none;
        border-right: none;
    }
    .im_self_pic{
        width: 35px;
        height: 35px;
        float: right;
        margin: 8px 10px 8px 10px;
        background-color: #ccc; 
        border-radius: 3px; 
    }
    .clearfix:after {
        content: ".";
        display: block;
        height: 0;
        font-size: 0;
        clear: both;
        visibility: hidden;
    }
    .im_self_msg{
        float: right;
        font-size: 13px;
        max-width: 335px;
        padding: 10px;
        position: relative;
        border-radius: 3px;
        line-height: 13px;
        margin-top: 8px;
        white-space: normal;
        word-break: break-all;
        background-color: #87daef;
    }
    .im_self_corner{
         width: 7px;
        height: 7px;
        border: 1px solid transparent;
        transform: rotate(45deg);
        position: absolute;
        top: 12px;
        right: -3px;
        background-color: #87daef;
        border-bottom: none;
        border-left: none;
    }
    .im_history{
        width: 332px;
        height: 448px;
        /* float: right; */
        position: absolute;
        /* right: 0px; */
        top: 0px;
        background-color: white;
    
    }
    .im_history_top{
        height: 46px;
        background-color: white;
        border-bottom: 1px solid #ccc;
    }
    .im_history_text,.im_history_file,.im_history_talk{ 
        display: inline-block;
        float: left;
        margin-left: 6px;
        font-size: 13px;
        text-align: center;
        width: 46px;
        cursor: pointer;
    }
    .im_history_talk{ 
      width: 72px;
    }
    .active_history{
        border: 1px solid #ccc;
        background-color: #e6e6e6;
        border-radius: 3px;
    }
    .im-history-wrap{
        position: absolute;
        top: 47px;
        bottom: 16px;
        left: 15px;
        right: 3px;
        padding-right: 17px;
        overflow: auto;
        background-color: white;
    }
    .im_history_info{
        border-bottom: 1px solid #e5e5e5;
        padding-top: 7px;
        padding-bottom: 7px;
        font-size: 14px;
    }
    .im_history_name{
        width: 277px;
        height: 20px;
    }
    .im_talk_content{
        width: 277px;
        height: 30px;  
    }
    .im_voice-wrap{
        width: 140px;
        height: 30px;
        border-radius: 2px;
        border: 1px solid #e5e5e5;
        line-height: 30px;
    }
    .im_voice_div{
        height: 26px;
        width: 26px;
        margin-top: 3px;
        display: inline-block;
    }
    .im_voice_num{
        float: right
    }
    .im_voice_img{
     height: 23px;
        width: 25px;
    }
    .im_voice_gif{
        height: 23px;
        width: 25px;
        z-index: 666;
        position: absolute;
        left: 1px;
        background-color: white;
    }
    .im_record_name{
       float: left;
    }
    .im_record_time{
        float: right
    }
    .im_history_title{
        font-size:5px;
    }
    .dele_file{
        border: 1px solid #ccc;
        width: 68px;
        text-align: center;
        cursor: pointer;
    }
    .im_get_img{
        max-width: 253px;
    }
    /* 地图 */
    #map{
        height: 100%
    }
    .other_div{
        height: 69px;
        width: 307px;
        position: absolute;
        right: 0px;
        top: 290px;
        z-index: 666;
        background-color: white;
        border-radius: 21px;
    }
    .sos_div{
        height: 69px;
        width: 307px;
        position: absolute;
        right: 0px;
        top: 200px;
        z-index: 666;
        background-color: red;
        border-radius: 21px;   
        text-align: center;
    }
    .sos_title{
        color: white;
    }
    .sos_body{
        margin-top: 3px
    }
    .other_answer{
        display: inline-block;
        float: right;
        margin-right: 23px;
        background-color: cornflowerblue;
        padding: 5px;
        border-radius: 5px;
        cursor: pointer;
    }
    /* 视频语音请求 */
    .video_apply{
        width: 524px;
        height: 176px;
        background-color: white;
        position: absolute;
        top: 24%;
        left: 41%;
    }
    .off_line_div{
        width: 274px;
        height: 176px;
        background-color: white;
        position: absolute;
        top: 24%;
        left: 41%;
        z-index: 55;
    }
    .chat_div{
        width: 274px;
        height: 176px;
        background-color: white;
        position: absolute;
        top: 24%;
        left: 41%;
        z-index: 999;   
    }
    .off_title{
        text-align: center;
        margin-top: 10px;
    }
    .chat_title{
        text-align: center;
        margin-top: 37px;
    }
    .off_body{
        text-align: center;
        margin-top: 30px;
    }
    .off_foot{
        text-align: center;
        margin-top: 40px;
    }
    .video_apply_title{
        text-align: center;
        font-size: 20px;
        margin-top: 45px;
    }
    .video_apply_body{
        text-align: center;
        margin-top: 46px;
    }
    .v-modal{
        /* display: none */
    }
    .out_talkback{
        float: right;
        margin-top: 10px;
        margin-right: 10px;
        background-color: #409eff;
        padding: 5px;
        cursor: pointer;
    }
    .talkback_icon{
        font-size: 100px;
    margin-top: 12px;
    display: inline-block;
    margin-left: 12px;
    }
    .talkback_div{
        display: inline-block;
        width: 140px;
        height: 140px;
        background-color: #ccc;
        text-align: center;
        margin-left: 147px;
        margin-top: 24px;
        border-radius: 101px;
        cursor: pointer;
    }
    .talkback_active{
        background-color: #409eff !important;
        color: white !important
    }
    .talkactive_tip{
        margin-top: 174px;
        text-align: center;
    }
    </style>
    