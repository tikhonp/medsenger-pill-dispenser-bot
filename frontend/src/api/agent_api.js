import HttpClient from './http_client.js'

let API_URL = import.meta.env.VITE_AGENT_URL || "";

class AgentApiClient {
  constructor(host) {
    this.client = new HttpClient(host)
  }

  async get2x2Layout() {
    let result
    try {
      result = await this.client.get('/confirm-layout/small')
    } catch (e) {
      console.log(e)
      throw new Error('Ошибка соединения с сервером.')
    }

    return result
  }

  async getLargeLayout() {
    let result
    try {
      result = await this.client.get('/confirm-layout/large')
    } catch (e) {
      console.log(e)
      throw new Error('Ошибка соединения с сервером.')
    }

    return result
  }

}

const api = new AgentApiClient(API_URL)

export default api
