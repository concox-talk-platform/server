import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);
// const homePage = () => import('./Foo.vue');
export default new Router({
    routes: [
        // {
        //     path: '/',
        //     redirect: '/login'
        // },
        // {
        //     path: '/login',
        //     component: resolve => require(['../components/page/Login.vue'], resolve)
        // },
        {
            path: '/',
            redirect: '/homePage',
            component: resolve => require(['../components/common/Home.vue'], resolve),
            meta: { title: '自述文件' },
            children:[
                {
                    path: '/homePage',
                    component: () => import('../components/page/homePage.vue'),
                    meta: { title: '账户中心' }
                },
                {
                    path: '/usergroup',
                    component: resolve => require(['../components/page/UserGroup.vue'], resolve),
                    meta: { title: '用户' }
                },
                {
                    path: '/mapControl',
                    component: resolve => require(['../components/page/MapControl.vue'], resolve),
                    meta: { title: '控制台' }
                },
                {
                    path: '/area',
                    component: resolve => require(['../components/page/Area.vue'], resolve),
                    meta: { title: '区域' }
                },
                {
                    // 设备组件
                    path: '/equipment',
                    component: resolve => require(['../components/page/Equipment.vue'], resolve),
                    meta: { title: '设备' }
                },
                // 统计组件
                {
                    path: '/monitor',
                    component: resolve => require(['../components/page/Monitor.vue'], resolve),
                    meta: { title: '监控' }
                },
                //     // 发出信息组件
                // {    path: '/sendMessage',
                //     component: resolve => require(['../components/page/SendMesssage.vue'], resolve),
                //     meta: { title: '发出消息记录' }
                // },
                // {
                //     // 收到信息组件
                //     path: '/getMessage',
                //     component: resolve => require(['../components/page/receiveMessage.vue'], resolve),
                //     meta: { title: '收到消息记录' }    
                // },
           
                // {
                //     // 人员权限组件
                //     path: '/peoPermissions',
                //     component: resolve => require(['../components/page/PeoPermissions.vue'], resolve),
                //     meta: { title: '人员权限' }
                // },
                // {
                //     // 设备权限组件
                //     path: '/eqPermissions',
                //     component: resolve => require(['../components/page/EqPermissions.vue'], resolve),
                //     meta: { title: '设备权限' }
                // },
                // {
                //     // 权限页面
                //     path: '/permission',
                //     component: resolve => require(['../components/page/Permission.vue'], resolve),
                //     meta: { title: '权限测试', permission: true }
                // },
                // {
                //     path: '/404',
                //     component: resolve => require(['../components/page/404.vue'], resolve),
                //     meta: { title: '404' }
                // },
                // {
                //     path: '/403',
                //     component: resolve => require(['../components/page/403.vue'], resolve),
                //     meta: { title: '403' }
                // }
            ]
        },
     {
            path: '/login',
            component: resolve => require(['../components/page/Login.vue'], resolve)
        },
        {
            path: '*',
            redirect: '/404'
        }
    ],
    base:'/'
})
