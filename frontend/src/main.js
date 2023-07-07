import { createApp } from 'vue'
import App from './App.vue'
import naive from 'naive-ui'
import router from './router/index'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'


const app = createApp(App)
app.use(naive)
app.use(router)
app.use(ElementPlus, {})
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.mount('#app')
