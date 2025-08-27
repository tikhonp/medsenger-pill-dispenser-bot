import {empty} from "./helpers"

class HttpClient {
    constructor(endpoint) {
        this.endpoint = endpoint
    }

    async makeRequest(action, method, body, headers) {
        let requestOptions = {
            method: method,
            body: body
        }

        if (headers) {
            requestOptions.headers = headers
        }

        try {
            let result = await fetch(this.endpoint + action, requestOptions)

            result = await result.text()

            try {
                return JSON.parse(result)
            } catch (e) {
                return result
            }
        } catch (e) {
            console.log("Request failed!")
            throw e
        }
    }

    async post(action, get_params, post_params) {
        if (!empty(get_params)) {
            action = action + '?' + new URLSearchParams(get_params).toString()
        }

        if (!post_params) {
            post_params = {}
        }

        for (const key of Object.keys(post_params)) {
            if (post_params[key] === null || post_params[key] === undefined) {
                console.log("Removing null value from post params", key)
                delete post_params[key]
            }
        }

        let data = new FormData()
        Object.keys(post_params).forEach((key) => {
            data.append(key, post_params[key])
        })

        return await this.makeRequest(action, 'POST', data)
    }

    async postJson(action, get_params, data) {
        if (!empty(get_params)) {
            action = action + '?' + new URLSearchParams(get_params).toString()
        }

        let body = JSON.stringify(data)

        return await this.makeRequest(action, 'POST', body, {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },)
    }

    async get(action, get_params) {
        if (!empty(get_params)) {
            action = action + '?' + new URLSearchParams(get_params).toString()
        }

        return await this.makeRequest(action, 'GET')
    }
}

export default HttpClient