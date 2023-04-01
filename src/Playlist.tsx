import * as React from "react"
import CreatePlaylist from "./CreatePlaylist"
import SongLI from "./SongLI"
import { AIPlaylist } from "./typings"
import { searchSongs, useSpotify } from "./utils/spotify"
import { useEffect } from "react"
import { SpotifyUser } from "./spotifyTypes"

interface Props {
    playlist: AIPlaylist
    spotifyToken: string
    setLoading: (v: string | boolean) => void
    setPlaylist: (playlist: AIPlaylist) => void
    user: SpotifyUser
}

const Playlist: React.FC<Props> = ({ playlist, setPlaylist, spotifyToken, setLoading, user }) => {
    useEffect(() => {
        const search = async () => {
            if (playlist && playlist.songs && spotifyToken && user) {
                setLoading("#1DB954")
                const p: AIPlaylist | null = await searchSongs(playlist, spotifyToken, user)
                setLoading(false)
                if (p) {
                    setPlaylist(p)
                    localStorage.setItem("playlist", JSON.stringify(p))
                }
            }
        }
        search()
        if (playlist) localStorage.setItem("playlist", JSON.stringify(playlist))
    }, [playlist])
    return (
        <div className={"playlist"}>
            <h4>{playlist.title}</h4>
            <ul>
                {playlist.songs.map((song, i) => (
                    <SongLI key={`song-${i}`} song={song} index={i} />
                ))}
            </ul>
            <p>{playlist.description}</p>
            <CreatePlaylist
                songs={playlist.songs}
                spotifyToken={spotifyToken}
                title={playlist.title}
                description={playlist.description}
                setLoading={setLoading}
                user={user}
            />
        </div>
    )
}

export default Playlist
