import React from "react"
import "./scss/modal.scss"

type Props = {
    isOpen: boolean
    closeable: boolean
    children?: React.ReactNode
    onClose: () => void
}

const Modal: React.FC<Props> = ({ isOpen, onClose, children, closeable = false }) => {
    if (!isOpen) {
        return null
    }

    return (
        <div
            className="modal-backdrop"
            onClick={() => {
                if (closeable) onClose()
            }}
        >
            <div className="modal-container">{children}</div>
        </div>
    )
}

export default Modal
