// 创建content 展示的小仓库
import { defineStore } from 'pinia'
import type { ContentState } from './types/types'
import { reqGetDiaries } from '../../APP_interface/diary_posts'

const useContentStore = defineStore('ContentStore', {
  state: (): ContentState => {
    return {
      contents: [],
    }
  },
  actions: {
    // 获取内容请求
    async getContent(page: number = 1, limit: number = 10) {
      let result = await reqGetDiaries(page, limit)
      if (result.code === 200) {
        console.log('result:', result)
        // update content
        this.contents = result.data.diaries
        return 'ok'
      } else {
        Promise.reject(new Error(result.message))
      }
    },
  },
})

export default useContentStore
