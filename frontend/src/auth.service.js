import axios from "axios";

const LOGIN_URL = "http://localhost:5000/login"
const REGISTER_URL = "http://localhost:5000/register"
const LOGOUT_URL = "http://localhost:5000/logout"


class AuthService {
    login(username, password) {
        return axios.post(LOGIN_URL, {
            username: username,
            password: password
        }).then(response => {
            console.log(response)
            if (response.data.token) {
                localStorage.setItem("token", JSON.stringify(response.data.token))
                localStorage.setItem("username", JSON.stringify(username))
                localStorage.setItem("enteredCheckerPassword", JSON.stringify(response.data.enteredCheckerPassword))
            }
            return response.data
        });
    }
    logout() {
        localStorage.removeItem("token")
    }
    register(username, password) {
        return axios.post(REGISTER_URL, {
            username: username,
            password: password
        });
    }
    getCurrentUser() {
        return JSON.parse(localStorage.getItem("user"));
    }
}

export default new AuthService();