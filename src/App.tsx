import React, { KeyboardEvent, useEffect, useRef, useState } from "react"
import "./App.css"
import { API_SERVER, AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE } from "./config"
import { ChatGPTResponse, ChatMessage } from "./typings"
import MusicLoader from "./loader"
import LoadingBar, { LoadingBarRef } from "react-top-loading-bar"
import ConnectSpotify from "./ConnectSpotify"

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
    const [prompt, setPrompt] = useState<string>("new noise by refused")
    const [songs, setSongs] = useState<string[]>([])
    const [loading, setLoading] = useState<boolean>(false)
    const ref = useRef<LoadingBarRef>(null)

    const [response, setResponse] = useState<ChatMessage | null>(null)

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

    useEffect(() => {
        if (loading && ref.current) {
            ref.current.continuousStart()
        } else if (ref.current) {
            ref.current.complete()
        }
    }, [loading])

    const chat = async (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" && prompt !== "") {
            setLoading(true)
            const res = await fetch(`${API_SERVER}/chat`, {
                method: "POST",
                headers: {
                    Accept: "application/json",
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ msg: prompt }),
            })
            const resp: ChatGPTResponse = await res.json()

            const st = resp.choices[0].message.content.replaceAll("\n", "")
            console.log(st)
            const msg: ChatMessage = JSON.parse(st)
            //setSongs(msg.split("\n"))
            setResponse(msg)
            setLoading(false)
        }
    }

    return (
        <div className="App">
            <LoadingBar color="#f11946" ref={ref} />
            <div>
                <input onKeyDown={chat} value={prompt} disabled={loading} onChange={(e) => setPrompt(e.target.value)} />
                {response && (
                    <div>
                        <h4>{response.title}</h4>
                        <ul>
                            {response.songs.map((song, i) => (
                                <li key={`song-${i}`}>{`${i + 1}: ${song.title} - ${song.artist}`}</li>
                            ))}
                        </ul>
                        <p>{response.description}</p>
                    </div>
                )}
            </div>
            {!token && <ConnectSpotify />}
        </div>
    )
}

export default App
