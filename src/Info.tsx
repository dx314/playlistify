import React from "react"

const Info: React.FC = () => (
    <div className={"info"}>
        <p>
            <strong>Welcome to plAIlist!</strong>
        </p>
        <p> plAIlist uses the ChatGPT AI to recommend playlists based on your prompts. To get started, follow these simple steps:</p>
        <ol>
            <li>Connect your Spotify account: Click the "Connect Spotify" button and log in to your Spotify account.</li>
            <li>
                Enter a prompt: In the text box, type a prompt that describes your music preferences, mood, song inspiration, or a theme for the playlist. E.g.
                "upbeat pop music for a workout" or "relaxing jazz for a rainy day."
            </li>
            <li>Click the ▷ button, and ChatGPT will generate a playlist based on your prompt.</li>
            <li>If you like the recommended playlist, click the "Add to Spotify" button to add it to your Spotify library.</li>
            <li>
                Open the playlist in Spotify: After the playlist has been added to your Spotify account, you can click the "Open in Spotify" button to open the
                playlist in Spotify's web player or desktop app.
            </li>
        </ol>
    </div>
)

export default Info
