import * as React from "react"
import { AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE } from "./config"

const ConnectSpotify = () => {
    return (
        <div>
            <button className={"spotify-connect"}>
                <a href={`${AUTH_ENDPOINT}?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&response_type=${RESPONSE_TYPE}`}>Create Playlist</a>
            </button>
        </div>
    )
}

export default ConnectSpotify
