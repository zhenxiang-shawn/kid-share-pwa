export const SET_USER_INFO = (
  username: string,
  displayName: string,
  avatar: string,
  relation: string,
) => {
  SET_USER_NAME(username)
  SET_USER_DISPLAY_NAME(displayName)
  SET_USER_AVATAR(avatar)
  SET_USER_RELATION(relation)
}
/* Getter & Setter of username */
export const SET_USER_NAME = (username: string) => {
  localStorage.setItem('user.name', username)
}

export const GET_USER_NAME = () => {
  return localStorage.getItem('user.name')
}
/* Getter & Setter of user displayname */
export const SET_USER_DISPLAY_NAME = (displayName: string) => {
  localStorage.setItem('user.displayName', displayName)
}
export const GET_USER_DISPLAY_NAME = () => {
  return localStorage.getItem('user.displayName')
}
/* Getter & Setter of user avatar */
export const SET_USER_AVATAR = (avatar: string) => {
  localStorage.setItem('user.avatar', avatar)
}
export const GET_USER_AVATAR = () => {
  return localStorage.getItem('user.avatar')
}
/* Getter & Setter of user relation */
export const SET_USER_RELATION = (relation: string) => {
  localStorage.setItem('user.relation', relation)
}
export const GET_USER_RELATION = () => {
  return localStorage.getItem('user.relation')
}
