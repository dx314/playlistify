import React, { useState, useEffect } from "react"

const EllipsisLoading: React.FC = () => {
    const [dots, setDots] = useState<string>("")

    useEffect(() => {
        const interval = setInterval(() => {
            setDots((prevDots) => (prevDots.length < 3 ? prevDots + "." : ""))
        }, 300)

        return () => clearInterval(interval)
    }, [])

    return <span>{dots}</span>
}

export default EllipsisLoading
