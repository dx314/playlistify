import * as React from "react"
import { Song } from "./typings"
import SpotifyLogo from "./SpotifyLogo"
import { SpotifyGreen } from "./utils"

const SongLI: React.FC<{ song: Song; index: number }> = ({ song, index }) => {
    return (
        <li className={"song"}>
            {song.spotifyId && (
                <a href={`https://open.spotify.com/track/${song.spotifyId}`} target={"_blank"}>
                    <SpotifyLogo height={"16px"} width={"16px"} color={SpotifyGreen} />
                </a>
            )}
            {!song.spotifyId && <SpotifyLogo height={"16px"} width={"16px"} color={"#2d2d2d"} />}
            <span>{`${index + 1}.`}</span>
            {`${song.title} - ${song.artist}`}{" "}
        </li>
    )
}

export default SongLI
