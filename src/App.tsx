import React, {useEffect, useState} from 'react'
import './App.css'
import {AUTH_ENDPOINT, CLIENT_ID, REDIRECT_URI, RESPONSE_TYPE} from "./config";

function App() {
    const [token, setToken] = useState<String | null>(null)

    useEffect(() => {
        const hash = window.location.hash
        let token = window.localStorage.getItem("token")

        if (!token && hash) {
            token = hash.substring(1).split("&").find(elem => elem.startsWith("access_token")).split("=")[1]

            window.location.hash = ""
            window.localStorage.setItem("token", token)
        }

        setToken(token)







    }, [])

    return (
        <div className="App">
            <div>
                {
                    !token && <a href={`${AUTH_ENDPOINT}?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&response_type=${RESPONSE_TYPE}`}>Login
                        to Spotify</a>}
                <input onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        console.log('do validate');
                    }
                }} />
            </div>
        </div>
    )
}

export default App
