export const getVars = (hash: string): { [key: string]: string } => {
    if (hash.substring(0, 1) === "#") {
        hash = hash.substring(1)
    }
    return hash.split("&").reduce(function (res: { [key: string]: string }, item) {
        let parts = item.split("=")
        res[parts[0]] = parts[1]
        return res
    }, {})
}

export function isString(maybeString: any | string): maybeString is string {
    return typeof maybeString === "string"
}

export const SpotifyGreen = "#1DB954"

export const prompts = [
    "Upbeat summer vibes",
    "New Noise by Refused",
    "UK drum & bass",
    "songs to eat burgers to",
    "Relaxing evening jazz",
    "High-energy workout tunes",
    "Classic rock anthems",
    "Chill lo-fi beats for studying",
    "Epic movie soundtracks",
    "Melodic acoustic guitar",
    "90s nostalgia trip",
    "Road trip sing-alongs",
    "Energetic dance party",
    "Cozy rainy day songs",
    "Feel-good indie folk",
    "Mellow R&B for winding down",
    "Inspirational power ballads",
    "Late-night electronic grooves",
    "Soothing piano melodies",
    "Uplifting motivational tracks",
    "Classic 80s synth-pop",
    "Funky and soulful jams",
    "Intense metal workout",
]

export function getRandomPrompt() {
    const randomIndex = Math.floor(Math.random() * prompts.length)
    return prompts[randomIndex]
}
