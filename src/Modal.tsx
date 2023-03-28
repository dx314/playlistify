import React from "react"
import "./scss/modal.scss"

type Props = {
    isOpen: boolean
    children?: React.ReactNode
    onClose: () => void
}

const Modal: React.FC<Props> = ({ isOpen, onClose, children }) => {
    if (!isOpen) {
        return null
    }

    return (
        <div className="modal-backdrop">
            <div className="modal-container">{children}</div>
        </div>
    )
}

export default Modal
