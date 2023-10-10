import {Channel, JobState, State} from "../utils/api";
import {Show} from "solid-js";

const ChannelActions = (props: { channel: Channel }) => {
    // TODO
    return (
        <div>
            <Show when={props.channel.state == State.COMPLETED}>
                <button title={"Clear"}>🗑</button>    
            </Show>
            <Show when={props.channel.job_state == JobState.RUNNING}>
                <button title={"Stop Job"}>🛑</button>
            </Show>
            <Show when={props.channel.job_state != undefined && props.channel.job_state != JobState.RUNNING}>
                <button title={"Clear broken job"}>👌</button>
            </Show>
        </div>
        
    )
}

export default ChannelActions