<script>
import { h } from "vue";
import { RouterLink } from "vue-router";
import { NIcon, NConfigProvider } from "naive-ui";
import {
  Home as HomeIcon,
  SettingsSharp as SettingIcon,
  AlertCircleSharp as AboutIcon,
  } from '@vicons/ionicons5'

import { lightTheme } from 'naive-ui'
import {ElNotification} from "element-plus";

const renderIcon = (icon) => {
  return () => h(NIcon, null, { default: () => h(icon) });
}

export default {
  components: {
    NConfigProvider
  },
  data() {
    return {
      collapsed: false,
      myTheme: null,
      switchTheme: false,
      menuOptions: [
        {
          label: () => h(
              RouterLink,
              {
                to: {
                  path: "/",
                }
              },
              { default: () => "主页" }
          ),
          key: "go-back-home",
          icon: renderIcon(HomeIcon)
        },
        {
          label: () => h(
              RouterLink,
              {
                to: {
                  name: 'fileList',
                  path: "/fileList",
                }
              },
              { default: () => "文件列表" }
          ),
          key: "go-back-fileList",
          icon: renderIcon(SettingIcon)
        },

        {
          label: () => h(
              RouterLink,
              {
                to: {
                  name: 'download',
                  path: "/download",
                }
              },
              { default: () => "下载列表" }
          ),
          key: "go-back-download",
          icon: renderIcon(SettingIcon)
        },
        {
          label: () => h(
              RouterLink,
              {
                to: {
                  name: 'setting',
                  path: "/setting",
                }
              },
              { default: () => "软件设置" }
          ),
          key: "go-back-setting",
          icon: renderIcon(SettingIcon)
        },
        {
          label: () => h(
              RouterLink,
              {
                to: {
                  name: 'about',
                  path: "/about",
                }
              },
              { default: () => "关于" }
          ),
          key: "go-back-about",
          icon: renderIcon(AboutIcon)
        },
      ],
      railStyle: ({
                    focused,
                    checked
                  }) => {
        const style = {};
        if (checked) {
          style.background = "#4B9D5F";
          if (focused) {
            style.boxShadow = "0 0 0 2px #d0305040";
          }
        } else {
          style.background = "#000000";
          if (focused) {
            style.boxShadow = "0 0 0 2px #2080f040";
          }
        }
        return style;
      }
    }
  },
  mounted() {
    this.myTheme = lightTheme
  },
  methods: {
    changeTheme(){
      ElNotification({
        title: '事件通知',
        message: 'This is a 挖掘机',
        type: 'info',
      })
    }
  }
}

</script>

<template>
  <n-config-provider :theme="myTheme">
    <n-space vertical size="large">
      <n-layout has-sider position="absolute">
        <n-layout-sider bordered collapse-mode="width" :collapsed-width="80" :width="150" :collapsed="collapsed"
                        show-trigger @collapse="collapsed = true" @expand="collapsed = false"
                        style="--wails-draggable:drag; opacity: 1;">
          <n-menu :options="menuOptions" :collapsed-width="64" :collapsed-icon-size="22" style="margin-top: 40px;" />
          <div class="switchBtnPar">
            <n-divider />
            <n-switch :rail-style="railStyle" v-model:value="switchTheme" @update:value="changeTheme()"
                      class="switchBtn">
              <template #checked>
                点
              </template>
              <template #unchecked>
                点
              </template>
            </n-switch>
          </div>
        </n-layout-sider>
        <n-layout-content>
          <router-view />
        </n-layout-content>
      </n-layout>
    </n-space>
  </n-config-provider>
</template>

<style>
.switchBtnPar {
  position: relative;
}

.switchBtn {
  position: absolute;
  left: 50%;
  transform: translate(-50%);
}

body {
  margin: 0;
}
</style>