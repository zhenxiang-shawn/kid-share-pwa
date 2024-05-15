export interface TimelineItemState {
  poster: string // 作者
  avatar: string // 头像
  title: string // 头衔
  date: string // 日期
  content: string // 内容
  images: Array<string> // 图片
}

export interface ContentState {
  contents: TimelineItemState[] //
}
