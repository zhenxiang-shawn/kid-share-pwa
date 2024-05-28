export interface diaryData {
  content: string
  username: string
  image_paths: Array<string>
  timestamp: Date
}

export interface diariesResponseData {
  data: {
    diaries: diaryData[]
    total: number
  }
  code: number
  message: string
  ok: boolean
}

export interface diariesRequestData {
  username: string
  content: string
  image_paths: Array<string>
}
