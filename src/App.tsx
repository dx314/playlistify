import React, { KeyboardEvent, useEffect, useMemo, useRef, useState } from "react"
import "./App.css"
import { API_SERVER } from "./config"
import { ChatGPTResponse, ChatMessage } from "./typings"
import LoadingBar, { LoadingBarRef } from "react-top-loading-bar"
import ConnectSpotify from "./ConnectSpotify"
import SpotifyLogo from "./SpotifyLogo"
import { searchSongs, useSpotify } from "./utils/spotify"
import CreatePlaylist from "./CreatePlaylist"
import SongLI from "./SongLI"
import Loader from "./Loader"
import Info from "./Info"
import { getRandomPrompt } from "./utils"
import Thinking from "./Thinking"

function App() {
    const [prompt, setPrompt] = useState<string>(getRandomPrompt())
    const [hadFocus, setHadFocus] = useState<boolean>(false)
    const [prompted, setPrompted] = useState<string | null>(null)
    const [loading, setLoading] = useState<boolean | string>(true)
    const [thinking, setThinking] = useState<boolean>(false)
    const ref = useRef<LoadingBarRef>(null)

    const [now, setNow] = useState<Date>(new Date())

    const response = useRef<ChatMessage | null>(null)
    const { token: spotifyToken, user } = useSpotify()

    useEffect(() => {
        const playlist = localStorage.getItem("playlist")
        if (playlist) {
            response.current = JSON.parse(playlist)
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
            const msg: ChatMessage = JSON.parse(st)
            //setSongs(msg.split("\n"))
            response.current = msg
            setPrompted(prompt)
            setLoading(false)
            setThinking(false)
            if (response.current && response.current.songs && spotifyToken && user) {
                setLoading("#1DB954")
                searchSongs(response.current.songs, spotifyToken, user, (loading: boolean) => {
                    if (ref.current) {
                        ref.current.complete()
                        if (loading) ref.current.continuousStart()
                    }
                    setNow(new Date())
                }).then(() => {
                    setLoading(false)
                    localStorage.setItem("playlist", JSON.stringify(response.current))
                })
            }
        },
        [prompt],
    )

    const chat = async (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" && prompt !== "") {
            fetchPlaylist()
        }
    }
    console.log(spotifyToken)

    return (
        <div className={!loading ? "App" : "App loading"}>
            <div className={"searchbox"}>
                <input
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
                    onChange={(e) => setPrompt(e.target.value)}
                />
                {!loading && <button onClick={fetchPlaylist}>{prompted === prompt ? `Go again` : `Fetch Playlist`}</button>}
                {loading && (
                    <button disabled>
                        <Loader />
                    </button>
                )}
            </div>
            {thinking && <Thinking />}
            <LoadingBar color={typeof loading === "string" ? loading : "#008080"} ref={ref} />

            {response.current && (
                <div className={"playlist"}>
                    <h4>
                        {response.current.title}
                        {response.current && (
                            <CreatePlaylist
                                songs={response.current.songs}
                                spotifyToken={spotifyToken}
                                title={response.current.title}
                                description={response.current.description}
                                setLoading={setLoading}
                                user={user}
                            />
                        )}
                    </h4>
                    <ul>
                        {response.current.songs.map((song, i) => (
                            <SongLI key={`song-${i}`} song={song} index={i} />
                        ))}
                    </ul>
                    <p>{response.current.description}</p>
                </div>
            )}
            {!spotifyToken && <ConnectSpotify />}
            <Info />
        </div>
    )
}

export default App
