import {JSX} from "solid-js";

const Title = (props: { children: JSX.Element, class?: string }) => {
    return (
        <h1 class={"leading-none tracking-tight text-gray-900 dark:text-white font-bold text-2xl md:text-3xl" + (props.class ? " " + props.class : "")}>{props.children}</h1>
    )
}

export default Title