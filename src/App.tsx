import React, { KeyboardEvent, useEffect, useState } from "react"
import "./App.css"
import { AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE } from "./config"

const getVars = (hash: string): { [key: string]: string } => {
    if (hash.substring(0, 1) === "#") {
        hash = hash.substring(1)
    }
    return hash.split("&").reduce(function (res: { [key: string]: string }, item) {
        let parts = item.split("=")
        res[parts[0]] = parts[1]
        return res
    }, {})
}

function App() {
    const [token, setToken] = useState<string | null>(null)
    const [prompt, setPrompt] = useState<string>("playlist based on new noise by refused")
    const [response, setResponse] = useState<string>("")

    useEffect(() => {
        const hash = window.location.hash
        let token = window.localStorage.getItem("token")

        if (!token && hash) {
            token = getVars(hash)["access_token"]
            if (token && token !== "") {
                window.location.hash = ""
                window.localStorage.setItem("token", token)
            }
        }

        setToken(token)
    }, [])

    const chat = async (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" && prompt !== "") {
            const res = await fetch(`/.netlify/functions/chat?msg=${prompt}`)
            const text = await res.text()
            setResponse(text)
        }
    }

    return (
        <div className="App">
            <div>
                {!token && <a href={`${AUTH_ENDPOINT}?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&response_type=${RESPONSE_TYPE}`}>Login to Spotify</a>}
                <input onKeyDown={chat} value={prompt} onChange={(e) => setPrompt(e.target.value)} />
            </div>
            <div>{response}</div>
        </div>
    )
}

export default App
