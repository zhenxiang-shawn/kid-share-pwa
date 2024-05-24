import request from '../../utils/request'
import { diariesRequestData, diariesResponseData } from './type'

enum API {
  // 日记
  DIARY_POST_URL = '/diary',
  DIARIES_GET_URL = '/diaries',
}

export const reqPostDiary = (data: diariesRequestData) =>
  request.post<any, any>(API.DIARY_POST_URL, data)

export const reqGetDiaries = (page: number, limit: number) =>
  request.get<any, diariesResponseData>(
    `${API.DIARIES_GET_URL}?page=${page}&limit=${limit}`,
  )
