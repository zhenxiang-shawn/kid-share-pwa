// 创建content 展示的小仓库
import { defineStore } from 'pinia'
import type { ContentState } from './types/types'
import { reqGetDiaries } from '../../APP_interface/diary_posts'
import { diaryData } from '../../APP_interface/diary_posts/type'

const GET_LOCAL_CONTENT = () => {
  return localStorage.getItem('content')
    ? JSON.parse(localStorage.getItem('content')!)
    : []
}
const SET_LOCAL_CONTENT = (content: Array<diaryData>) => {
  // sort content by time
  // 根据 date 属性排序
  // content.sort((a, b) => {
  //   const dateA = new Date(a.date).getTime() // 将日期字符串转换为时间戳
  //   const dateB = new Date(b.date).getTime() // 同上

  //   // 升序排序
  //   // return dateA - dateB

  //   // 如果需要降序排序，可以交换 a 和 b 的位置
  //   return dateB - dateA
  // })
  // append content to localstorage if it's not empty
  if (localStorage.getItem('content')) {
    let localContent = JSON.parse(localStorage.getItem('content')!)
    localContent.push(...content)
    localStorage.setItem('content', JSON.stringify(localContent))
  } else {
    // initial content to local storage if it's empty
    localStorage.setItem('content', JSON.stringify(content))
  }
}
export const CLEAR_LOCAL_CONTENT = () => {
  localStorage.removeItem('content')
}
const useContentStore = defineStore('ContentStore', {
  state: (): ContentState => {
    return {
      contents: GET_LOCAL_CONTENT(),
      contentsNum: 0,
    }
  },
  actions: {
    // 获取内容请求
    async updateContent(page: number = 1, limit: number = 20) {
      let result = await reqGetDiaries(page, limit)
      if (result.code === 200) {
        console.log('result:', result)
        if (result.data.diaries === null) {
          Promise.resolve('No more diaries')
          return
        }
        // update content
        SET_LOCAL_CONTENT(result.data.diaries)
        this.contentsNum = result.data.total
        return 'ok'
      } else {
        Promise.reject(new Error(result.message))
      }
    },
  },
})

export default useContentStore
