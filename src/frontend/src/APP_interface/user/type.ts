export interface loginFormData {
  username: string
  password: string
}

//定义登录接口返回数据类型
export interface loginResponseData {
  code: number
  message: string
  ok: boolean
  data: {
    token: string
    username: string
    display_name: string
    avatar: string
    relation: string
  }
}

//定义获取用户信息返回数据类型
export interface userInfoReponseData {
  code: number
  message: string
  ok: boolean

  data: {
    displayName: string
    avatar: string
    relation: string
  }
}
