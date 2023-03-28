import * as React from "react"
import { MdPlayArrow, MdReplay } from "react-icons/md"
import EllipsisLoading from "./EllipsisLoading"
import "./scss/input.scss"

type Props = {
    onButton: () => void
    isLoading: boolean
    isRefresh: boolean
} & React.InputHTMLAttributes<HTMLInputElement>

const Input: React.FC<Props> = ({ onButton, isLoading, isRefresh, ...props }) => {
    return (
        <div className="plailist-input">
            <input className="" type="text" {...props}></input>
            {
                <button disabled={props.disabled} onClick={onButton}>
                    {isLoading ? <EllipsisLoading /> : !isRefresh ? <MdPlayArrow /> : <MdReplay />}
                </button>
            }
        </div>
    )
}

export default Input
