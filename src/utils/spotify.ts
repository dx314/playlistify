import { useEffect } from "react"
import { getVars } from "./index"
import { Song } from "../typings"
import { atom, useAtom } from "jotai"
import { AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE } from "../config"
import { LoadingBarRef } from "react-top-loading-bar"

const tokenKey = "spotify_token"
const userKey = "user"
export const spotifyTokenAtom = atom<string | null>(localStorage.getItem(tokenKey))
const initialUser: SpotifyUser | null = JSON.parse(`${localStorage.getItem(userKey)}`)
export const userAtom = atom<SpotifyUser | null>(initialUser)

const scopes = encodeURIComponent(
    ["user-read-private", "playlist-modify-private", "playlist-modify-public", "playlist-read-private", "playlist-read-collaborative"].join(" "),
)
export const spotifyURL = `${AUTH_ENDPOINT}?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&response_type=${RESPONSE_TYPE}&scope=${scopes}`

export const searchSongs = async (songs: Song[], access_token: string, user: SpotifyUser, loader: (loading: boolean) => void): Promise<Song[]> => {
    for (const song of songs) {
        loader(false)
        console.log("searching!!")
        song.spotifyId = await searchSpotify(song.title, song.artist, access_token, user.country)
        loader(true)
    }

    return songs
}

async function searchSpotify(title: string, artist: string, accessToken: string, market: string = "US"): Promise<string | null> {
    const query = `track:${encodeURIComponent(title)} artist:${encodeURIComponent(artist)}`
    const types = "track"
    const limit = 1
    const include_external = "audio"

    const url = `https://api.spotify.com/v1/search?q=${query}&type=${types}&market=${market}&limit=${limit}&include_external=${include_external}`

    const response = await fetch(url, {
        method: "GET",
        headers: {
            Authorization: `Bearer ${accessToken}`,
            "Content-Type": "application/json",
        },
    })

    if (!response.ok) {
        throw new Error(`Spotify API request failed with status ${response.status}`)
    }

    const data = await response.json()
    const track = data.tracks.items[0]

    return track ? track.id : null
}

export async function createPlaylistAndAddTracks(
    user: SpotifyUser,
    accessToken: string,
    playlistName: string,
    trackIds: string[],
    description: string = "",
): Promise<string> {
    const createPlaylistUrl = `https://api.spotify.com/v1/users/${user.id}/playlists`
    const playlistData = await createPlaylist(accessToken, createPlaylistUrl, playlistName, description)
    const playlistId = playlistData.id

    const addTracksUrl = `https://api.spotify.com/v1/playlists/${playlistId}/tracks`
    await addTracksToPlaylist(accessToken, addTracksUrl, trackIds)

    return playlistId
}

async function createPlaylist(accessToken: string, url: string, name: string, description: string): Promise<any> {
    const response = await fetch(url, {
        method: "POST",
        headers: {
            Authorization: `Bearer ${accessToken}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, description }),
    })

    if (!response.ok) {
        throw new Error(`Spotify API request failed with status ${response.status}`)
    }

    return response.json()
}

async function addTracksToPlaylist(accessToken: string, url: string, trackIds: string[]): Promise<void> {
    const uris = trackIds.map((id) => `spotify:track:${id}`)

    const response = await fetch(url, {
        method: "POST",
        headers: {
            Authorization: `Bearer ${accessToken}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ uris }), // Send the 'uris' as a JSON array in the request body
    })

    if (!response.ok) {
        throw new Error(`Spotify API request failed with status ${response.status}`)
    }
}

async function fetchUser(accessToken: string): Promise<SpotifyUser> {
    const url = "https://api.spotify.com/v1/me"

    const response = await fetch(url, {
        method: "GET",
        headers: {
            Authorization: `Bearer ${accessToken}`,
            "Content-Type": "application/json",
        },
    })

    if (!response.ok) {
        throw new Error(`Spotify API request failed with status ${response.status}`)
    }

    const data: SpotifyUser = await response.json()
    return data
}

interface SpotifyProperties {
    token: string | null
    user: SpotifyUser | null
}

export const useSpotify = (): SpotifyProperties => {
    const [token, setSpotifyToken] = useAtom(spotifyTokenAtom)
    const [user, setUser] = useAtom(userAtom)

    useEffect(() => {
        const hash = window.location.hash
        let token: string | null = localStorage.getItem(tokenKey)

        if (!token && hash) {
            token = getVars(hash)["access_token"]
            if (token && token !== "") {
                window.location.hash = ""
                window.localStorage.setItem(tokenKey, token)
            }
        }

        if (token) {
            fetchUser(token).then((user) => {
                setUser(user)
                localStorage.setItem(userKey, JSON.stringify(user))
            })
        }

        setSpotifyToken(token)
        if (token) localStorage.setItem(tokenKey, token)
        else localStorage.removeItem(tokenKey)
    }, [])
    return { token, user }
}

export interface SpotifyUser {
    country: string
    display_name: string
    email: string
    external_urls: {
        spotify: string
    }
    followers: {
        href: string | null
        total: number
    }
    href: string
    id: string
    images: {
        height: number | null
        url: string
        width: number | null
    }[]
    product: "free" | "open" | "premium" | "unknown"
    type: "user"
    uri: string
}
