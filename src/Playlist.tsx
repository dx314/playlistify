import * as React from "react"
import CreatePlaylist from "./CreatePlaylist"
import SongLI from "./SongLI"
import { ChatMessage } from "./typings"
import { SpotifyUser } from "./utils/spotify"

interface Props {
    playlist: ChatMessage
    spotifyToken: string
    setLoading: (v: string | boolean) => void
    user: SpotifyUser
}

const Playlist: React.FC<Props> = ({ playlist, spotifyToken, setLoading, user }) => (
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

export default Playlist
