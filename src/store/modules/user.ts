import { defineStore } from "pinia"
import { GET_TOKEN, REMOVE_TOKEN } from "../../utils/token";

const useUserStore = defineStore('User', {
    state: () => {
        return {
            token: GET_TOKEN,
            username: '',
            avatar: '',
            title: '',

        }
    },
    actions: {
        // 用户登陆
        // 重置用户信息
        resetUserInfo() {
            REMOVE_TOKEN();
            this.username = ''
            this.avatar = ''
            this.title = ''
        }
    }
});

export default useUserStore;