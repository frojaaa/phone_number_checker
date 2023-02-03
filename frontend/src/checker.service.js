import axios from "axios";
import authHeader from "./auth.headers";

const CHECKER_URL = "http://localhost:5000/checker/run/"
const PASSWORD_CHECKER = "http://localhost:5000/checker/checkPassword"

class CheckerService {
    runChecker(numWorkers, lkLogin, lkPassword, botToken, tgUserID, inputFileDir, outputFileDir) {
        return axios.post(CHECKER_URL + `?numWorkers=${numWorkers}`, {
            lkLogin: lkLogin,
            lkPassword: lkPassword,
            botToken: botToken,
            tgUserID: tgUserID,
            inputFileDir: inputFileDir,
            outputFileDir: outputFileDir
        }, {headers: authHeader()});
    }

    checkPassword(username, password) {
        return axios.post(PASSWORD_CHECKER, {
            username: username,
            password: password
        }, {headers: authHeader()}).then(r => {
            r.status === 200 ? localStorage.setItem("enteredCheckerPassword", "true") : console.log(r)
        })
    }
}

export default new CheckerService();