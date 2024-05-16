// 创建content 展示的小仓库
import { defineStore } from 'pinia'
import type { ContentState } from './types/types'

const useContentStore = defineStore('ContentStore', {
  state: (): ContentState => {
    return {
      contents: [],
    }
  },
  actions: {
    // 获取内容请求
    async contentUpdate() {
      // query 请求
      let result = { code: 200, data: 'data' }
      if (result.code === 200) {
        // update content
        return 'ok'
      } else {
        Promise.reject(new Error(result.data))
      }
    },
  },
})

export default useContentStore
