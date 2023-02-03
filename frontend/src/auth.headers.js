export default function authHeader() {
    const token = JSON.parse(localStorage.getItem("token"))
    console.log(token)
    if (token !== null) {
        return { Authorization: 'Bearer ' + token };
    } else {
        return {};
    }
}
