import { useEffect, useState } from "react"
import { getVars } from "./index"
import { AIPlaylist, Song } from "../typings"
import { atom, useAtom } from "jotai"
import { API_SERVER, AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE } from "../config"
import { LoadingBarRef } from "react-top-loading-bar"
import { useLocation } from "react-router-dom"
import { SearchResults, SpotifyUser } from "../spotifyTypes"

const tokenKey = "spotify_token"
const userKey = "user"
export const spotifyTokenAtom = atom<string | null>(localStorage.getItem(tokenKey))
const initialUser: SpotifyUser | null = JSON.parse(`${localStorage.getItem(userKey)}`)
export const userAtom = atom<SpotifyUser | null>(initialUser)
export const errAtom = atom<string | null>(null)

const scopes = encodeURIComponent(
    ["user-read-private", "playlist-modify-private", "playlist-modify-public", "playlist-read-private", "playlist-read-collaborative"].join(" "),
)
export const spotifyURL = `${AUTH_ENDPOINT}?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&response_type=${RESPONSE_TYPE}&scope=${scopes}`

const fetchSpotify = async <T>(spotify_path: string): Promise<T> => {
    const url = `/api/spotify/api${spotify_path}`
    const response = await fetch(url, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    })
    if (!response.ok) {
        throw new Error(`Spotify API request failed with status ${response.status}`)
    }
    return await response.json()
}

export const searchSongs = async (playlist: AIPlaylist, access_token: string, user: SpotifyUser): Promise<AIPlaylist | null> => {
    const songs: Song[] = [...playlist.songs]
    let count = 0
    for (const song of playlist.songs) {
        if (song.spotifyId || song.spotifyId === null) {
            continue
        }
        try {
            console.log("searching spotify for " + song.title)
            song.spotifyId = await searchSpotify(song.title, song.artist, access_token, user.country)
            count++
        } catch (err) {
            console.error(err)
            continue
        }
    }

    if (count === 0) return null

    const newPlaylist = { ...playlist, songs }
    return newPlaylist
}

async function searchSpotify(title: string, artist: string, accessToken: string, market: string = "US"): Promise<string | null> {
    const query = `track:${encodeURIComponent(title)} artist:${encodeURIComponent(artist)}`
    const types = "track"
    const limit = 1
    const include_external = "audio"

    const url = `/search?q=${query}&type=${types}&market=${market}&limit=${limit}&include_external=${include_external}`
        .replace(/ *\([^)]*\) */g, "")
        .replaceAll("'", "")

    const data: { tracks: SearchResults } = await fetchSpotify(url)

    if (!data.tracks.items || data.tracks.items.length === 0) {
        return null
    }

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
    console.log(user, "user")
    const createPlaylistUrl = `https://api.spotify.com/v1/users/${user.spotify_id}/playlists`
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
        let text: string = ""
        try {
            text = await response.text()
        } catch {
            text = "no error body"
        }

        throw new Error(`Spotify API request failed with status ${response.statusText}: ${text}`)
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
        let text: string = ""
        try {
            text = await response.text()
        } catch {
            text = "no error body"
        }

        throw new Error(`Spotify API request failed with status ${response.statusText}: ${text}`)
    }
}

async function fetchUser(accessToken: string, reToken: () => void): Promise<SpotifyUser> {
    const url = "https://api.spotify.com/v1/me"

    const response = await fetch(url, {
        method: "GET",
        headers: {
            Authorization: `Bearer ${accessToken}`,
            "Content-Type": "application/json",
        },
    })

    if (!response.ok) {
        reToken()
        throw new Error(`Spotify API request failed with status ${response.status}`)
    }

    const data: SpotifyUser = await response.json()
    return data
}

interface SpotifyProperties {
    token: string | null
    user: SpotifyUser | null
    clearSpotify: () => void
    verified: boolean
}

type RefreshTokenResponse = {
    access_token: string
}

export const useSpotify = (): SpotifyProperties => {
    const [token, setAccessToken] = useAtom(spotifyTokenAtom)
    const [user, setUser] = useAtom(userAtom)
    const [verified, setVerified] = useState<boolean>(false)
    const location = useLocation()
    const [err, setError] = useAtom(errAtom)

    useEffect(() => {
        const getUser = async () => {
            const response = await fetch("/api/me", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            })

            setVerified(true)

            if (!response.ok) {
                let text: string = ""
                try {
                    text = await response.text()
                } catch {
                    text = "no error body"
                }
                throw new Error(`Error retrieving user: ${response.statusText}: ${text}`)
            }

            try {
                const me: SpotifyUser = await response.json()
                setUser(me)
                setAccessToken(me.access_token)
            } catch (err) {
                console.error("unable to save user", err)
            }
        }
        getUser()
    }, [])

    const clearSpotify = () => {
        setAccessToken(null)
        setUser(null)
    }

    useEffect(() => {
        // Get the access_token parameter from the URL
        const searchParams = new URLSearchParams(location.search)
        const token = searchParams.get("access_token")
        setAccessToken(token === "" ? null : token)
    }, [location.search])

    return { token, user, clearSpotify, verified }
}
