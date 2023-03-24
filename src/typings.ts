export interface ChatGPTResponse {
    id: string
    object: string
    created: number
    choices: Choice[]
}

export interface Choice {
    index: number
    message: Message
    finish_reason: string
    usage: Usage
}

export interface Message {
    role: string
    content: string
}

export interface Usage {
    prompt_tokens: number
    completion_tokens: number
    total_tokens: number
}

export interface Song {
    artist: string
    title: string
    spotifyId?: string | null
}

export interface ChatMessage {
    songs: Song[]
    description: string
    title: string
}
