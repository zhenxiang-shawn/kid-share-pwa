import request from '../../utils/request'
import { loginFormData, loginResponseData, userInfoReponseData } from './type'

enum API {
  // 用户
  USER_LOGIN = '/user/login',
  USER_REGISTER = '/user/register',
  USER_INFO = '/user/info',
  USER_UPDATE = '/user/update',
  USER_LOGOUT = '/user/logout',
}

// user login
export const reqLogin = (data: loginFormData) =>
  request.post<any, loginResponseData>(API.USER_LOGIN, data)

// user sign up
export const reqRegister = (data: any) =>
  request.post<any, any>(API.USER_REGISTER, data)

// grab user information
export const reqUserInfo = () =>
  request.get<any, userInfoReponseData>(API.USER_INFO)

// logout
export const reqLogout = () => request.post<any, any>(API.USER_LOGOUT)
