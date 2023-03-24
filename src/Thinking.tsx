import React, { useState, useEffect } from "react"

interface LoadingMessageProps {
    message: string
}

const thinkingMessages = [
    "Crunching neural circuits",
    "Beaming AI brainwaves",
    "Synthesizing deep thoughts",
    "Rummaging through the data vault",
    "Assembling a symphony of ideas",
    "Revving up the cognitive engines",
    "Feeding the knowledge hamsters",
    "Generating sparks of creativity",
    "Churning the idea gears",
    "Computing melodies and harmonies",
    "Delving into the musical matrix",
    "Blending genres and rhythms",
    "Diving into the ocean of tunes",
    "Pondering the infinite soundscapes",
    "Unlocking the vault of musical wisdom",
    "Whipping up a symphonic storm",
    "Activating sonic superpowers",
    "Warming up the imagination turbines",
    "Fine-tuning the auditory algorithms",
    "Weaving a tapestry of musical gems",
    "Replacing human taste makers",
]

const getRandomMessage = (): string => {
    const randomIndex = Math.floor(Math.random() * thinkingMessages.length)
    return thinkingMessages[randomIndex]
}

const Thinking: React.FC = () => {
    const [dots, setDots] = useState<string>(".")
    const [message, setMessage] = useState<string>("")

    useEffect(() => {
        setMessage(getRandomMessage())
        const interval = setInterval(() => {
            setDots((prevDots) => (prevDots.length < 3 ? prevDots + "." : "."))
        }, 500)

        return () => {
            clearInterval(interval)
        }
    }, [])

    return (
        <div className={"thinking"}>
            {message} {dots}
        </div>
    )
}

export default Thinking
