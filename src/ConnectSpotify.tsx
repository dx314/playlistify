import * as React from "react"
import { AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE } from "./config"
import { spotifyURL } from "./utils/spotify"

const ConnectSpotify: React.FC = () => {
    return (
        <div>
            <button className={"spotify-connect"}>
                <a href={spotifyURL}>Connect Spotify</a>
            </button>
        </div>
    )
}

export default ConnectSpotify
