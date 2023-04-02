import * as React from "react"
import { Song } from "./typings"
import { createPlaylistAndAddTracks, errAtom } from "./utils/spotify"
import { isString, SpotifyGreen } from "./utils"
import { SpotifyUser } from "./spotifyTypes"
import { useAtom } from "jotai/index"

interface Props {
    songs: Song[]
    title: string
    description: string
    user?: SpotifyUser | null
    spotifyToken?: string | null
    setLoading: (v: string | boolean) => void
}

const CreatePlaylist: React.FC<Props> = ({ songs, title, description, user, spotifyToken, setLoading }) => {
    const songIds: string[] = songs.map((song) => song.spotifyId).filter(isString)
    const [, setError] = useAtom(errAtom)
    return (
        <span
            className={"spotify-add-playlist" + (!user || !spotifyToken ? " disabled" : "")}
            onClick={() => {
                if (!user || !spotifyToken) return
                setLoading(SpotifyGreen)
                createPlaylistAndAddTracks(user, spotifyToken, title, songIds, description)
                    .then((playlistId) => {
                        const playlistUrl = `https://open.spotify.com/playlist/${playlistId}`
                        window.open(playlistUrl, "_blank")
                        setLoading(false)
                    })
                    .catch((err: any) => {
                        setError(err.message)
                    })
            }}
        >
            add to spotify
        </span>
    )
}

export default CreatePlaylist
