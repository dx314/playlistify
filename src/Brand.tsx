import * as React from "react"
import { SVGProps } from "react"
import "./scss/brand.scss"
import LogoMark from "./LogoMark"
import Logo from "./Logo"

const Brand = () => (
    <div className={"brand"}>
        <div>
            <LogoMark width={"40px"} />
            <Logo width={"70px"} />
        </div>
    </div>
)

export default Brand
