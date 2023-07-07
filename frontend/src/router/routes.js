import Home from '../views/Home.vue'
// 务必要使用此种方法引用一个组件，否则wails打包后不能正常显示组件，我也不清楚怎么肥事儿
const routes = [
    {
        path: '/',
        component: Home
    },
    {
        name: 'fileList',
        path: '/fileList',
        component: () => import('../views/FileList.vue')
    },
    {
        name: 'download',
        path: '/download',
        component: () => import('../views/Download.vue')
    },
    {
        name: 'setting',
        path: '/setting',
        component: () => import('../views/Setting.vue')
    },
    {
        name: 'about',
        path: '/about',
        component: () => import('../views/About.vue')
    }
]

export default routes