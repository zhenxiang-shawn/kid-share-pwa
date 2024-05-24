import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import VueSetupExtend from 'vite-plugin-vue-setup-extend'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  let env = loadEnv(mode, process.cwd())
  return {
    plugins: [vue(), VueSetupExtend()],
    resolve: {
      alias: {
        '@': '/src', // 相对路径别名配置，使用 @ 代替 src
      },
    },
    // server: {
    //   host: '0.0.0.0',
    //   port: 8818,
    // },
    server: {
      proxy: {
        [env.VITE_APP_BASE_API]: {
          // target: env.VITE_SERVE,
          target: 'http://117.72.64.223:8818/',
          changeOrigin: true,
          // rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
      port: 8818,
    },
  }
})
