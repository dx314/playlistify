import {Handler} from '@netlify/functions'
import {ChatGPTAPI} from 'chatgpt'

const handler: Handler = async (event, context) => {
    const api = new ChatGPTAPI({
        apiKey: process.env.OPENAI_API_KEY
    })
    const {msg} = event.queryStringParameters
    const res = await api.sendMessage(msg)

    return {
        statusCode: 200,
        body: JSON.stringify({message: "Hello World"})
    }
}

export {handler}