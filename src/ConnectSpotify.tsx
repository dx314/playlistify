import * as React from "react"
import { AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE } from "./config"
import { spotifyURL } from "./utils/spotify"
import Modal from "./Modal"
import Info from "./Info"

const ConnectSpotify: React.FC = () => {
    return (
        <Modal isOpen={true} onClose={() => undefined}>
            <Info />
            <a className="spotify-connect" href={spotifyURL}>
                Connect Spotify
            </a>
        </Modal>
    )
}

export default ConnectSpotify
