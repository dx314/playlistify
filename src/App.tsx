import React, { KeyboardEvent, useEffect, useMemo, useRef, useState } from "react"
import "./scss/App.scss"
import "./media.css"
import { API_SERVER } from "./config"
import { ChatGPTResponse, AIPlaylist } from "./typings"
import LoadingBar, { LoadingBarRef } from "react-top-loading-bar"
import ConnectSpotify from "./ConnectSpotify"
import { errAtom, searchSongs, useSpotify } from "./utils/spotify"
import Info from "./Info"
import { getRandomPrompt } from "./utils"
import Thinking from "./Thinking"
import Input from "./Input"
import Playlist from "./Playlist"
import Brand from "./Brand"
import Loader from "./Loader"
import { useAtom } from "jotai"
import Modal from "./Modal"

function App() {
    const [prompt, setPrompt] = useState<string>(getRandomPrompt())
    const [hadFocus, setHadFocus] = useState<boolean>(false)
    const [prompted, setPrompted] = useState<string | null>(null)
    const [loading, setLoading] = useState<boolean | string>(true)
    const [thinking, setThinking] = useState<boolean>(false)
    const [err, setError] = useAtom(errAtom)

    const ref = useRef<LoadingBarRef>(null)

    const [now, setNow] = useState<Date>(new Date())

    const [playlist, setPlaylist] = useState<AIPlaylist | null>(null)
    const { token: spotifyToken, user, verified } = useSpotify()

    useEffect(() => {
        const playlist = localStorage.getItem("playlist")
        if (playlist) {
            setPlaylist(JSON.parse(playlist))
        }
        setLoading(false)
    }, [])

    useEffect(() => {
        if (loading && ref.current) {
            ref.current.continuousStart()
        } else if (ref.current) {
            ref.current.complete()
        }
    }, [loading])

    const fetchPlaylist = useMemo(
        () => async () => {
            if (!prompt || prompt === "" || prompt.length < 3) return
            setLoading(true)
            setThinking(true)

            try {
                const res = await fetch(`${API_SERVER}/chat`, {
                    method: "POST",
                    headers: {
                        Accept: "application/json",
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({ msg: prompt }),
                })
                if (!res.ok) {
                    let text: string = ""
                    try {
                        text = await res.text()
                    } catch {
                        text = "no error body"
                    }
                    throw new Error(`${res.statusText}: ${text}`)
                }
                const resp: ChatGPTResponse = await res.json()

                const st = resp.choices[0].message.content.replaceAll("\n", "")

                const msg: AIPlaylist = JSON.parse(st)
                //setSongs(msg.split("\n"))

                setPlaylist(msg)
                setPrompted(prompt)
                setLoading(false)
                setThinking(false)
            } catch (err: any) {
                setError(err.message)
                return
            }
        },
        [prompt],
    )

    const chat = async (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" && prompt !== "") {
            fetchPlaylist()
        }
    }

    if (!verified) {
        return (
            <div className={!loading ? "App" : "App loading"}>
                <LoadingBar color={typeof loading === "string" ? loading : "#008080"} ref={ref} />
                <Loader />
            </div>
        )
    }

    return (
        <div className={!loading ? "App" : "App loading"}>
            <Brand />
            <div className={"plailist-container"}>
                <div className={"searchbox"}>
                    <Input
                        onFocus={() => {
                            if (!hadFocus) {
                                setPrompt("")
                                setHadFocus(true)
                            }
                        }}
                        placeholder={`enter a prompt, such as "${getRandomPrompt()}"`}
                        onKeyDown={chat}
                        value={prompt}
                        disabled={loading === true}
                        isLoading={!!loading}
                        isRefresh={prompted === prompt}
                        onButton={fetchPlaylist}
                        onChange={(e) => setPrompt(e.target.value)}
                    />
                </div>
                <div className={"spacer"} />
                {thinking && <Thinking />}
                <LoadingBar color={typeof loading === "string" ? loading : "#008080"} ref={ref} />

                {playlist && spotifyToken && user && (
                    <Playlist setPlaylist={setPlaylist} playlist={playlist} spotifyToken={spotifyToken} setLoading={setLoading} user={user} />
                )}
            </div>
            {!user && <ConnectSpotify />}
            {err && err !== "" && (
                <Modal isOpen={true} onClose={() => setError(null)} closeable={true}>
                    <p>{err}</p>
                </Modal>
            )}
            <Info />
        </div>
    )
}

export default App
