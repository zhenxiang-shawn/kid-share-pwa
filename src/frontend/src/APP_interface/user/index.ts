import request from '../../utils/request'
import { loginFormData, loginResponseData, userInfoReponseData } from './type'

enum API {
  // 用户
  USER_LOGIN_URL = '/user/login',
  USER_REGISTER_URL = '/user/register',
  USER_INFO_URL = '/user',
  USER_UPDATE_URL = '/user/update',
  USER_LOGOUT_URL = '/user/logout',
}

// user login
export const reqLogin = (data: loginFormData) =>
  request.post<any, loginResponseData>(API.USER_LOGIN_URL, data)

// user sign up
export const reqRegister = (data: any) =>
  request.post<any, any>(API.USER_REGISTER_URL, data)

// grab user information
export const reqUserInfo = (username: string) =>
  request.get<any, userInfoReponseData>(`${API.USER_INFO_URL}/${username}`)

// logout
export const reqLogout = () => request.post<any, any>(API.USER_LOGOUT_URL)
